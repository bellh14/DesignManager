package jobs

import (
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"sync"
)

func RunSimulation(jobSubmission *types.JobSubmissionType, simID int) (types.SimulationResult) {
	// Create simulation object
	// simulation := NewSimulation(jobSubmission, simID)
	// Run simulation
	// simulation.Run()
	// return simulation.SimulationResult
	return types.SimulationResult{}
}

func HandleSimulations(Results *[]types.SimulationResult, numSims int) {
	var wg sync.WaitGroup
	results := make(chan types.SimulationResult, numSims)

	for i := 0; i < numSims; i++ {
		wg.Add(1)
		go func(simID int) {
			defer wg.Done()
			// simsParams := SampleDesignParameters()
			// simResult := RunSimulation(simsParams)
		}(i)
		wg.Wait()
		close(results)

		for result := range results {
			*Results = append(*Results, result)
		}	
	}
}