package inputs

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/bellh14/DesignManager/pkg/types"
	"github.com/bellh14/DesignManager/pkg/utils"
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

func CalculateStep(min float64, max float64, numSims int) float64 {
	if min == max {
		return 0
	}
	return (math.Abs(min) + math.Abs(max)) / float64(numSims-1)
}

func GenerateSimInputs(designParameters []types.DesignParameter) []SimInput {
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

func GenerateStudyInputs(simInputs []SimInput) StudyInput {
	var studyInputs StudyInput
	numParams := len(simInputs)
	if numParams == 0 {
		return studyInputs
	}

	// Assuming all SimInput have the same NumSims for simplicity
	numSims := simInputs[0].NumSims
	studyInputs.SimInputSamples = make([][]float64, numSims)
	for i := range studyInputs.SimInputSamples {
		studyInputs.SimInputSamples[i] = make([]float64, numParams)
	}

	for i, si := range simInputs {
		studyInputs.SimInputNames = append(studyInputs.SimInputNames, si.Name)
		for j := 0; j < si.NumSims; j++ {
			studyInputs.SimInputSamples[j][i] = si.Min + float64(j)*si.Step
		}
	}

	return studyInputs
}

func GenerateSimInputCSV(studyInput StudyInput, fileName string) error {
	inputFile, err := os.Create(fileName + ".csv")
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer inputFile.Close()
	utils.WriteParameterCsvHeader(studyInput.SimInputNames, inputFile)
	utils.WriteParameterCsv(studyInput.SimInputSamples, inputFile)
	return nil
}

func HandleSimInputs(designParameters []types.DesignParameter, fileName string) error {
	simInputs := GenerateSimInputs(designParameters)
	studyInputs := GenerateStudyInputs(simInputs)
	err := GenerateSimInputCSV(studyInputs, fileName)
	if err != nil {
		return err
	}
	return nil
}

func SimInputByJobNumber(inputFileName string, jobNumber int) (SimInputIteration, error) {
	simInputIteration := SimInputIteration{}
	file, err := os.Open(inputFileName + ".csv")
	if err != nil {
		return simInputIteration, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return simInputIteration, err
	}

	if jobNumber >= len(records)-1 {
		return simInputIteration, fmt.Errorf("Job number %d is out of range", jobNumber)
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
