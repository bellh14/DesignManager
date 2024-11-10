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

func TestParseNodeListSinglePlusArray(t *testing.T) {
	t.Helper()
	hostName := "stampede3.tacc.utexas.edu"
	nodeList := "c511-[014,021-023]"
	nodes, err := batchsystem.ParseNodeList(nodeList, hostName)
	if err != nil {
		t.Errorf("Failed to parse node list %s", err)
	}
	if len(nodes) != 4 {
		t.Errorf("failed to parse node list: icorrect number of nodes %s", nodes)
	}
	if nodes[0] != "c511-014.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[0], "c511-014.stampede.tacc.utexas.edu")
	}
	if nodes[1] != "c511-021.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[1], "c511-021.stampede.tacc.utexas.edu")
	}
	if nodes[2] != "c511-022.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[2], "c511-022.stampede.tacc.utexas.edu")
	}
	if nodes[3] != "c511-023.stampede3.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[3], "c511-023.stampede.tacc.utexas.edu")
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

func TestParseNodeListFull(t *testing.T) {
	t.Helper()
	hostName := "ls6.tacc.utexas.edu"
	nodeList := "c304-[013,023-024],c479-092,c486-112,c460-[222-223],c461-[201-203,231]"
	nodes, err := batchsystem.ParseNodeList(nodeList, hostName)
	if err != nil {
		t.Errorf("Failed to parse node list %s", err)
	}

	if nodes[0] != "c304-013.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[0], "c304-013.ls6.tacc.utexas.edu")
	}
	if nodes[1] != "c304-023.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[1], "c304-023.ls6.tacc.utexas.edu")
	}
	if nodes[2] != "c304-024.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[2], "c304-024.ls6.tacc.utexas.edu")
	}
	if nodes[3] != "c479-092.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[3], "c479-092.ls6.tacc.utexas.edu")
	}
	if nodes[4] != "c486-112.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[4], "c486-112.ls6.tacc.utexas.edu")
	}
	if nodes[5] != "c460-222.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[5], "c460-222.ls6.tacc.utexas.edu")
	}
	if nodes[6] != "c460-223.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[6], "c460-223.ls6.tacc.utexas.edu")
	}
	if nodes[7] != "c461-201.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[7], "c461-201.ls6.tacc.utexas.edu")
	}
	if nodes[8] != "c461-202.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[8], "c461-202.ls6.tacc.utexas.edu")
	}
	if nodes[9] != "c461-203.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[9], "c461-203.ls6.tacc.utexas.edu")
	}
	if nodes[10] != "c461-231.ls6.tacc.utexas.edu" {
		t.Errorf("Got %s : Wanted %s", nodes[10], "c461-231.ls6.tacc.utexas.edu")
	}
}
