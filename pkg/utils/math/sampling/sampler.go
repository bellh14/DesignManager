package sampling

import (
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"github.com/bellh14/DFRDesignManager/pkg/utils"
	"github.com/bellh14/DFRDesignManager/pkg/utils/math/probability"
)

type Sampler struct {
	DesignParameters []types.DesignParameter
	//Distribution   string  TODO: implement this
}

func NewSampler(config types.ConfigFile) *Sampler {
	utils.SeedRand()
	return &Sampler{
		DesignParameters: config.DesignManagerInputParameters.DesignParameters,
	}
}

func (sampler *Sampler) SampleParameter(designParameter types.DesignParameter) float64 {
	return probability.NormalDistribution(designParameter.Min, designParameter.Max)
}

func (sampler *Sampler) Sample() types.ParameterSamples {
	samples := make([]float64, len(sampler.DesignParameters))
	for i, designParameter := range sampler.DesignParameters {
		samples[i] = sampler.SampleParameter(designParameter)
	}
	return types.ParameterSamples{
		Samples: samples,
	}
}
