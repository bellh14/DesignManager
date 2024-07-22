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
	ConfigFile     config.ConfigFile
	Logger         *log.Logger
	InputGenerator inputs.SimInputGenerator
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
	dm.HandleInputs()
	dm.HandleDesignStudy(dm.ConfigFile.DesignStudyConfig.StudyType)
}

func (dm *DesignManager) HandleAeroMap() {
	// buffered channel 2nd param is for number of sweeps to run in parallel
	numberOfSweeps := dm.ConfigFile.DesignStudyConfig.DesignParameters[0].NumSims
	jobs := make(chan int, numberOfSweeps)

	// start sweeps
	for i := 0; i < numberOfSweeps; i++ {
		newDM := dm
		inputOffset := i * numberOfSweeps
		go newDM.HandleSweep(inputOffset, jobs)
	}

	// drain the channel
	for i := 0; i < numberOfSweeps; i++ {
		<-jobs // wait for task to complete
	}

	dm.Logger.Log("Finished Running AeroMap")
}

func (dm *DesignManager) HandleInputs() {
	dm.Logger.Log("Creating Input Parameter File")
	jobSubmission := jobscript.CreateJobSubmission(dm.ConfigFile)

	inputFileName := jobSubmission.WorkingDir + "/" + "Inputs.csv"

	dm.InputGenerator = *inputs.NewSimInputGenerator(
		dm.ConfigFile.DesignStudyConfig.DesignParameters,
		inputFileName,
		dm.ConfigFile.DesignStudyConfig.NumSims,
	)
	err := dm.InputGenerator.HandleSimInputs()
	if err != nil {
		dm.Logger.Error("Failed to HandleSimInputs", err)
	}
}

func (dm *DesignManager) HandleSweep(offset int, c chan int) {
	jobSubmission := jobscript.CreateJobSubmission(dm.ConfigFile)

	for i := 1; i <= dm.ConfigFile.DesignStudyConfig.NumSims; i++ {
		simNum := offset + i
		inputs, err := dm.InputGenerator.SimInputByJobNumber(simNum)
		if err != nil {
			fmt.Printf("Error obtaining siminput by job number %s", err)
			dm.Logger.Error(fmt.Sprintf("Error Obtaining siminput for job number %d", simNum), err)
		}
		simLogger := log.NewLogger(0, fmt.Sprintf("Simulation: %d", simNum), "63")
		sim := simulations.NewSimulation(&jobSubmission, simNum, inputs, simLogger)
		sim.Run()
	}
	dm.Logger.Log("Finished running design sweep")
	c <- 1 // signals sweep is finished
}

func (dm *DesignManager) HandleDesignStudy(studyType string) {
	switch studyType {
	case "AeroMap":
		dm.Logger.Log("Running AeroMap")
		dm.HandleAeroMap()
	case "Pareto":
		// dm.HandlePareto()
		fmt.Println("TODO: Implement Pareto")
	case "Sweep":
		dm.Logger.Log("Running design sweep")
		c := make(chan int, 1)
		dm.HandleSweep(0, c)
	default:
		fmt.Println("Error: Study type not supported")
		os.Exit(1)
	}
}
