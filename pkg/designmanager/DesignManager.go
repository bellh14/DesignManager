package designmanager

import (
	"fmt"
	"os"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/simulations"
)

type DesignManager struct {
	ConfigFile config.ConfigFile
}

func NewDesignManager(config config.ConfigFile) *DesignManager {
	return &DesignManager{
		ConfigFile: config,
	}
}

func (dm *DesignManager) Run() {
	dm.HandleDesignStudy(dm.ConfigFile.DesignStudyConfig.StudyType)
}

func (dm *DesignManager) HandleSweep() {
	jobSubmission := jobscript.CreateJobSubmission(dm.ConfigFile)

	inputFileName := jobSubmission.WorkingDir + "/" + "Inputs.csv"

	inputGenerator := inputs.NewSimInputGenerator(
		dm.ConfigFile.DesignStudyConfig.DesignParameters,
		inputFileName,
	)
	inputGenerator.HandleSimInputs()

	for i := 1; i <= dm.ConfigFile.DesignStudyConfig.NumSims; i++ {
		inputs, err := inputGenerator.SimInputByJobNumber(i)
		if err != nil {
			fmt.Printf("Error obtaining siminput by job number %s", err)
		}
		sim := simulations.NewSimulation(&jobSubmission, i, inputs)
		sim.Run()
	}
	fmt.Println("Finished running design sweep")
}

func (dm *DesignManager) HandleDesignStudy(studyType string) {
	switch studyType {
	case "Pareto":
		// dm.HandlePareto()
		fmt.Println("TODO: Implement Pareto")
	case "Sweep":
		fmt.Println("Running Design Sweep")
		dm.HandleSweep()
	default:
		fmt.Println("Error: Study type not supported")
		os.Exit(1)
	}
}
