package simulations

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	e "github.com/bellh14/DesignManager/pkg/err"
	"github.com/bellh14/DesignManager/pkg/generator/batchsystem"
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
	SlurmConfig            batchsystem.SlurmConfig
	MachineFile            string
	HostNodes              string
	TestFunction           string
	Successful             bool
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
	slurmConfig batchsystem.SlurmConfig,
	hostNodes string,
	testFunction string,
) *Simulation {
	return &Simulation{
		JobNumber:       simID,
		JobSubmission:   *jobSubmission,
		InputParameters: inputParameters,
		Logger:          logger,
		SlurmConfig:     slurmConfig,
		HostNodes:       hostNodes,
		TestFunction:    testFunction,
	}
}

func (simulation *Simulation) SetWorkingDir() {
	simulation.JobDir = simulation.JobSubmission.WorkingDir + "/" + fmt.Sprint(simulation.JobNumber)
}

func (simulation *Simulation) Run() {
	time.Sleep(time.Second)
	simulation.Logger.LogSimulation(simulation.LogValue(), "Running Simulation")
	simulation.SetWorkingDir()
	simulation.CreateSimulationDirectory()
	simulation.CopySimulationFiles()
	simulation.CreateSimulationInputFile()
	simulation.CreateSimulationMachineFile()
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
	err := os.MkdirAll(simulation.JobDir, 0o777)
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simulation.Logger.Error("Error Creating simulation directory", err)
		simError.SimError()
	}
	time.Sleep(time.Second)
}

func (simulation *Simulation) CopySimulationFiles() {
	utils.CopyFile(
		simulation.JobSubmission.WorkingDir+"/"+simulation.JobSubmission.SimFile,
		simulation.JobDir+"/"+simulation.JobSubmission.SimFile,
	)
	utils.CopyFile(
		simulation.JobSubmission.WorkingDir+"/"+simulation.JobSubmission.JavaMacro,
		simulation.JobDir+"/"+simulation.JobSubmission.JavaMacro,
	)
	if simulation.TestFunction != "" {
		utils.CopyFile("MOOT", simulation.JobDir+"/MOOT")
		_ = os.Chmod(fmt.Sprintf("%s/MOOT", simulation.JobDir), 0o777)
	}
	time.Sleep(time.Second)
}

func (simulation *Simulation) CreateSimulationInputFile() {
	// create input file
	simulation.Logger.LogSimulation(simulation.LogValue(), "Creating Input CSV")
	inputFile, err := os.Create(simulation.JobDir + "/InputParams.csv")
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
		simulation.Logger.Error("Failed to create input csv", err)
	}
	defer inputFile.Close()
	utils.WriteParameterCsvHeader(simulation.InputParameters.Name, inputFile)
	utils.WriteSimulationInputCSV(simulation.InputParameters.Value, inputFile)
	inputFile.Close()
	time.Sleep(time.Second)
}

func (simulation *Simulation) CreateSimulationMachineFile() {
	simulation.Logger.LogSimulation(simulation.LogValue(), "Creating machinefile")
	simulation.MachineFile = fmt.Sprintf("%d.txt", simulation.JobNumber)
	err := jobscript.CreateMachineFile(
		fmt.Sprintf("%s/%s", simulation.JobDir, simulation.MachineFile),
		simulation.HostNodes,
		simulation.JobSubmission.NtasksPerNode,
	)
	if err != nil {
		simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
		simError.SimError()
		simulation.Logger.Error("Failed to create machine file", err)
	}
	time.Sleep(time.Second)
}

func (simulation *Simulation) CreateJobScript() {
	simulation.Logger.LogSimulation(simulation.LogValue(), "Creating Job Script")
	jobscript.GenerateJobScript(
		simulation.JobSubmission,
		simulation.JobNumber,
		simulation.InputParameters,
		simulation.SlurmConfig,
		simulation.HostNodes,
		simulation.TestFunction,
	)
	time.Sleep(time.Second)
}

func (simulation *Simulation) RunSimulation() {
	// exec job script
	simulation.Logger.LogSimulation(simulation.LogValue(), "Starting StarCCM+")

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		simulation.JobDir+"/sim_"+fmt.Sprint(simulation.JobNumber)+".sh",
	)
	// cmd := exec.Command("sbatch", simulation.JobDir+"sim_"+fmt.Sprint(simulation.JobNumber)+".sh")
	errChan := make(chan error, 1)
	go func() {
		err := cmd.Run()
		errChan <- err
	}()

	select {
	case err := <-errChan:
		if err != nil {
			simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
			simError.SimError()
			simulation.Logger.Error(simError.SimError(), err)
			simulation.Successful = false
		}
		simulation.Successful = true
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			simulation.Logger.LogSimulation(
				simulation.LogValue(),
				"1 hour timeout reaced. Killing simulation.",
			)
			if killErr := cmd.Process.Kill(); killErr != nil {
				simulation.Logger.Error("Failed to kill simulation", killErr)
			}
		}
	}
	// _, err := cmd.CombinedOutput()
	// if err != nil {
	// 	simError := e.SimulationError{JobNumber: simulation.JobNumber, Err: err}
	// 	simError.SimError()
	// 	fmt.Printf(simError.SimError() + "\n")
	// 	simulation.Logger.Error(simError.SimError(), err)
	// }
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
	fmt.Println(simulation.DesignObjectiveResults)
	return parameterNames, floatResults
}

func (simulation *Simulation) SaveResults() {
	// save results
}
