package jobscript_test

import (
	"testing"

	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/types"

)

func TestGenerateJobScript(t *testing.T) {
	jobScriptInputs := types.JobSubmissionType{
		WorkingDir: "../../../scripts/",
		Ntasks:     4,
		Path:       "/opt/Siemens/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin/",
		PodKey:     "1234-5678-9012-3456",
		JavaMacro:  "DMPareto.java",
		SimFile:    "sim.sim",
	}
	jobscript.GenerateJobScript(jobScriptInputs, 1)
}
