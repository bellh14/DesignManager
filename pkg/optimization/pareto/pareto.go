package pareto

import (
	"github.com/bellh14/DesignManager/pkg/types"
	// "github.com/bellh14/DesignManager/pkg/utils"
	// "github.com/bellh14/DesignManager/pkg/jobs"
)

type ParetoHandler struct {
	DesignManagerInputParameters *types.DesignManagerInputParameters
	JobSubmissionType            *types.JobSubmissionType
	Results *[]types.SimulationResult
}

func NewPareto(designManagerInputParams types.DesignManagerInputParameters, jobSubmissionParams types.JobSubmissionType) *ParetoHandler {
	return &ParetoHandler{
		DesignManagerInputParameters: &designManagerInputParams,
		JobSubmissionType:            &jobSubmissionParams,
	}
}

func (paretoHandler *ParetoHandler) Run() {
	// utils.PrintStruct(paretoHandler.DesignManagerInputParameters)
	// utils.PrintStruct(paretoHandler.JobSubmissionType)
	// jobs.HandleSimulations(paretoHandler.JobSubmissionType, paretoHandler.Results, 64)
}