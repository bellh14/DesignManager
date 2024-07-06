package simulations

import (
	"fmt"
	"os"
	"os/exec"

	e "github.com/bellh14/DesignManager/pkg/err"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/utils"

)

type Simulation struct {
	JobNumber              int
	JobSubmission          jobscript.JobSubmission
	InputParameters        inputs.SimInputIteration
	JobDir                 string
	DesignObjectiveResults []float64
}

func NewSimulation(jobSubmission *jobscript.JobSubmission, simID int, inputParameters inputs.SimInputIteration) *Simulation {
	return &Simulation{
		JobNumber:       simID,
		JobSubmission:   *jobSubmission,
		InputParameters: inputParameters,
	}
}

func (simulation *Simulation) SetWorkingDir() {
	simulation.JobDir = simulation.JobSubmission.WorkingDir + "/" + fmt.Sprint(simulation.JobNumber)
}

func (simulation *Simulation) Run() {
	simulation.SetWorkingDir()
	simulation.CreateSimulationDirectory()
	simulation.CopySimulationFiles()
	fmt.Print(simulation.InputParameters)
	simulation.CreateSimulationInputFile()
	simulation.CreateJobScript()
	simulation.RunSimulation()
	// simulation.DesignObjectiveResults = simulation.RunSimulation()
}

// func (simulation *Simulation) SampleDesignParameters() []types.SimInput {
// 	sampler := sampling.NewSampler(simulation.JobSubmission)
// 	samples := sampler.Sample()
// 	return samples
// }

func (simulation *Simulation) SimulationInputs() {}

func (simulation *Simulation) CreateSimulationDirectory() {
	// create directory
	err := os.MkdirAll(simulation.JobDir, 0o777)
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
	}
}

func (simulation *Simulation) CopySimulationFiles() {
	// copy files
	fmt.Printf("Copying files to %s\n", simulation.JobDir)
	utils.CopyFile(simulation.JobSubmission.WorkingDir+"/"+simulation.JobSubmission.SimFile, simulation.JobDir+"/"+simulation.JobSubmission.SimFile)
	utils.CopyFile(simulation.JobSubmission.WorkingDir+"/"+simulation.JobSubmission.JavaMacro, simulation.JobDir+"/"+simulation.JobSubmission.JavaMacro)
}

func (simulation *Simulation) CreateSimulationInputFile() {
	// create input file
	inputFile, err := os.Create(simulation.JobDir + "/InputParams.csv")
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
	}
	utils.WriteParameterCsvHeader(simulation.InputParameters.Name, inputFile)
	utils.WriteSimulationInputCSV(simulation.InputParameters.Value, inputFile)
}

func (simulation *Simulation) CreateJobScript() {
	jobscript.GenerateJobScript(simulation.JobSubmission, simulation.JobNumber)
}

func (simulation *Simulation) RunSimulation() {
	// exec job script
	cmd := exec.Command(simulation.JobDir + "/job_" + fmt.Sprint(simulation.JobNumber) + ".sh")
	_, err := cmd.CombinedOutput()
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
	}
}

func (simulation *Simulation) ParseSimulationResults() {
	// parse results
}

func (simulation *Simulation) SaveResults() {
	// save results
}
