package batchsystem_test

import (
	"testing"

	"github.com/bellh14/DesignManager/pkg/generator/batchsystem"
)

func TestGenerateSlurmScript(t *testing.T) {
	slurmInputs := batchsystem.SlurmConfig{
		WorkingDir: "../../../test/testoutput/",
		JobName:    "testjob",
		Nodes:      1,
		Ntasks:     4,
		Partition:  "icx",
		WallTime:   "24:00:00",
		Email:      "test@gmail.com",
		MailType:   "all",
		OutputFile: "output.txt",
		ErrorFile:  "error.txt",
	}
	batchsystem.GenerateSlurmScript(slurmInputs, "configfile")
}
