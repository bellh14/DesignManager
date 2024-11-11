package designmanager

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/discord"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/optimization/custom"
	"github.com/bellh14/DesignManager/pkg/optimization/genetic"
	"github.com/bellh14/DesignManager/pkg/simulations"
	"github.com/bellh14/DesignManager/pkg/utils"
	"github.com/bellh14/DesignManager/pkg/utils/log"
)

type DesignManager struct {
	ConfigFile      config.ConfigFile
	Logger          *log.Logger
	InputGenerator  inputs.SimInputGenerator
	SimResultParams []string
	SimResults      [][]float64
	DiscordHook     discord.DiscordHook
}

func NewDesignManager(config config.ConfigFile, logger *log.Logger) *DesignManager {
	return &DesignManager{
		ConfigFile:      config,
		Logger:          logger,
		SimResultParams: make([]string, config.DesignStudyConfig.NumSims),
		SimResults:      make([][]float64, config.DesignStudyConfig.NumSims),
	}
}

func (dm *DesignManager) Run() {
	if !dm.ConfigFile.UseDM {
		dm.Logger.Log("Use DM set to false. Exiting")
		return
	}
	if dm.ConfigFile.DesignStudyConfig.StudyType != "Pareto" {
		dm.HandleInputs()
	}
	dm.HandleDesignStudy(dm.ConfigFile.DesignStudyConfig.StudyType)
	dm.SaveCompiledResults("")
	dm.DiscordHook.PayloadJson.Content = "Finished running design study"
	dm.DiscordHook.CallWebHook()
}

