package generator

import (
	"fmt"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"github.com/bellh14/DFRDesignManager/pkg/utils"
	"os"
	"reflect"
)

func GenerateJobScript(jobScriptInputs types.JobSubmissionType) {

	//TODO: make this less painful to read

	jobScript, err := os.Create("JobScript.sh")
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer jobScript.Close()

	jobScript.WriteString("#!/bin/bash\n\n")

	jobSubmissionValues := reflect.ValueOf(jobScriptInputs)

	utils.WriteStructOfBashVariables(jobSubmissionValues, jobScript)

	jobScript.WriteString("mkdir $WorkingDir/$JobNumber\n\n")

	jobScript.WriteString(`$Path/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JobNumber/$JavaMacro $WorkingDir/$JobNumber/$SimFile -np $Ntasks -time -batch-report`)
}
