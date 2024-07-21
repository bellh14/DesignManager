package batchsystem

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

type SlurmConfig struct {
	JobName    string `json:"JobName"`
	Partition  string `json:"Partition"`
	Nodes      int    `json:"Nodes"`
	Ntasks     int    `json:"Ntasks"`
	WallTime   string `json:"WallTime"` // "hh:mm:ss"
	Email      string `json:"Email"`
	MailType   string `json:"MailType"` // "begin", "end", "fail", "all"
	OutputFile string `json:"OutputFile"`
	ErrorFile  string `json:"ErrorFile"`
	WorkingDir string `json:"WorkingDir"`
}

func WriteSlurmVariable(file *os.File, name string, value any) {
	switch name {

	case "JobName":
		fmt.Fprintf(file, "#SBATCH\t\t-J \"%s\"\n", value)
	case "Partition":
		fmt.Fprintf(file, "#SBATCH\t\t-p %s\n", value)
	case "Nodes":
		fmt.Fprintf(file, "#SBATCH\t\t-N %d\n", value)
	case "Ntasks":
		fmt.Fprintf(file, "#SBATCH\t\t-n %d\n", value)
	case "WallTime":
		fmt.Fprintf(file, "#SBATCH\t\t-t %s\n", value)
	case "Email":
		fmt.Fprintf(file, "#SBATCH\t\t-mail-user=%s\n", value)
	case "MailType":
		fmt.Fprintf(file, "#SBATCH\t\t-mail-type=%s\n", value)
	case "OutputFile":
		fmt.Fprintf(file, "#SBATCH\t\t-o \"%s\"\n", value)
	case "ErrorFile":
		fmt.Fprintf(file, "#SBATCH\t\t-e \"%s\"\n\n", value)
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

func GenerateSlurmScript(slurmConfig SlurmConfig, configFile string) {
	// TODO: make this less painful to read

	slurmScript, err := os.Create(
		fmt.Sprintf("%s%s.sh", slurmConfig.WorkingDir, slurmConfig.JobName),
	)
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer slurmScript.Close()

	slurmScript.WriteString("#!/bin/bash\n\n")

	slurmConfigValues := reflect.ValueOf(slurmConfig)

	WriteStructOfSlurmVariables(slurmConfigValues, slurmScript)

	fmt.Fprintf(slurmScript, "./DesignManager -config %s", configFile)

	err = os.Chmod(fmt.Sprintf("%s%s.sh", slurmConfig.WorkingDir, slurmConfig.JobName), 0o777)
	if err != nil {
		log.Fatal(err)
	}
}
