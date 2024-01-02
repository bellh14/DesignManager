package generator

import (
	"fmt"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"os"
)

func GenerateJobScript(jobScriptInputs types.JobSubmissionType) {

	//TODO: make this less painful to read

	jobScript, err := os.Create("jobScript.sh")
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer jobScript.Close()

	jobScript.WriteString("#!/bin/bash\n\n")

	jobScript.WriteString("WORKING_DIR=")
	jobScript.WriteString(jobScriptInputs.WorkingDir)
	jobScript.WriteString("\n")

	jobScript.WriteString("NCPU=")
	jobScript.WriteString(fmt.Sprint(jobScriptInputs.Ntasks))
	jobScript.WriteString("\n")

	jobScript.WriteString("PODKEY=")
	jobScript.WriteString(jobScriptInputs.StarCCM.PodKey)
	jobScript.WriteString("\n")

	jobScript.WriteString("JAVA_MACRO=")
	jobScript.WriteString(jobScriptInputs.StarCCM.JavaMacro)
	jobScript.WriteString("\n")

	jobScript.WriteString("SIM_FILE=")
	jobScript.WriteString(jobScriptInputs.StarCCM.SimFile)
	jobScript.WriteString("\n")

	jobScript.WriteString("JOB_NUMBER=")
	jobScript.WriteString(fmt.Sprint(jobScriptInputs.JobNumber))
	jobScript.WriteString("\n\n")

	jobScript.WriteString("mkdir -p $WORKING_DIR/$JOB_NUMBER\n\n")

	jobScript.WriteString("module load starccm/17.04.007\n")

	jobScript.WriteString(`starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PODKEY -batch $WORKING_DIR/$JOB_NUMBER/$JAVA_MACRO $WORKING_DIR/$JOB_NUMBER/$SIM_FILE -np $NCPU -bs slurm -time -batch-report`)
}
