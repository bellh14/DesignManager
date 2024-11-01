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

func TestParseNodeListSingles(t *testing.T) {
	t.Helper()
	hostName := "stampede3.tacc.utexas.edu"
	nodeList := "c479-092,c486-112"
	nodes, err := batchsystem.ParseNodeList(nodeList, hostName)
	if err != nil {
		t.Errorf("Failed to parse node list %s", err)
	}
	if nodes[0] != "c479-092.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[0], "c479-092.stampede3.tacc.utexas.edu")
	}
	if nodes[1] != "c486-112.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[1], "c486-112.stampede3.tacc.utexas.edu")
	}
	if len(nodes) == 0 {
		t.Errorf("Failed to parse node list %s", nodes)
	}
}

func TestParseNodeListArrays(t *testing.T) {
	t.Helper()
	hostName := "stampede3.tacc.utexas.edu"
	nodeList := "c460-[222-223],c461-[201-202]"
	nodes, err := batchsystem.ParseNodeList(nodeList, hostName)
	if err != nil {
		t.Errorf("Failed to parse node list %s", err)
	}
	if nodes[0] != "c460-222.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[0], "c460-222.stampede3.tacc.utexas.edu")
	}
	if nodes[1] != "c460-223.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[1], "c460-223.stampede3.tacc.utexas.edu")
	}
	if nodes[2] != "c461-201.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[2], "c461-201.stampede3.tacc.utexas.edu")
	}
	if nodes[3] != "c461-202.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[3], "c461-202.stampede3.tacc.utexas.edu")
	}
	if len(nodes) == 0 {
		t.Errorf("Failed to parse node list %s", nodes)
	}
}
