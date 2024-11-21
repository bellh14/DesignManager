package custom_test

import (
	"math"
	"testing"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/optimization/custom"
	"github.com/bellh14/DesignManager/pkg/optimization/genetic"
	"github.com/bellh14/DesignManager/pkg/simulations"
)

func TestCalculateObjectivesPeaks(t *testing.T) {
	population := genetic.Population{
		genetic.Individual{
			Sim: &simulations.Simulation{
				Successful: true,
				DesignObjectiveResults: map[string]float64{
					"Aero Efficiency":     2.2,
					"Downforce":           120,
					"Rear Axle Downforce": 80,
				},
			},
			Fitness: 0.0,
		},
		genetic.Individual{
			Sim: &simulations.Simulation{
				Successful: true,
				DesignObjectiveResults: map[string]float64{
					"Aero Efficiency":     2.5,
					"Downforce":           110,
					"Rear Axle Downforce": 70,
				},
			},
			Fitness: 0.0,
		},
		genetic.Individual{
			Sim: &simulations.Simulation{
				Successful: true,
				DesignObjectiveResults: map[string]float64{
					"Aero Efficiency":     2.0,
					"Downforce":           140,
					"Rear Axle Downforce": 90,
				},
			},
			Fitness: 0.0,
		},
	}
	minValues, maxValues := custom.CalculateObjectivesPeaks(population)

	if minValues["Aero Efficiency"] != 2.0 {
		t.Errorf(
			"Expected min value of Aero Efficiency to be 2.0, got %f",
			minValues["Aero Efficiency"],
		)
	}
	if maxValues["Aero Efficiency"] != 2.5 {
		t.Errorf(
			"Expected max value of Aero Efficiency to be 2.5, got %f",
			maxValues["Aero Efficiency"],
		)
	}

	if minValues["Downforce"] != 110 {
		t.Errorf("Expected min value of Downforce to be 110, got %f", minValues["Downforce"])
	}
	if maxValues["Downforce"] != 140 {
		t.Errorf("Expected max value of Downforce to be 140, got %f", maxValues["Downforce"])
	}

	if minValues["Rear Axle Downforce"] != 70 {
		t.Errorf(
			"Expected min value of Rear Axle Downforce to be 70, got %f",
			minValues["Rear Axle Downforce"],
		)
	}
	if maxValues["Rear Axle Downforce"] != 90 {
		t.Errorf(
			"Expected max value of Rear Axle Downforce to be 90, got %f",
			maxValues["Rear Axle Downforce"],
		)
	}
}

func TestNormalize(t *testing.T) {
	normalized := custom.Normalize(120, 110, 140, "Maximize")
	if math.Abs(normalized-0.333333) > 0.0001 {
		t.Errorf("Expected normalized value to be 0.333333, got %f", normalized)
	}

	normalizedMinimize := custom.Normalize(120, 110, 140, "Minimize")
	if math.Abs(normalizedMinimize-0.666666) > 0.0001 {
		t.Errorf("Expected normalized value to be 0.666666, got %f", normalizedMinimize)
	}
}

func TestCalculateFitness(t *testing.T) {
	individual := genetic.Individual{
		Sim: &simulations.Simulation{
			DesignObjectiveResults: map[string]float64{
				"Aero Efficiency":      2.2,
				"Front Axle Downforce": 88,
				"Rear Axle Downforce":  78,
			},
		},
		Fitness: 0.0,
	}
	dsc := config.DesignStudyConfig{
		DesignObjectives: []config.DesignObjective{
			{
				Name:   "Aero Efficiency",
				Target: 0.0,
				Goal:   "Maximize",
				Weight: 0.75,
			},
			{
				Name:   "Front Axle Downforce",
				Target: 85,
				Goal:   "Maximize",
				Weight: 0.75,
			},
			{
				Name:   "Rear Axle Downforce",
				Target: 85,
				Goal:   "Maximize",
				Weight: 1,
			},
		},
	}
	minValues := map[string]float64{
		"Aero Efficiency":      2.0,
		"Front Axle Downforce": 78,
		"Rear Axle Downforce":  60,
	}
	maxValues := map[string]float64{
		"Aero Efficiency":      2.5,
		"Front Axle Downforce": 94,
		"Rear Axle Downforce":  90,
	}
	custom.CalculateFitness(&individual, dsc, minValues, maxValues)

	if math.Abs(individual.Fitness-0.698750) > 0.0001 {
		t.Errorf("Expected fitness value to be 0.698750, got %f", individual.Fitness)
	}
}
