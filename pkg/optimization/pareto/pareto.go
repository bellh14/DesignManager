package pareto

import (
	"github.com/bellh14/DFRDesignManager/pkg/types"
	// "github.com/bellh14/DFRDesignManager/pkg/utils"
)

type ParetoHandler struct {
	DesignManagerInputParameters *types.DesignManagerInputParameters
	JobSubmissionType            *types.JobSubmissionType
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
}