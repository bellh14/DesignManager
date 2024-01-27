package nsgaii

import (
	"math"
)

// implements the NSGA-II algorithm

type NSGAII struct {
	PopulationSize    int
	MaxGenerations    int
	CurrentGeneration int
	CurrentFront      *Population
	NextFront         *Population
	ObjectiveMaxes    []float64
	ObjectiveMins     []float64
}
type Population struct {
	Solutions []*Solution
}

type Solution struct {
	Rank               int
	DesignObjectives   []float64
	DominationCount    int
	DominatedSolutions []*Solution
	Distance           float64
}

func (nsgaii *NSGAII) Run() {

}

func (nsgaii *NSGAII) SortByObjective(objective int) (sortedSolutions []*Solution) {

	return sortedSolutions

}

func (nsgaii *NSGAII) FindObjectiveMaxes() {

	for i := 0; i < len(nsgaii.CurrentFront.Solutions[0].DesignObjectives); i++ {
		max := math.Inf(-1)
		for _, s := range nsgaii.CurrentFront.Solutions {
			if s.DesignObjectives[i] > max {
				max = s.DesignObjectives[i]
			}
		}
		nsgaii.ObjectiveMaxes[i] = max
	}

}

func (nsgaii *NSGAII) FindObjectiveMins() {

	for i := 0; i < len(nsgaii.CurrentFront.Solutions[0].DesignObjectives); i++ {
		min := math.Inf(1)
		for _, s := range nsgaii.CurrentFront.Solutions {
			if s.DesignObjectives[i] < min {
				min = s.DesignObjectives[i]
			}
		}
		nsgaii.ObjectiveMins[i] = min
	}

}

func (nsgaii *NSGAII) CalculateCrowdingDistance() {

	for _, s := range nsgaii.CurrentFront.Solutions {
		s.Distance = 0
	}

	for i := 0; i < len(nsgaii.CurrentFront.Solutions[0].DesignObjectives); i++ {
		// nsgaii.CurrentFront.Solutions = nsgaii.SortByObjective(i)
		nsgaii.CurrentFront.Solutions[0].Distance = 1000000000
		nsgaii.CurrentFront.Solutions[len(nsgaii.CurrentFront.Solutions)-1].Distance = 1000000000
		for j := 1; j < len(nsgaii.CurrentFront.Solutions)-1; j++ {
			if i == 0 {
				//total df so want higher values
				nsgaii.CurrentFront.Solutions[j].Distance += (nsgaii.CurrentFront.Solutions[j+1].DesignObjectives[i] - nsgaii.CurrentFront.Solutions[j-1].DesignObjectives[i]) / math.Abs(nsgaii.ObjectiveMins[i]-nsgaii.ObjectiveMaxes[i])
			} else {
				// drag so want lower values
				nsgaii.CurrentFront.Solutions[j].Distance += (nsgaii.CurrentFront.Solutions[j+1].DesignObjectives[i] - nsgaii.CurrentFront.Solutions[j-1].DesignObjectives[i]) / (nsgaii.ObjectiveMaxes[i] - nsgaii.ObjectiveMins[i])

			}
		}
	}

}

func (nsgaii *NSGAII) CompareCrowdingDistance(s1 *Solution, s2 *Solution) bool {

	if s1.Rank < s2.Rank || (s1.Rank == s2.Rank && s1.Distance > s2.Distance) {
		return false
	}
	return true

}

func (nsgaii *NSGAII) InitalizePopulation() {

}

func (nsgaii *NSGAII) RankSolutions() {
	for _, s1 := range nsgaii.CurrentFront.Solutions {

		s1.DominationCount = 0
		s1.DominatedSolutions = make([]*Solution, 0)

		for _, s2 := range nsgaii.CurrentFront.Solutions {
			if Dominates(s1, s2) {
				s1.DominatedSolutions = append(s1.DominatedSolutions, s2)
			} else if Dominates(s2, s1) {
				s1.DominationCount++
			}
		}
		if s1.DominationCount == 0 {
			s1.Rank = 1
			nsgaii.NextFront.Solutions = append(nsgaii.NextFront.Solutions, s1)
		}
	}
	for i := 1; len(nsgaii.NextFront.Solutions) > 0; i++ {
		var nextFront []*Solution
		for _, s1 := range nsgaii.NextFront.Solutions {
			for _, s2 := range s1.DominatedSolutions {
				s2.DominationCount--
				if s2.DominationCount == 0 {
					s2.Rank = i + 1
					nextFront = append(nextFront, s2)
				}
			}
		}
		nsgaii.NextFront.Solutions = nextFront
	}

}

func Dominates(solution1 *Solution, solution2 *Solution) bool {
	// for i := range solution1.DesignObjectives {
	// 	if solution1.DesignObjectives[i] > solution2.DesignObjectives[i] {
	// 		return false
	// 	}
	// }
	if solution1.DesignObjectives[0] < solution2.DesignObjectives[0] {
		// for now means solution 1 has a lower total df
		return false
	}
	if solution1.DesignObjectives[1] < solution2.DesignObjectives[1] {
		// for now means solution 1 has higher df and lower drag
		return true
	}
	return false
}
