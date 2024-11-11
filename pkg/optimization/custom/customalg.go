package custom

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/optimization/genetic"
	"github.com/bellh14/DesignManager/pkg/simulations"
	"github.com/bellh14/DesignManager/pkg/utils/log"
)

func CalculateObjectivesPeaks(
	population genetic.Population,
) (map[string]float64, map[string]float64) {
	minValues := make(map[string]float64)
	maxValues := make(map[string]float64)

	for _, individual := range population {
		for objective, value := range individual.Sim.DesignObjectiveResults {
			if _, exists := minValues[objective]; !exists {
				minValues[objective] = math.Inf(1)
				maxValues[objective] = math.Inf(-1)
			}
			if value < minValues[objective] {
				minValues[objective] = value
			}
			if value > maxValues[objective] {
				maxValues[objective] = value
			}
		}
	}
	return minValues, maxValues
}

func Normalize(result, min, max float64, goal string) float64 {
	if goal == "Maximize" {
		return (result - min) / (max - min)
	}
	return (max - result) / (max - min)
}

func CalculateFitness(
	individual *genetic.Individual,
	dsc config.DesignStudyConfig,
	minValues, maxValues map[string]float64,
) {
	i := 0
	for objective, result := range individual.Sim.DesignObjectiveResults {
		goal := dsc.DesignObjectives[i].Goal
		target := dsc.DesignObjectives[i].Target
		weight := dsc.DesignObjectives[i].Weight

		normalizedResult := Normalize(
			result,
			minValues[objective],
			maxValues[objective],
			goal,
		)

		if goal == "Maximize" {
			if target != 0 {
				if result >= float64(target) {
					individual.Fitness += normalizedResult * float64(weight)
				} else {
					individual.Fitness -= normalizedResult * float64(weight)
				}
			} else {
				individual.Fitness += normalizedResult * float64(weight)
			}
		} else {
			if target != 0 {
				if result <= float64(target) {
					individual.Fitness += normalizedResult * float64(weight)
				} else {
					individual.Fitness -= normalizedResult * float64(weight)
				}
			}
		}
		i += 1
	}
}

func Evaluate(
	population genetic.Population,
	dsc config.DesignStudyConfig,
	logger *log.Logger,
) genetic.Population {
	minValues, maxValues := CalculateObjectivesPeaks(population)
	for i := range population {
		CalculateFitness(&population[i], dsc, minValues, maxValues)
	}

	sort.Sort(population)

	for _, p := range population {
		logger.Log(fmt.Sprintf("Sim: %d, Fitness: %f\n", p.Sim.JobNumber, p.Fitness))
	}
	return population
}

func HandleSim(sim *simulations.Simulation, dsc config.DesignStudyConfig) map[string]float64 {
	designObjectives := make(
		map[string]float64,
		len(dsc.DesignObjectives),
	)
	for _, objective := range dsc.DesignObjectives {
		designObjectives[objective.Name] = 0.0
	}

	sim.DesignObjectiveResults = designObjectives
	sim.Run()
	if sim.Successful {
		_, _ = sim.ParseSimulationResults()
	}
	return sim.DesignObjectiveResults
}

func UpdateInputs(population *genetic.Population, dsc config.DesignStudyConfig) {
	for i := range len(dsc.DesignParameters) {
		bestValue := (*population)[population.Len()-1].Sim.InputParameters.Value[i]

		for j, ind := range *population {
			distance := math.Abs(bestValue - ind.Sim.InputParameters.Value[i])
			scalingFactor := dsc.DesignParameters[i].ScalingFactor
			updateValue := ind.Sim.InputParameters.Value[i] + scalingFactor*distance*float64(
				population.Len()-j,
			)/float64(
				population.Len(),
			)
			newInd := ind
			newInd.Sim.InputParameters.Value[i] = updateValue
			(*population)[j] = newInd
		}
	}
}

func HandleCustomAlg(config config.ConfigFile, logger *log.Logger) {
	dsc := config.DesignStudyConfig
	mooConfig := dsc.MOOConfig
	numSimsPerGen := mooConfig.NumSimsPerGeneration

	logger.Log("Initializing the Population")
	population := genetic.InitializePopulation(numSimsPerGen, config)

	for generation := 0; generation < mooConfig.NumGenerations; generation++ {
		logger.Log(fmt.Sprintf("Starting Generation: %d\n", generation))
		if generation == 0 {
			HandleGeneration(&population, dsc)
			logger.Log("Finished running generation 0")
			logger.Log("Evaluating Population")
			population = Evaluate(population, dsc, logger)
			PrintResults(population, logger)
			continue
		}
		logger.Log("Updating simulation parameters")
		UpdateInputs(&population, dsc)
	}
}

func HandleGeneration(population *genetic.Population, dsc config.DesignStudyConfig) {
	numSimsPerGen := dsc.MOOConfig.NumSimsPerGeneration
	jobs := make(chan int, numSimsPerGen)
	// results := make(chan []float64, numSimsPerGen)
	wg := sync.WaitGroup{}

	for i := range numSimsPerGen {
		wg.Add(1)
		jobs <- 1
		go func(i int) {
			defer wg.Done()
			time.Sleep(10 * time.Second)
			HandleSim((*population)[i].Sim, dsc)
			<-jobs
		}(i)
	}

	wg.Wait()
}

func PrintResults(population genetic.Population, logger *log.Logger) {
	for _, ind := range population {
		logger.Log(
			fmt.Sprintf(
				"Parameters: %v\n, ObjectiveResults: %v\n, Fitness: %f\n",
				ind.Sim.InputParameters,
				ind.Sim.DesignObjectiveResults,
				ind.Fitness,
			),
		)
	}
}
