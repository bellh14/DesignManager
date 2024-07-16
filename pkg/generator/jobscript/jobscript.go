package jobscript

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/bellh14/DesignManager/pkg/utils"
)

type JobSubmission struct {
	WorkingDir string
	Ntasks     int
	StarPath   string
	PodKey     string
	JavaMacro  string
	SimFile    string
}

func GenerateJobScript(jobScriptInputs JobSubmission, jobNumber int) {
	// TODO: make this less painful to read
	jobDir := jobScriptInputs.WorkingDir + "/" + fmt.Sprint(jobNumber)
	jobScript, err := os.Create(fmt.Sprintf("%s/sim_%d.sh", jobDir, jobNumber))
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer jobScript.Close()

	jobScript.WriteString("#!/bin/bash\n\n")

	jobSubmissionValues := reflect.ValueOf(jobScriptInputs)

	utils.WriteStructOfBashVariables(jobSubmissionValues, jobScript, []string{})

	// jobScript.WriteString("mkdir $WorkingDir/$JobNumber\n\n")

	jobScript.WriteString(`$StarPath/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JavaMacro $WorkingDir/$SimFile -np $Ntasks -time -batch-report -bs slurm`)

	jobScript.WriteString("\n\n")
	jobScript.WriteString("exit_code=$?\n")
	jobScript.WriteString("if [ $exit_code -ne 0 ]; then\n")
	jobScript.WriteString("    echo \"Error: StarCCM+ exited with non-zero exit code: $exit_code\" >&2\n")
	jobScript.WriteString("    exit $exit_code\n")
	jobScript.WriteString("fi\n\n")

	err = os.Chmod(fmt.Sprintf("%s/sim_%d.sh", jobDir, jobNumber), 0o777)
	if err != nil {
		log.Fatal(err)
	}
}
