package sampling

import (
	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/utils"
	"github.com/bellh14/DesignManager/pkg/utils/math/probability"
)

type Sampler struct {
	DesignParameters []config.DesignParameter
	// Distribution   string  TODO: implement this
}

func NewSampler(config config.DesignStudyConfig) *Sampler {
	utils.SeedRand()
	return &Sampler{
		DesignParameters: config.DesignParameters,
	}
}

func (sampler *Sampler) SampleParameter(designParameter config.DesignParameter) float64 {
	return probability.NormalDistribution(designParameter.Mean, designParameter.StdDev)
}
