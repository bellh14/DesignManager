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

	jobScript.WriteString("mkdir -p $WORKING_DIR/$JOB_NUMBER\n\n")

	jobScript.WriteString(`STARCCM_PATH/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PODKEY -batch $WORKING_DIR/$JOB_NUMBER/$JAVA_MACRO $WORKING_DIR/$JOB_NUMBER/$SIM_FILE -np $NCPU -bs slurm -time -batch-report`)
}
