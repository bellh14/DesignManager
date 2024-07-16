package designmanager

import (
	"fmt"
	"os"

	"github.com/bellh14/DesignManager/config"
	// "github.com/bellh14/DesignManager/pkg/optimization/pareto"
	// "github.com/bellh14/DesignManager/pkg/utils"
)

type DesignManager struct {
	ConfigFile config.ConfigFile
}

func NewDesignManager(config config.ConfigFile) *DesignManager {
	return &DesignManager{
		ConfigFile: config,
	}
}

func (designManager *DesignManager) Run() {
}

func (dm *DesignManager) HandleDesignStudy(studyType string) {
	switch studyType {
	case "Pareto":
		// dm.HandlePareto()
		fmt.Println("TODO: Implement Pareto")
	default:
		fmt.Println("Error: Study type not supported")
		os.Exit(1)
	}
}
