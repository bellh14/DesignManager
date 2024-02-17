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

func NewSampler(job types.JobSubmissionType) *Sampler {
	utils.SeedRand()
	return &Sampler{
		DesignParameters: job.DesignParameters,
	}
}

func (sampler *Sampler) SampleParameter(designParameter types.DesignParameter) float64 {
	return probability.NormalDistribution(designParameter.Mean, designParameter.StdDev)
}

func (sampler *Sampler) Sample() []types.SimInput {
	samples := make([]types.SimInput, len(sampler.DesignParameters))
	for i, designParameter := range sampler.DesignParameters {
		samples[i].Name = designParameter.Name
		samples[i].Value = sampler.SampleParameter(designParameter)
	}
	return samples
}
