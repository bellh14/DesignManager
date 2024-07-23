package inputs

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/utils"
	extraMath "github.com/bellh14/DesignManager/pkg/utils/math"
)

type SimInput struct {
	Name    string
	Min     float64
	Max     float64
	Step    float64
	NumSims int
}

type SimInputIteration struct {
	Name  []string
	Value []float64
}

type StudyInput struct {
	SimInputNames   []string
	SimInputSamples [][]float64
}

type SimInputGenerator struct {
	SimInputsFileName string
	DesignParameters  []config.DesignParameter
	TotalNumSims      int
}

func NewSimInputGenerator(
	designParameters []config.DesignParameter,
	fileName string,
	totalNumSims int,
) *SimInputGenerator {
	return &SimInputGenerator{
		SimInputsFileName: fileName,
		DesignParameters:  designParameters,
		TotalNumSims:      totalNumSims,
	}
}

func CalculateStep(min float64, max float64, numSims int) float64 {
	if min == max {
		return 0
	}
	var step float64
	if min > 0 && numSims <= 2 {
		step = (max - min)
	} else if min > 0 && numSims > 2 {
		step = (max - min) / float64(numSims-1)
	} else {
		step = (math.Abs(min) + math.Abs(max)) / float64(numSims-1)
	}
	return step
}

func GenerateSimInputs(designParameters []config.DesignParameter) []SimInput {
	var simInputs []SimInput
	for _, dp := range designParameters {
		simInputs = append(simInputs, SimInput{
			Name:    dp.Name,
			Min:     dp.Min,
			Max:     dp.Max,
			Step:    CalculateStep(dp.Min, dp.Max, dp.NumSims),
			NumSims: dp.NumSims,
		})
	}
	return simInputs
}

// Generates all input values for given parameter
func GenerateValues(simInput SimInput) []float64 {
	var values []float64

	for val := simInput.Min; val < simInput.Max || extraMath.AlmostEqual(val, simInput.Max); val += simInput.Step {
		val = extraMath.RoundToDecimalPlaces(val, 4)
		values = append(values, val)
	}
	return values
}

func CombineCombinations(combinations [][]float64, newValues []float64) [][]float64 {
	var newCombinations [][]float64
	for _, combination := range combinations {
		for _, newValue := range newValues {
			newCombination := append([]float64{}, combination...)
			newCombination = append(newCombination, newValue)
			newCombinations = append(newCombinations, newCombination)
		}
	}
	return newCombinations
}

func GenerateStudyInputs(simInputs []SimInput, numSims int) StudyInput {
	var studyInputs StudyInput
	numParams := len(simInputs)
	if numParams == 0 {
		return studyInputs
	}

	// Assuming all SimInput have the same NumSims for simplicity
	// numSims = simInputs[0].NumSims
	// for i := range studyInputs.SimInputSamples {
	// 	studyInputs.SimInputSamples[i] = make([]float64, numParams)
	// }

	initialValues := GenerateValues(simInputs[0])
	studyInputs.SimInputSamples = make([][]float64, len(initialValues))

	studyInputs.SimInputNames = append(studyInputs.SimInputNames, simInputs[0].Name)
	for i, val := range initialValues {
		studyInputs.SimInputSamples[i] = []float64{val}
	}

	for i := 1; i < len(simInputs); i++ {
		studyInputs.SimInputNames = append(studyInputs.SimInputNames, simInputs[i].Name)
		studyInputs.SimInputSamples = CombineCombinations(
			studyInputs.SimInputSamples,
			GenerateValues(simInputs[i]),
		)
	}
	//
	// for i, si := range simInputs {
	// 	studyInputs.SimInputNames = append(studyInputs.SimInputNames, si.Name)
	// 	for j := 0; j < si.NumSims; j++ {
	// 		logger.Debug(fmt.Sprintf("i: %d, j: %d", i, j))
	// 		studyInputs.SimInputSamples[j][i] = si.Min + float64(j)*si.Step
	// 	}
	// }

	return studyInputs
}

func GenerateSimInputCSV(studyInput StudyInput, fileName string) error {
	inputFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer inputFile.Close()
	utils.WriteParameterCsvHeader(studyInput.SimInputNames, inputFile)
	utils.WriteParameterCsv(studyInput.SimInputSamples, inputFile)
	return nil
}

func (simInputGenerator *SimInputGenerator) HandleSimInputs() error {
	// pathetic why did I write these seperate, am confused now but will fix later
	simInputs := GenerateSimInputs(simInputGenerator.DesignParameters)
	studyInputs := GenerateStudyInputs(simInputs, simInputGenerator.TotalNumSims)
	err := GenerateSimInputCSV(studyInputs, simInputGenerator.SimInputsFileName)
	if err != nil {
		return err
	}
	return nil
}

func (simInputGenerator *SimInputGenerator) SimInputByJobNumber(
	jobNumber int,
) (SimInputIteration, error) {
	simInputIteration := SimInputIteration{}
	file, err := os.Open(simInputGenerator.SimInputsFileName)
	if err != nil {
		return simInputIteration, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return simInputIteration, err
	}

	if jobNumber > len(records)-1 {
		return simInputIteration, fmt.Errorf("job number %d is out of range", jobNumber)
	}

	for i, record := range records {
		if i == 0 {
			simInputIteration.Name = append(simInputIteration.Name, record...)
		} else {
			if i == jobNumber {
				for _, value := range record {
					parsedValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return simInputIteration, err
					}
					simInputIteration.Value = append(simInputIteration.Value, parsedValue)
				}
			}
		}
	}
	return simInputIteration, nil
}
