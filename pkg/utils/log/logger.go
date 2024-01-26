package log

import (
	"fmt"
	"log"
	"os"

	"github.com/bellh14/DFRDesignManager/pkg/types"
	"github.com/bellh14/DFRDesignManager/pkg/utils/err"
)

type Logger struct {
	*log.Logger
}

func Setup(logFile string){
	// Create log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func (logger *Logger) Info(message string) {
	logger.Output(2, message)
}

func (logger *Logger) Error(message string) {
	logger.Output(2, message)
}

func (logger *Logger) SimError (jobNumber int, e error) {
	simError := &err.SimulationError{
		JobNumber: jobNumber,
		Err: e,
	}
	logger.Output(2, simError.SimError())
}

func (logger *Logger) GenerationResults (generationResults types.GenerationResults){
	logger.Output(2, fmt.Sprintf("Generation Results: %v", generationResults))
	if len(generationResults.FailedSims) > 0 {
		for _, sim := range generationResults.FailedSims {
			logger.Output(2, fmt.Sprintf("Simulation: %v failed from %s", sim.JobNumber, sim.Cause))
		}
	}

	for sim := range generationResults.SucceededSims {
		logger.Output(2, fmt.Sprintf("Succedded Simulation: %v", sim))
	}
}