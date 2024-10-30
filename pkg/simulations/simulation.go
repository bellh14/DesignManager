package simulations

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
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
) *Simulation {
	return &Simulation{
		JobNumber:       simID,
		JobSubmission:   *jobSubmission,
		InputParameters: inputParameters,
		Logger:          logger,
		SlurmConfig:     slurmConfig,
		HostNodes:       hostNodes,
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
	time.Sleep(time.Second * time.Duration(simulation.JobNumber) / 2)
	simulation.Logger.LogSimulation(simulation.LogValue(), "Running Simulation")
	simulation.SetWorkingDir()
	simulation.CreateSimulationDirectory()
	simulation.CopySimulationFiles()
	simulation.CreateSimulationInputFile()
	simulation.CreateSimulationMachineFile()
	simulation.CreateJobScript()
	simulation.RunSimulation()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	simulation.Logger.LogSimulation(
		simulation.LogValue(),
		fmt.Sprintf("Alloc = %v MiB\n", m.Alloc/1024/1024),
	)
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
	time.Sleep(time.Second)
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
		simulation.MachineFile,
	)
	time.Sleep(time.Second)
}

func (simulation *Simulation) RunSimulation() {
	// exec job script
	simulation.Logger.LogSimulation(simulation.LogValue(), "Starting StarCCM+")

	cmd := exec.Command(simulation.JobDir + "/sim_" + fmt.Sprint(simulation.JobNumber) + ".sh")
	// cmd := exec.Command("sbatch", simulation.JobDir+"sim_"+fmt.Sprint(simulation.JobNumber)+".sh")
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
	fmt.Println(simulation.DesignObjectiveResults)
	return parameterNames, floatResults
}

func (simulation *Simulation) SaveResults() {
	// save results
}
