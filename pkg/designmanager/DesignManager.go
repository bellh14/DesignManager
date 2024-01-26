package designmanager

import (
	"fmt"
	"os"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"github.com/bellh14/DFRDesignManager/pkg/optimization/pareto"
	"github.com/bellh14/DFRDesignManager/pkg/utils"
)

type DesignManager struct {
	ConfigFile    types.ConfigFile
}

func NewDesignManager(config types.ConfigFile) (*DesignManager) {
	return &DesignManager{
		ConfigFile: config,
	}
}

func (designManager *DesignManager) Run() {
}

func (dm *DesignManager) handleDesignStudy(studyType string) {
	switch studyType {
	case "Pareto":
		dm.handlePareto()
	default:
		fmt.Println("Error: Study type not supported")
		os.Exit(1)
	}
}

func (dm *DesignManager) handlePareto() {

	jobSubmission := utils.CreateJobSubmission(dm.ConfigFile.SystemResources, dm.ConfigFile.WorkingDir, dm.ConfigFile.StarCCM, 0)

	// Create pareto object
	paretoHandler := pareto.NewPareto(dm.ConfigFile.DesignManagerInputParameters, jobSubmission)

	// Run pareto
	paretoHandler.Run()
}