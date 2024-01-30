package generator

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/bellh14/DFRDesignManager/pkg/types"
	"github.com/bellh14/DFRDesignManager/pkg/utils"
)

func GenerateJobScript(jobScriptInputs types.JobSubmissionType, jobNumber int) {

	//TODO: make this less painful to read

	jobScript, err := os.Create(fmt.Sprintf("%s/job_%d.sh", jobScriptInputs.WorkingDir, jobNumber))
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer jobScript.Close()

	jobScript.WriteString("#!/bin/bash\n\n")

	jobSubmissionValues := reflect.ValueOf(jobScriptInputs)

	utils.WriteStructOfBashVariables(jobSubmissionValues, jobScript)

	// jobScript.WriteString("mkdir $WorkingDir/$JobNumber\n\n")

	jobScript.WriteString(`$Path/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JobNumber/$JavaMacro $WorkingDir/$JobNumber/$SimFile -np $Ntasks -time -batch-report`)

	jobScript.WriteString("\n\n")
	jobScript.WriteString("exit_code=$?\n")
	jobScript.WriteString("if [ $exit_code -ne 0 ]; then\n")
	jobScript.WriteString("    echo \"Error: StarCCM+ exited with non-zero exit code: $exit_code\" >&2\n")
	jobScript.WriteString("    exit $exit_code\n")
	jobScript.WriteString("fi\n\n")

	err = os.Chmod(fmt.Sprintf("%s/job_%d.sh", jobScriptInputs.WorkingDir, jobNumber), 0777)
	if err != nil {
		log.Fatal(err)
	}
}
