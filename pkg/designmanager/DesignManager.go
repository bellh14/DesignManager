package designmanager

import (
	"fmt"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"os"
	// "github.com/bellh14/DFRDesignManager/pkg/optimization/pareto"
	"github.com/bellh14/DFRDesignManager/pkg/jobs"
	"github.com/bellh14/DFRDesignManager/pkg/utils"
)

type DesignManager struct {
	ConfigFile types.ConfigFile
}

func NewDesignManager(config types.ConfigFile) *DesignManager {
	return &DesignManager{
		ConfigFile: config,
	}
}

func (designManager *DesignManager) Run() {
}

func (dm *DesignManager) HandleDesignStudy(studyType string) {
	switch studyType {
	case "Pareto":
		dm.HandlePareto()
	default:
		fmt.Println("Error: Study type not supported")
		os.Exit(1)
	}
}

func (dm *DesignManager) HandlePareto() {

	jobSubmission := utils.CreateJobSubmission(dm.ConfigFile.SystemResources, dm.ConfigFile.WorkingDir, dm.ConfigFile.StarCCM)

	// Create pareto object
	// paretoHandler := pareto.NewPareto(dm.ConfigFile.DesignManagerInputParameters, jobSubmission)

	// // Run pareto
	// paretoHandler.Run()

	results := make([]types.SimulationResult, 0)
	for i := 0; i < 8; i++ {
		jobs.HandleSimulations(&jobSubmission, &results, 52)
	}

}
