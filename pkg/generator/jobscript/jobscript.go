package jobscript

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/bellh14/DesignManager/pkg/types"
	"github.com/bellh14/DesignManager/pkg/utils"

)

func GenerateJobScript(jobScriptInputs types.JobSubmissionType, jobNumber int) {
	// TODO: make this less painful to read

	jobScript, err := os.Create(fmt.Sprintf("%s/job.sh", jobScriptInputs.WorkingDir))
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer jobScript.Close()

	jobScript.WriteString("#!/bin/bash\n\n")

	jobSubmissionValues := reflect.ValueOf(jobScriptInputs)

	utils.WriteStructOfBashVariables(jobSubmissionValues, jobScript)

	// jobScript.WriteString("mkdir $WorkingDir/$JobNumber\n\n")

	jobScript.WriteString("module load starccm/17.04.007\n")

	jobScript.WriteString(`starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JavaMacro $WorkingDir/$SimFile -np $Ntasks -time -batch-report`)

	jobScript.WriteString("\n\n")
	jobScript.WriteString("exit_code=$?\n")
	jobScript.WriteString("if [ $exit_code -ne 0 ]; then\n")
	jobScript.WriteString("    echo \"Error: StarCCM+ exited with non-zero exit code: $exit_code\" >&2\n")
	jobScript.WriteString("    exit $exit_code\n")
	jobScript.WriteString("fi\n\n")

	err = os.Chmod(fmt.Sprintf("%s/job.sh", jobScriptInputs.WorkingDir), 0o777)
	if err != nil {
		log.Fatal(err)
	}
}
