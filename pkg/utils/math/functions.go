package math

import (
	"math"

	"github.com/bellh14/DesignManager/pkg/types"
)


func CalculateMean(values []float64) float64 {
	var sum float64
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func CalculateStandardDeviation(values []float64) float64 {
	mean := CalculateMean(values)
	var sum float64
	for _, value := range values {
		sum += math.Pow(value-mean, 2)
	}
	return math.Sqrt(sum / float64(len(values)))
}

func CalculateVariance(values []float64) float64 {
	mean := CalculateMean(values)
	var sum float64
	for _, value := range values {
		sum += math.Pow(value-mean, 2)
	}
	return sum / float64(len(values))
}

func CalculateRange(min, max float64) float64 {
	return max - min
}

func CalculateNumSamples(parameterRange, step float64) int {
	return int(math.Floor(parameterRange / step))
}

func CalculateParamterPopulation(designParameter *types.DesignParameter) []float64 {
	parameterRange := CalculateRange(designParameter.Min, designParameter.Max)
	numSamples := CalculateNumSamples(parameterRange, designParameter.Step)
	population := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		population[i] = designParameter.Min + float64(i)*designParameter.Step
	}
	return population
}
