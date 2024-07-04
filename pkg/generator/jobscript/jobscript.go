package jobscript

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/bellh14/DesignManager/pkg/types"
	"github.com/bellh14/DesignManager/pkg/utils"
)

type StarCCM struct {
	StarPath  string
	PodKey    string
	JavaMacro string
}

func GenerateJobScript(jobScriptInputs types.JobSubmissionType, jobNumber int) {
	// TODO: make this less painful to read

	jobScript, err := os.Create(fmt.Sprintf("%ssim_%d.sh", jobScriptInputs.WorkingDir, jobNumber))
	if err != nil {
		// TODO: handle error

		fmt.Println(err)
	}
	defer jobScript.Close()

	jobScript.WriteString("#!/bin/bash\n\n")

	jobSubmissionValues := reflect.ValueOf(jobScriptInputs)

	utils.WriteStructOfBashVariables(jobSubmissionValues, jobScript, []string{"DesignParameters"})

	// jobScript.WriteString("mkdir $WorkingDir/$JobNumber\n\n")

	jobScript.WriteString(`$StarPath/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JavaMacro $WorkingDir/$SimFile -np $Ntasks -time -batch-report`)

	jobScript.WriteString("\n\n")
	jobScript.WriteString("exit_code=$?\n")
	jobScript.WriteString("if [ $exit_code -ne 0 ]; then\n")
	jobScript.WriteString("    echo \"Error: StarCCM+ exited with non-zero exit code: $exit_code\" >&2\n")
	jobScript.WriteString("    exit $exit_code\n")
	jobScript.WriteString("fi\n\n")

	err = os.Chmod(fmt.Sprintf("%ssim_%d.sh", jobScriptInputs.WorkingDir, jobNumber), 0o777)
	if err != nil {
		log.Fatal(err)
	}
}
