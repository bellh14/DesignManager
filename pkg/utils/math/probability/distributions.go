package probability

import (
	"math"
	"math/rand"
)

func UniformDistribution(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func NormalDistribution(mean, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + mean
}

func LogNormalDistribution(mean, stdDev float64) float64 {
	return math.Exp(rand.NormFloat64()*stdDev + mean)
}
