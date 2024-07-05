package batchsystem

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

type SlurmConfig struct {
	JobName    string `json:"jobName"`
	Partition  string `json:"partition"`
	Nodes      int    `json:"nodes"`
	Ntasks     int    `json:"ntasks"`
	WallTime   string `json:"wallTime"` // "hh:mm:ss"
	Email      string `json:"email"`
	MailType   string `json:"mailType"` // "begin", "end", "fail", "all"
	OutputFile string `json:"outputFile"`
	ErrorFile  string `json:"errorFile"`
	WorkingDir string `json:"workingDir"`
}

func WriteSlurmVariable(file *os.File, name string, value any) {
	switch name {

	case "JobName":
		file.WriteString(fmt.Sprintf("#SBATCH 	-J \"%s\"\n", value))
	case "Partition":
		file.WriteString(fmt.Sprintf("#SBATCH 	-p %s\n", value))
	case "Nodes":
		file.WriteString(fmt.Sprintf("#SBATCH 	-N %d\n", value))
	case "Ntasks":
		file.WriteString(fmt.Sprintf("#SBATCH 	-n %d\n", value))
	case "WallTime":
		file.WriteString(fmt.Sprintf("#SBATCH 	-t %s\n", value))
	case "Email":
		file.WriteString(fmt.Sprintf("#SBATCH 	-mail-user=%s\n", value))
	case "MailType":
		file.WriteString(fmt.Sprintf("#SBATCH 	-mail-type=%s\n", value))
	case "OutputFile":
		file.WriteString(fmt.Sprintf("#SBATCH 	-o \"%s\"\n", value))
	case "ErrorFile":
		file.WriteString(fmt.Sprintf("#SBATCH 	-e \"%s\"\n", value))
	case "WorkingDir":
		return
	default:
		return
	}
}

func WriteStructOfSlurmVariables(values reflect.Value, file *os.File) {
	for i := 0; i < values.NumField(); i++ {
		value := values.Field(i)
		name := values.Type().Field(i).Name
		WriteSlurmVariable(file, name, value.Interface())
	}
}

func GenerateSlurmScript(slurmConfig SlurmConfig) {
	// TODO: make this less painful to read

	slurmScript, err := os.Create(fmt.Sprintf("%s%s.sh", slurmConfig.WorkingDir, slurmConfig.JobName))
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer slurmScript.Close()

	slurmScript.WriteString("#!/bin/bash\n\n")

	slurmConfigValues := reflect.ValueOf(slurmConfig)

	WriteStructOfSlurmVariables(slurmConfigValues, slurmScript)

	err = os.Chmod(fmt.Sprintf("%s%s.sh", slurmConfig.WorkingDir, slurmConfig.JobName), 0o777)
	if err != nil {
		log.Fatal(err)
	}
}
