package jobs

import (
	"fmt"
	"github.com/bellh14/DesignManager/pkg/simulations"
	"github.com/bellh14/DesignManager/pkg/types"
	"sync"
)

func RunSimulation(jobSubmission *types.JobSubmissionType, simID int) types.SimulationResult {
	// Create simulation object
	// simulation := NewSimulation(jobSubmission, simID)
	// Run simulation
	// simulation.Run()
	// return simulation.SimulationResult
	return types.SimulationResult{}
}

func HandleSimulations(jobSubmission *types.JobSubmissionType, Results *[]types.SimulationResult, numSims int) {
	var wg sync.WaitGroup
	// results := make(chan types.SimulationResult, numSims)

	for i := 0; i < numSims; i++ {
		wg.Add(1)
		go func(simID int) {
			fmt.Println("Running simulation: ", simID)
			defer wg.Done()
			simulations.NewSimulation(jobSubmission, simID).Run() //TODO: fix this
		}(i)
		wg.Wait()
		// close(results)

		// for result := range results {
		// 	*Results = append(*Results, result)
		// }
	}
}
