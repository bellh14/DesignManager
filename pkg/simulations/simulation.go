package simulations

import (
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"github.com/bellh14/DFRDesignManager/pkg/jobs/generator"
)


type Simulation struct {
	JobNumber int
	JobSubmissionType types.JobSubmissionType
	InputParameters []types.SimInput
	DesignObjectiveResults []float64
}

func NewSimulation(jobSubmission *types.JobSubmissionType, simID int) *Simulation {
	return &Simulation{
		JobNumber: simID,
		JobSubmissionType: *jobSubmission,
	}
}

func (simulation *Simulation) Run() {
	// simulation.InputParameters = simulation.SampleDesignParameters()
	// simulation.DesignObjectiveResults = simulation.RunSimulation()
}

func (simulation *Simulation) SampleDesignParameters() []types.SimInput {
	// return []types.SimInput{}
	return []types.SimInput{}
}

func (simulation *Simulation) CreateSimulationDirectory() {
	// create directory
}

func (simulation *Simulation) CreateSimulationInputFile() {
	// create input file
}

func (simulation *Simulation) CreateJobScript() {
	generator.GenerateJobScript(simulation.JobSubmissionType, simulation.JobNumber)
}

func (simulation *Simulation) RunSimulation() {
	// exec job script
}

func (simulation *Simulation) ParseSimulationResults() {
	// parse results
}

func (simulation *Simulation) SaveResults() {
	// save results
}
