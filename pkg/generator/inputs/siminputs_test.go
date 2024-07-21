package inputs_test

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/utils/math"
)

func TestCalculateStep(t *testing.T) {
	min := -1.3
	max := 1.3
	numSims := 9
	expected := 0.325
	actual := inputs.CalculateStep(min, max, numSims)
	if actual != expected {
		t.Errorf("Expected %f, got %f", expected, actual)
	}
}

func TestCalculateStepZero(t *testing.T) {
	min := 0.0
	max := 0.0
	numSims := 9
	expected := 0.0
	actual := inputs.CalculateStep(min, max, numSims)
	if actual != expected {
		t.Errorf("Expected %f, got %f", expected, actual)
	}
}

func TestGenerateSimInputs(t *testing.T) {
	designParameters := []config.DesignParameter{
		{
			Name:    "Parameter1",
			Min:     -1.3,
			Max:     1.3,
			NumSims: 3,
		},
		{
			Name:    "Parameter2",
			Min:     -1.3,
			Max:     1.3,
			NumSims: 3,
		},
	}
	expected := []inputs.SimInput{
		{
			Name:    "Parameter1",
			Min:     -1.3,
			Max:     1.3,
			Step:    1.3,
			NumSims: 3,
		},
		{
			Name:    "Parameter2",
			Min:     -1.3,
			Max:     1.3,
			Step:    1.3,
			NumSims: 3,
		},
	}
	actual := inputs.GenerateSimInputs(designParameters)
	if len(actual) != len(expected) {
		t.Errorf("Expected %d, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, got %v", expected[i], actual[i])
		}
	}
}

func TestGenerateStudyInputs(t *testing.T) {
	simInputs := []inputs.SimInput{
		{
			Name:    "Angles",
			Min:     -1.3,
			Max:     1.3,
			Step:    1.3,
			NumSims: 3,
		},
		{
			Name:    "Heaves",
			Min:     -1.69,
			Max:     0.31,
			Step:    1.0,
			NumSims: 3,
		},
	}

	expected := inputs.StudyInput{
		SimInputNames: []string{"Angles", "Heaves"},
		SimInputSamples: [][]float64{
			{-1.3, -1.69},
			{0.0, -1.69},
			{1.3, -1.69},
			{-1.3, -0.69},
			{0.0, -0.69},
			{1.3, -0.69},
			{-1.3, 0.31},
			{0.0, 0.31},
			{1.3, 0.31},
		},
	}
	actual := inputs.GenerateStudyInputs(simInputs, simInputs[0].NumSims)

	if len(actual.SimInputNames) != len(expected.SimInputNames) {
		t.Errorf("Expected %d, got %d", len(expected.SimInputNames), len(actual.SimInputNames))
	}
	for i := range expected.SimInputNames {
		if actual.SimInputNames[i] != expected.SimInputNames[i] {
			t.Errorf("Expected %s, got %s", expected.SimInputNames[i], actual.SimInputNames[i])
		}
	}
	if len(actual.SimInputSamples) != len(expected.SimInputSamples) {
		t.Errorf("Expected %d, got %d", len(expected.SimInputSamples), len(actual.SimInputSamples))
	}
	for i := range expected.SimInputSamples {
		if len(actual.SimInputSamples[i]) != len(expected.SimInputSamples[i]) {
			t.Errorf(
				"Expected %d, got %d",
				len(expected.SimInputSamples[i]),
				len(actual.SimInputSamples[i]),
			)
		}
	}
}

func TestGenerateSimInputCSV(t *testing.T) {
	testFileName := "../../../test/testoutput/testInputs.csv"
	expectedStudyInputs := inputs.StudyInput{
		SimInputNames: []string{"Angles", "Heaves"},
		SimInputSamples: [][]float64{
			{-1.3, -1.69},
			{0.0, -1.69},
			{1.3, -1.69},
			{-1.3, -0.69},
			{0.0, -0.69},
			{1.3, -0.69},
			{-1.3, 0.31},
			{0.0, 0.31},
			{1.3, 0.31},
		},
	}

	// remove old file
	if _, err := os.Stat(testFileName); err == nil {
		err := os.Remove(testFileName)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}

	err := inputs.GenerateSimInputCSV(expectedStudyInputs, testFileName)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	file, err := os.Open(testFileName)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(records)-1 != len(expectedStudyInputs.SimInputSamples) {
		t.Errorf("Expected %d, got %d", len(expectedStudyInputs.SimInputSamples), len(records)-1)
	}

	for i, record := range records {
		if i == 0 {
			for j, header := range record {
				if header != expectedStudyInputs.SimInputNames[j] {
					t.Errorf("Expected %s, got %s", expectedStudyInputs.SimInputNames[j], header)
				}
			}
		} else {
			for j, value := range record {
				// expectedValue := strconv.FormatFloat(expectedStudyInputs.SimInputSamples[i-1][j], 'f', -1, 64)
				expectedValue := expectedStudyInputs.SimInputSamples[i-1][j]
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					t.Errorf("Error parsing input file, %s", err)
				}
				if !math.AlmostEqual(floatValue, expectedValue) {
					t.Errorf("Expected %f, got %f", expectedValue, floatValue)
				}
			}
		}
	}
}
