package inputs

import (
	"fmt"
	"math"
	"os"

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
	Name  string
	Value float64
}

type StudyInput struct {
	SimInputNames   []string
	SimInputSamples [][]float64
}

func CalculateStep(min float64, max float64, numSims int) float64 {
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
