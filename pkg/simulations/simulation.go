package simulations

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"

	e "github.com/bellh14/DesignManager/pkg/err"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/utils"
	"github.com/bellh14/DesignManager/pkg/utils/log"
)

type Simulation struct {
	JobNumber              int
	JobSubmission          jobscript.JobSubmission
	InputParameters        inputs.SimInputIteration
	JobDir                 string
	DesignObjectiveResults map[string]float64
	Logger                 *log.Logger
}

func LogSimParameters(inputParameters inputs.SimInputIteration) string {
	parameterSlice := make([]string, len(inputParameters.Name))

	for i := range inputParameters.Name {
		parameterSlice[i] = fmt.Sprintf(
			"%s: %.4f",
			inputParameters.Name[i],
			inputParameters.Value[i],
		)
	}
	return strings.Join(parameterSlice, ", ")
}

func (sim *Simulation) LogValue() slog.Value {
	parameters := LogSimParameters(sim.InputParameters)
	return slog.GroupValue(
		slog.String("JobNum:", fmt.Sprint(sim.JobNumber)),
		slog.String("Inputs", parameters),
	)
}

func NewSimulation(
	jobSubmission *jobscript.JobSubmission,
	simID int,
	inputParameters inputs.SimInputIteration,
	logger *log.Logger,
) *Simulation {
	return &Simulation{
		JobNumber:       simID,
		JobSubmission:   *jobSubmission,
		InputParameters: inputParameters,
		Logger:          logger,
	}
}

func (simulation *Simulation) SetWorkingDir() {
	simulation.JobDir = simulation.JobSubmission.WorkingDir + "/" + fmt.Sprint(simulation.JobNumber)
	simulation.Logger.LogSimulation(
		simulation.LogValue(),
		fmt.Sprintf("Setting Working Directory %s", simulation.JobDir),
	)
}

func (simulation *Simulation) Run() {
	simulation.Logger.LogSimulation(simulation.LogValue(), "Running Simulation")
	simulation.SetWorkingDir()
	simulation.CreateSimulationDirectory()
	simulation.CopySimulationFiles()
	simulation.CreateSimulationInputFile()
	simulation.CreateJobScript()
	simulation.RunSimulation()
	simulation.Logger.LogSimulation(simulation.LogValue(), "Finished running simulation\n\n")
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
	simulation.Logger.LogSimulation(simulation.LogValue(), "Creating simulation directory")
	err := os.MkdirAll(simulation.JobDir, 0o777)
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simulation.Logger.Error("Error Creating simulation directory", err)
		simError.SimError()
	}
}

func (simulation *Simulation) CopySimulationFiles() {
	// copy files
	simulation.Logger.LogSimulation(
		simulation.LogValue(),
		fmt.Sprintf("Copying files to %s\n", simulation.JobDir),
	)
	utils.CopyFile(
		simulation.JobSubmission.WorkingDir+"/"+simulation.JobSubmission.SimFile,
		simulation.JobDir+"/"+simulation.JobSubmission.SimFile,
	)
	utils.CopyFile(
		simulation.JobSubmission.WorkingDir+"/"+simulation.JobSubmission.JavaMacro,
		simulation.JobDir+"/"+simulation.JobSubmission.JavaMacro,
	)
}

func (simulation *Simulation) CreateSimulationInputFile() {
	// create input file
	simulation.Logger.LogSimulation(simulation.LogValue(), "Creating Input CSV")
	inputFile, err := os.Create(simulation.JobDir + "/InputParams.csv")
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
		simulation.Logger.Error("Failed to created input csv", err)
	}
	defer inputFile.Close()
	utils.WriteParameterCsvHeader(simulation.InputParameters.Name, inputFile)
	utils.WriteSimulationInputCSV(simulation.InputParameters.Value, inputFile)
	inputFile.Close()
}

func (simulation *Simulation) CreateJobScript() {
	simulation.Logger.LogSimulation(simulation.LogValue(), "Creating Job Script")
	jobscript.GenerateJobScript(simulation.JobSubmission, simulation.JobNumber)
}

func (simulation *Simulation) RunSimulation() {
	// exec job script
	simulation.Logger.LogSimulation(simulation.LogValue(), "Starting StarCCM+")
	cmd := exec.Command(simulation.JobDir + "/sim_" + fmt.Sprint(simulation.JobNumber) + ".sh")
	_, err := cmd.CombinedOutput()
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
		fmt.Printf(simError.SimError() + "\n")
		simulation.Logger.Error(simError.SimError(), err)
	}
}

func (simulation *Simulation) ParseSimulationResults() ([]string, []float64) {
	// parse results
	simName := strings.TrimSuffix(simulation.JobSubmission.SimFile, ".sim")
	reportName := simulation.JobDir + "/" + simName + "_Report.csv"
	file, err := os.Open(reportName)
	if err != nil {
		simulation.Logger.Error("Failed to parse simulation results", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	parameterNames, err := csvReader.Read()
	if err != nil {
		simulation.Logger.Error("Failed to read report header", err)
		return nil, nil
	}

	parameterResults, err := csvReader.Read()
	if err != nil {
		simulation.Logger.Error("Failed to read report values", err)
		return nil, nil
	}

	for i, parameterName := range parameterNames {
		if _, exists := simulation.DesignObjectiveResults[parameterName]; exists {
			result, err := strconv.ParseFloat(parameterResults[i], 64)
			if err != nil {
				simulation.Logger.Error("Error parsing float value", err)
				continue
			}
			simulation.DesignObjectiveResults[parameterName] = result
		}
	}
	floatResults, err := utils.ConvertStringSliceToFloat(parameterResults)
	if err != nil {
		simulation.Logger.Error("Failed to parse results into float slice", err)
		return nil, nil
	}
	return parameterNames, floatResults
}

func (simulation *Simulation) SaveResults() {
	// save results
}
