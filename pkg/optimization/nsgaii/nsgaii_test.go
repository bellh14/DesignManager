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

func TestCompareCrowdingDistance(t *testing.T){
	nsgaii := &nsgaii.NSGAII{
		CurrentFront: &nsgaii.Population{
			Solutions: []*nsgaii.Solution{
				{
					DesignObjectives: []float64{1, 2},
					Rank: 1,
					Distance: 1,
				},
				{
					DesignObjectives: []float64{4, 5},
					Rank: 1,
					Distance: 2,
				},
				{
					DesignObjectives: []float64{7, 8},
					Rank: 2,
					Distance: 3,
				},
			},
		},
	}
	if !nsgaii.CompareCrowdingDistance(nsgaii.CurrentFront.Solutions[0], nsgaii.CurrentFront.Solutions[1]) {
		t.Errorf("Expected %v, got %v", true, nsgaii.CompareCrowdingDistance(nsgaii.CurrentFront.Solutions[0], nsgaii.CurrentFront.Solutions[1]))
	}

	if nsgaii.CompareCrowdingDistance(nsgaii.CurrentFront.Solutions[1], nsgaii.CurrentFront.Solutions[0]) {
		t.Errorf("Expected %v, got %v", false, nsgaii.CompareCrowdingDistance(nsgaii.CurrentFront.Solutions[1], nsgaii.CurrentFront.Solutions[0]))
	}
	
	if !nsgaii.CompareCrowdingDistance(nsgaii.CurrentFront.Solutions[1], nsgaii.CurrentFront.Solutions[2]) {
		t.Errorf("Expected %v, got %v", true, nsgaii.CompareCrowdingDistance(nsgaii.CurrentFront.Solutions[1], nsgaii.CurrentFront.Solutions[2]))
		t.Logf("Solution 1 rank: %v", nsgaii.CurrentFront.Solutions[1].Rank)
		t.Logf("Solution 2 rank: %v", nsgaii.CurrentFront.Solutions[2].Rank)
		t.Logf("Solution 1 distance: %v", nsgaii.CurrentFront.Solutions[1].Distance)
		t.Logf("Solution 2 distance: %v", nsgaii.CurrentFront.Solutions[2].Distance)
	}
}