func (dm *DesignManager) HandleAeroMap() {
	numberOfSweeps := dm.ConfigFile.DesignStudyConfig.DesignParameters[0].NumSims
	offset := dm.ConfigFile.DesignStudyConfig.DesignParameters[1].NumSims

	// buffered channel 2nd param is for number of sweeps to run in parallel
	jobs := make(chan int, numberOfSweeps)
	wg := sync.WaitGroup{}
	var mu sync.Mutex

	// start sweeps
	for i := range numberOfSweeps {
		wg.Add(1)
		jobs <- 1
		newDM := dm
		go func(i int) {
			inputOffset := i * offset
			newDM.HandleSweep(inputOffset, offset, i)
			<-jobs
			mu.Lock()
			dm.SimResultParams = newDM.SimResultParams
			dm.SimResults = append(dm.SimResults, newDM.SimResults...)
			mu.Unlock()
			defer wg.Done()
		}(i)

	}
	wg.Wait()

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

func (dm *DesignManager) HandleSweep(offset int, numSims int, hostIndex int) {
	jobSubmission := jobscript.CreateJobSubmission(dm.ConfigFile)

	for i := 1; i <= numSims; i++ {
		simNum := offset + i
		inputs, err := dm.InputGenerator.SimInputByJobNumber(simNum)
		if err != nil {
			fmt.Printf("Error obtaining siminput by job number %s", err)
			dm.Logger.Error(fmt.Sprintf("Error Obtaining siminput for job number %d", simNum), err)
		}
		simLogger := log.NewLogger(0, fmt.Sprintf("Simulation: %d", simNum), "63")

		designObjectives := make(
			map[string]float64,
			len(dm.ConfigFile.DesignStudyConfig.DesignObjectives),
		)
		for _, objective := range dm.ConfigFile.DesignStudyConfig.DesignObjectives {
			designObjectives[objective.Name] = 0.0
		}
		sim := simulations.NewSimulation(
			&jobSubmission,
			simNum,
			inputs,
			simLogger,
			dm.ConfigFile.SlurmConfig,
			dm.ConfigFile.SlurmConfig.NodeList[hostIndex],
			dm.ConfigFile.Test.Function,
		)
		sim.Run()
		simParams, simResults := sim.ParseSimulationResults()
		dm.SimResultParams = simParams
		dm.SimResults = append(dm.SimResults, simResults)
	}
	dm.Logger.Log("Finished running design sweep")
}

// yes this is dumb
func (dm *DesignManager) HandleSim(sim *simulations.Simulation) {
	designObjectives := make(
		map[string]float64,
		len(dm.ConfigFile.DesignStudyConfig.DesignObjectives),
	)
	for _, objective := range dm.ConfigFile.DesignStudyConfig.DesignObjectives {
		designObjectives[objective.Name] = 0.0
	}

	sim.DesignObjectiveResults = designObjectives
	sim.Run()
	simParams, simResults := sim.ParseSimulationResults()
	dm.SimResultParams = simParams
	dm.SimResults = append(dm.SimResults, simResults)
	dm.Logger.Log("Finished handling sim")
}

func (dm *DesignManager) HandlePareto() {
	dsc := dm.ConfigFile.DesignStudyConfig
	mooConfig := dm.ConfigFile.DesignStudyConfig.MOOConfig
	numSimsPerGeneration := mooConfig.NumSimsPerGeneration

	dm.Logger.Log("Initializing the population")
	population := genetic.InitializePopulation(numSimsPerGeneration, dm.ConfigFile)

	for generation := 0; generation < mooConfig.NumGenerations; generation++ {
		dm.Logger.Log(fmt.Sprintf("Starting Generation: %d\n", generation))
		if generation == 0 {
			jobs := make(chan int, numSimsPerGeneration)
			results := make(chan []float64, numSimsPerGeneration)
			wg := sync.WaitGroup{}

			for i := range numSimsPerGeneration {
				wg.Add(1)
				jobs <- 1
				newDM := dm
				go func(i int) {
					defer wg.Done()
					newDM.HandleSim(population[i].Sim)
					<-jobs
					dm.SimResultParams = newDM.SimResultParams
					// dm.SimResults = append(dm.SimResults, newDM.SimResults...)
					results <- newDM.SimResults[0]
				}(i)
			}
			wg.Wait()
			close(results)

			dm.Logger.Log("Finished running generation 1")

			for result := range results {
				dm.SimResults = append(dm.SimResults, result)
			}

			dm.Logger.Log("Evaluating Best and sorting population")
			population = genetic.Evaluate(population, dsc)
			dm.Logger.Log("Sorted Population, top 2 will be parents of next: \n")
			for _, ind := range population {
				dm.Logger.Log(
					fmt.Sprintf(
						"Parameters: %v\n, ObjectiveResults: %v\n, Fitness: %f\n",
						ind.Sim.InputParameters,
						ind.Sim.DesignObjectiveResults,
						ind.Fitness,
					),
				)
				// dm.SimResults = append(dm.SimResults, ind.Sim.DesignObjectiveResults)
			}
			dm.Logger.Log("Saving compiled generation results")
			dm.SaveCompiledResults(
				fmt.Sprintf(
					"Gen_%d_sim_%s",
					generation,
					strings.TrimSuffix(dm.ConfigFile.StarCCM.SimFile, ".sim"),
				),
			)

			// clear slices
			dm.SimResults = nil
			dm.SimResultParams = nil

			continue
		}
		newPopulation := make(genetic.Population, 0, numSimsPerGeneration)
		i := 1
		for len(newPopulation) < numSimsPerGeneration {
			parent1 := population[len(population)-1]
			parent2 := population[len(population)-2]

			// create sim for child
			jobSubmission := jobscript.CreateJobSubmission(dm.ConfigFile)
			simInputs := genetic.SampleInputs(dsc) // temp until crossover and mutate
			simNum := (generation * numSimsPerGeneration) + i - 1
			simLogger := log.NewLogger(0, fmt.Sprintf("Simulation: %d", simNum), "63")
			sim := simulations.NewSimulation(
				&jobSubmission,
				simNum,
				simInputs,
				simLogger,
				dm.ConfigFile.SlurmConfig,
				dm.ConfigFile.SlurmConfig.NodeList[i-1],
				dm.ConfigFile.Test.Function,
			)

			child := genetic.Individual{
				Sim:     sim,
				Fitness: 0.0,
			}

			genetic.Crossover(parent1, parent2, &child, dsc)
			genetic.Mutate(&child, mooConfig.MutationRate, dsc)
			newPopulation = append(newPopulation, child)

			i += 1
		}
		population = nil

		jobs := make(chan int, numSimsPerGeneration)
		results := make(chan []float64, numSimsPerGeneration)
		wg := sync.WaitGroup{}

		for i := range numSimsPerGeneration {
			wg.Add(1)
			jobs <- 1
			newDM := dm
			go func(i int) {
				defer wg.Done()
				newDM.HandleSim(newPopulation[i].Sim)
				<-jobs
				dm.SimResultParams = newDM.SimResultParams
				// dm.SimResults = append(dm.SimResults, newDM.SimResults...)
				results <- newDM.SimResults[0]
			}(i)

		}
		wg.Wait()
		close(results)

		dm.Logger.Log(fmt.Sprintf("Finshed running generation: %d", generation))

		for result := range results {
			dm.SimResults = append(dm.SimResults, result)
		}

		dm.Logger.Log("Evaluating Best and sorting population")
		newPopulation = genetic.Evaluate(newPopulation, dsc)
		dm.Logger.Log("Sorted Population, top 2 will be parents of next: \n")
		for _, ind := range newPopulation {
			dm.Logger.Log(
				fmt.Sprintf(
					"Parameters: %v\n, ObjectiveResults: %v\n, Fitness: %f\n",
					ind.Sim.InputParameters,
					ind.Sim.DesignObjectiveResults,
					ind.Fitness,
				),
			)
		}

		dm.Logger.Log("Saving compiled generation results")
		dm.SaveCompiledResults(
			fmt.Sprintf(
				"Gen_%d_sim_%s",
				generation,
				strings.TrimSuffix(dm.ConfigFile.StarCCM.SimFile, ".sim"),
			),
		)

		population = newPopulation

		// clear slices
		dm.SimResults = nil
		dm.SimResultParams = nil
	}
	dm.Logger.Log("Finished running last generation\n\nFinal Population:")
	for _, ind := range population {
		dm.Logger.Log(
			fmt.Sprintf(
				"Parameters: %v\n, ObjectiveResults: %v\n, Fitness: %f\n",
				ind.Sim.InputParameters,
				ind.Sim.DesignObjectiveResults,
				ind.Fitness,
			),
		)
	}
}

func (dm *DesignManager) HandleCustom() {
	custom.HandleCustomAlg(dm.ConfigFile, dm.Logger, dm.DiscordHook)
}

func (dm *DesignManager) HandleDesignStudy(studyType string) {
	dm.DiscordHook.PayloadJson.Content = "Running Design Study"
	dm.DiscordHook.CallWebHook()
	switch studyType {
	case "AeroMap":
		dm.Logger.Log("Running AeroMap")
		dm.HandleAeroMap()
	case "Pareto":
		if dm.ConfigFile.DesignStudyConfig.MOOConfig.OptimizationAlgorithm == "Genetic" {
			dm.HandlePareto()
			dm.Logger.Log("Running Pareto Study")
		} else if dm.ConfigFile.DesignStudyConfig.MOOConfig.OptimizationAlgorithm == "Custom" {
			dm.Logger.Log("Running MOO study with Custom PSO algorithm")
			dm.HandleCustom()
		}
	case "Sweep":
		dm.Logger.Log("Running design sweep")
		dm.HandleSweep(0, dm.ConfigFile.DesignStudyConfig.NumSims, 0)
	default:
		fmt.Println("Error: Study type not supported")
		os.Exit(1)
	}
}

func (dm *DesignManager) SaveCompiledResults(fileName string) {
	resultsFile, err := os.Create("Compiled_" + fileName + "_Report.csv")
	if err != nil {
		dm.Logger.Error("Failed to create results file", err)
	}
	utils.WriteParameterCsvHeader(dm.SimResultParams, resultsFile)

	csvWriter := csv.NewWriter(resultsFile)
	defer csvWriter.Flush()
	for _, row := range dm.SimResults {

		strRow := make([]string, len(row))
		for j, value := range row {
			strRow[j] = strconv.FormatFloat(value, 'f', 4, 64)
		}
		if err := csvWriter.Write(strRow); err != nil {
			dm.Logger.Error("Failed to write result csv row", err)
		}
	}
}
