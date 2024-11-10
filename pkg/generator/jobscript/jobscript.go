package jobscript

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/generator/batchsystem"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/utils"
)

type JobSubmission struct {
	WorkingDir     string
	Ntasks         int
	StarPath       string
	PodKey         string
	JavaMacro      string
	SimFile        string
	Params         string
	StarWorkingDir string // I hate this, but will change later
	MachineFile    string
	NtasksPerNode  int
}

func CreateJobSubmission(config config.ConfigFile) JobSubmission {
	jobSumssion := JobSubmission{
		WorkingDir:     config.WorkingDir,
		Ntasks:         config.DesignStudyConfig.NtasksPerSim,
		StarPath:       config.StarCCM.StarPath,
		PodKey:         config.StarCCM.PodKey,
		JavaMacro:      config.StarCCM.JavaMacro,
		SimFile:        config.StarCCM.SimFile,
		StarWorkingDir: config.StarCCM.WorkingDir,
		NtasksPerNode:  config.DesignStudyConfig.NtasksPerNode,
	}
	return jobSumssion
}

func CreateParamsString(inputs inputs.SimInputIteration) string {
	paramString := ""
	for i := range inputs.Name {
		paramString += fmt.Sprintf("-param \"%s\" %f ", inputs.Name[i], inputs.Value[i])
	}
	return paramString
}

func GenerateJobScript(
	jobScriptInputs JobSubmission,
	jobNumber int,
	inputs inputs.SimInputIteration,
	slurmConfig batchsystem.SlurmConfig,
	hostNodes string,
	testFunction string,
) {
	paramString := CreateParamsString(inputs)

	// TODO: make this less painful to read
	jobDir := jobScriptInputs.WorkingDir + "/" + fmt.Sprint(jobNumber)
	jobScriptInputs.StarWorkingDir += "/" + fmt.Sprint(jobNumber)
	jobScript, err := os.Create(fmt.Sprintf("%s/sim_%d.sh", jobDir, jobNumber))
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer jobScript.Close()

	jobScriptInputs.WorkingDir = jobScriptInputs.StarWorkingDir

	jobScript.WriteString("#!/bin/bash\n\n")

	batchsystem.WriteStructOfSlurmVariables(reflect.ValueOf(slurmConfig), jobScript)

	jobSubmissionValues := reflect.ValueOf(jobScriptInputs)

	utils.WriteStructOfBashVariables(jobSubmissionValues, jobScript, []string{})

	jobScript.WriteString("\nmodule load starccm/18.06.006\ncd $StarWorkingDir\n\n")

	// jobScript.WriteString("mkdir $WorkingDir/$JobNumber\n\n")

	// coreOffset := jobScriptInputs.Ntasks * (jobNumber % 2) // TODO temp since we are running 56x2

	if testFunction != "" {
		simName := strings.TrimSuffix(jobScriptInputs.SimFile, ".sim")
		reportName := simName + "_Report.csv"
		jobScript.WriteString(
			fmt.Sprintf(
				"./MOOT -f InputParams.csv -o %s --test %s > $WorkingDir/output.txt 2>&1",
				reportName,
				testFunction,
			),
		)

	} else {
		if jobScriptInputs.StarPath == "" {
			jobScript.WriteString(
				fmt.Sprintf(
					"starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JavaMacro $WorkingDir/$SimFile -np $Ntasks %s -machinefile $WorkingDir/%s -time -batch-report > $WorkingDir/output.txt 2>&1",
					// jobScriptInputs.Ntasks,
					// coreOffset,
					paramString,
					hostNodes,
				),
			)
		} else {
			jobScript.WriteString(
				fmt.Sprintf(
					"$StarPath/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JavaMacro $WorkingDir/$SimFile -np $Ntasks %s -on %s -time -batch-report > $WorkingDir/output.txt 2>&1",
					// jobScriptInputs.Ntasks,
					// coreOffset,
					paramString,
					hostNodes,
				),
			)
		}
	}
	// jobScript.WriteString(
	// 	fmt.Sprintf(
	// 		"$StarPath/starccm+ -power -licpath 1999@flex.cd-adapco.com -podkey $PodKey -batch $WorkingDir/$JavaMacro $Working/Dir/$SimFile -np $Ntasks %s -time -batch-report > $WorkingDir/output.txt 2>&1",
	// 		paramString,
	// 	),
	// )

	jobScript.WriteString("\n\n")
	jobScript.WriteString("exit_code=$?\n")
	jobScript.WriteString("if [ $exit_code -ne 0 ]; then\n")
	jobScript.WriteString(
		"    echo \"Error: StarCCM+ exited with non-zero exit code: $exit_code\" >&2\n",
	)
	jobScript.WriteString("    exit $exit_code\n")
	jobScript.WriteString("fi\n\n")

	err = os.Chmod(fmt.Sprintf("%s/sim_%d.sh", jobDir, jobNumber), 0o777)
	if err != nil {
		log.Fatal(err)
	}
}
