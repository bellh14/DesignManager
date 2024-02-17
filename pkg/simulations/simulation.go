package simulations

import (
	"fmt"
	"github.com/bellh14/DFRDesignManager/pkg/jobs/generator"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"github.com/bellh14/DFRDesignManager/pkg/utils"
	e "github.com/bellh14/DFRDesignManager/pkg/utils/err"
	"github.com/bellh14/DFRDesignManager/pkg/utils/math/sampling"
	"os"
	"os/exec"
)

type Simulation struct {
	JobNumber              int
	JobSubmissionType      types.JobSubmissionType
	InputParameters        []types.SimInput
	DesignObjectiveResults []float64
}

func NewSimulation(jobSubmission *types.JobSubmissionType, simID int) *Simulation {
	return &Simulation{
		JobNumber:         simID,
		JobSubmissionType: *jobSubmission,
	}
}

func (simulation *Simulation) Run() {
	simulation.CreateSimulationDirectory()
	simulation.CopySimulationFiles()
	simulation.InputParameters = simulation.SampleDesignParameters()
	simulation.CreateSimulationInputFile()
	simulation.CreateJobScript()

	// simulation.DesignObjectiveResults = simulation.RunSimulation()
}

func (simulation *Simulation) SampleDesignParameters() []types.SimInput {
	sampler := sampling.NewSampler(simulation.JobSubmissionType)
	samples := sampler.Sample()
	return samples
}

func (simulation *Simulation) CreateSimulationDirectory() {
	// create directory
	err := os.MkdirAll(simulation.JobSubmissionType.WorkingDir, 0777)
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
	}
}

func (simulation *Simulation) CopySimulationFiles() {
	// copy files
	utils.CopyFile(simulation.JobSubmissionType.SimFile, simulation.JobSubmissionType.WorkingDir)
	utils.CopyFile(simulation.JobSubmissionType.JavaMacro, simulation.JobSubmissionType.WorkingDir)
}

func (simulation *Simulation) CreateSimulationInputFile() {
	// create input file
	inputFile, err := os.Create(simulation.JobSubmissionType.WorkingDir + "/input.csv")
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
	}
	utils.WriteParameterCsvHeader(simulation.InputParameters, inputFile)
	utils.WriteParameterCsv(simulation.InputParameters, inputFile)
}

func (simulation *Simulation) CreateJobScript() {
	generator.GenerateJobScript(simulation.JobSubmissionType, simulation.JobNumber)
}

func (simulation *Simulation) RunSimulation() {
	// exec job script
	cmd := exec.Command(simulation.JobSubmissionType.WorkingDir + "/job_" + fmt.Sprint(simulation.JobNumber) + ".sh")
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
