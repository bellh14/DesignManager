package genetic

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/simulations"
	"github.com/bellh14/DesignManager/pkg/utils/log"
	"github.com/bellh14/DesignManager/pkg/utils/math/probability"
)

type Individual struct {
	Sim     *simulations.Simulation
	Fitness float64
}

type Population []Individual

func (p Population) Len() int {
	return len(p)
}

func (p Population) Less(i, j int) bool {
	return p[i].Fitness < p[j].Fitness
}

func (p Population) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func SampleInputs(inputConfig config.DesignStudyConfig) inputs.SimInputIteration {
	inputIteration := inputs.SimInputIteration{
		Name:  make([]string, len(inputConfig.DesignParameters)),
		Value: make([]float64, len(inputConfig.DesignParameters)),
	}

	for i, param := range inputConfig.DesignParameters {
		inputIteration.Name[i] = param.Name
		inputIteration.Value[i] = probability.UniformDistribution(param.Min, param.Max)
	}

	return inputIteration
}

func CalculateFitness(individual *Individual, dsc config.DesignStudyConfig) {
	i := 0
	for _, result := range individual.Sim.DesignObjectiveResults {
		if dsc.DesignObjectives[i].Goal == "Maximize" {
			individual.Fitness += result * float64(dsc.DesignObjectives[i].Weight)
		} else { // Minimize
			individual.Fitness -= result * float64(dsc.DesignObjectives[i].Weight)
		}
		i += 1
	}
}

func Evaluate(population Population, dsc config.DesignStudyConfig) Population {
	for i := range population {
		CalculateFitness(&population[i], dsc)
		fmt.Printf("Sim: %d, Fitness: %f\n", population[i].Sim.JobNumber, population[i].Fitness)
	}

	sort.Sort(population)
	for _, p := range population {
		fmt.Printf("Sim: %d, Fitness: %f\n", p.Sim.JobNumber, p.Fitness)
	}
	return population
}

func InitializePopulation(size int, dmConfig config.ConfigFile) Population {
	population := make(Population, size)
	jobSubmission := jobscript.CreateJobSubmission(dmConfig)

	for i := 0; i < size; i++ {
		simInputs := SampleInputs(dmConfig.DesignStudyConfig)
		simLogger := log.NewLogger(0, fmt.Sprintf("Simulation: %d", i), "63")
		sim := simulations.NewSimulation(
			&jobSubmission,
			i,
			simInputs,
			simLogger,
			dmConfig.SlurmConfig,
			dmConfig.SlurmConfig.NodeList[i],
		)
		individual := Individual{
			Sim:     sim,
			Fitness: 1.0,
		}
		population[i] = individual
	}

	return population
}

func Crossover(parent1, parent2 Individual, child *Individual, dsc config.DesignStudyConfig) {
	for i := range parent1.Sim.InputParameters.Value {
		alpha := rand.Float64()
		inRange := false
		for !inRange {
			child.Sim.InputParameters.Value[i] = alpha*parent1.Sim.InputParameters.Value[i] + (1-alpha)*parent2.Sim.InputParameters.Value[i]
			if child.Sim.InputParameters.Value[i] >= dsc.DesignParameters[i].Min &&
				child.Sim.InputParameters.Value[i] <= dsc.DesignParameters[i].Max {
				inRange = true
			}
		}
	}
}

func Mutate(individual *Individual, mutatationRate float32, dsc config.DesignStudyConfig) {
	for i := range individual.Sim.InputParameters.Value {
		if rand.Float32() < mutatationRate {
			individual.Sim.InputParameters.Value[i] = probability.UniformDistribution(
				dsc.DesignParameters[i].Min,
				dsc.DesignParameters[i].Max,
			)
		}
	}
}
