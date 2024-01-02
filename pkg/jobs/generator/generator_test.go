package generator_test

import (
	"github.com/bellh14/DFRDesignManager/pkg/jobs/generator"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"testing"
)

func TestGenerateJobScript(t *testing.T) {
	jobScriptInputs := types.JobSubmissionType{
		WorkingDir: "../../../scripts/run_simulation.sh",
		Ntasks: 4,
		StarCCM: types.StarCCM{
			PodKey: "1234-5678-9012-3456",
			JavaMacro: "macro.java",
			SimFile: "sim.sim",
		},
		JobNumber: 1,
	}
	generator.GenerateJobScript(jobScriptInputs)


}