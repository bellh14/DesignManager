package designmanager

import (
	"fmt"
	"os"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/simulations"
	"github.com/bellh14/DesignManager/pkg/utils/log"
)

type DesignManager struct {
	ConfigFile config.ConfigFile
	Logger     *log.Logger
}

func NewDesignManager(config config.ConfigFile, logger *log.Logger) *DesignManager {
	return &DesignManager{
		ConfigFile: config,
		Logger:     logger,
	}
}

func (dm *DesignManager) Run() {
	if !dm.ConfigFile.UseDM {
		dm.Logger.Log("Use DM set to false. Exiting")
		return
	}
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
			dm.Logger.Error(fmt.Sprintf("Error Obtaining siminput for job number %d", i), err)
		}
		simLogger := log.NewLogger(0, fmt.Sprintf("Simulation: %d", i), "63")
		sim := simulations.NewSimulation(&jobSubmission, i, inputs, simLogger)
		sim.Run()
	}
	dm.Logger.Log("Finished running design sweep")
}

func (dm *DesignManager) HandleDesignStudy(studyType string) {
	switch studyType {
	case "Pareto":
		// dm.HandlePareto()
		fmt.Println("TODO: Implement Pareto")
	case "Sweep":
		dm.Logger.Log("Running design sweep")
		dm.HandleSweep()
	default:
		fmt.Println("Error: Study type not supported")
		os.Exit(1)
	}
}
