package nsgaii_test

import (
	"testing"
	"github.com/bellh14/DFRDesignManager/pkg/optimization/nsgaii"
)

func TestFindMaxes(t *testing.T) {
	nsgaii := &nsgaii.NSGAII{
		CurrentFront: &nsgaii.Population{
			Solutions: []*nsgaii.Solution{
				{
					DesignObjectives: []float64{1, 2},
				},
				{
					DesignObjectives: []float64{4, 5},
				},
				{
					DesignObjectives: []float64{7, 8},
				},
			},
		},
		ObjectiveMaxes: make([]float64, 2),
	}
	nsgaii.FindObjectiveMaxes()
	for i, max := range nsgaii.ObjectiveMaxes {
		if max != float64(i + 7) {
			t.Errorf("Expected %v, got %v", float64(i + 7), max)
		}
	}
}

func TestFindMins(t *testing.T) {
	nsgaii := &nsgaii.NSGAII{
		CurrentFront: &nsgaii.Population{
			Solutions: []*nsgaii.Solution{
				{
					DesignObjectives: []float64{1, 2},
				},
				{
					DesignObjectives: []float64{4, 5},
				},
				{
					DesignObjectives: []float64{7, 8},
				},
			},
		},
		ObjectiveMins: make([]float64, 2),
	}
	nsgaii.FindObjectiveMins()
	for i, min := range nsgaii.ObjectiveMins {
		if min != float64(i + 1) {
			t.Errorf("Expected %v, got %v", float64(i + 1), min)
		}
	}
}

func TestCalculateCrowdingDistance(t *testing.T) {
	nsgaii := &nsgaii.NSGAII{
		CurrentFront: &nsgaii.Population{
			Solutions: []*nsgaii.Solution{
				{
					DesignObjectives: []float64{1, 2},
				},
				{
					DesignObjectives: []float64{4, 5},
				},
				{
					DesignObjectives: []float64{7, 8},
				},
			},
		},
		ObjectiveMaxes: make([]float64, 2),
		ObjectiveMins: make([]float64, 2),
	}
	nsgaii.FindObjectiveMins()
	nsgaii.FindObjectiveMaxes()
	nsgaii.CalculateCrowdingDistance()
	for _, s := range nsgaii.CurrentFront.Solutions {
		if s.Distance != 0 {
			t.Logf("Expected %v, got %v", 0, s.Distance)
		}
	}
}