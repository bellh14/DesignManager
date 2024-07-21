package simulations_test

import (
	"os"
	"testing"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/generator/inputs"
	"github.com/bellh14/DesignManager/pkg/generator/jobscript"
	"github.com/bellh14/DesignManager/pkg/simulations"
	"github.com/bellh14/DesignManager/pkg/utils/log"
)

func TestHandleSimulation(t *testing.T) {
	t.Helper()

	sim := testCreateNewSimulation(t)
	testSetWorkingDir(t, sim)

	sim.CreateSimulationDirectory()
	testCreateSimulationDir(t)

	sim.CopySimulationFiles()
	testCopySimulationFiles(t, sim)

	sim.CreateSimulationInputFile()
	testCreateSimulationInputFile(t)

	sim.CreateJobScript()
	testCreateJobScript(t)
}

func testCreateNewSimulation(t *testing.T) *simulations.Simulation {
	t.Helper()
	jobSubmission := jobscript.JobSubmission{
		WorkingDir: "../../test/testoutput",
		Ntasks:     80,
		StarPath:   "/opt/Siemens/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin",
		PodKey:     "1234-5678-9012-3456",
		JavaMacro:  "AirfoilAOA.java",
		SimFile:    "testsim.sim",
	}

	designParameters := []config.DesignParameter{
		{
			Name:    "Parameter1",
			Min:     -1.3,
			Max:     1.3,
			NumSims: 9,
		},
		{
			Name:    "Parameter2",
			Min:     -1.3,
			Max:     1.3,
			NumSims: 9,
		},
	}
	inputFileName := jobSubmission.WorkingDir + "/" + "testInputs.csv"
	inputGenerator := inputs.NewSimInputGenerator(designParameters, inputFileName)
	inputGenerator.HandleSimInputs()
	inputs, err := inputGenerator.SimInputByJobNumber(1)
	if err != nil {
		t.Errorf("Error obtaining siminput by job number %s", err)
	}
	logger := log.NewLogger(0, "Simulation Test", "63")
	return simulations.NewSimulation(&jobSubmission, 1, inputs, logger)
}

func testSetWorkingDir(t *testing.T, sim *simulations.Simulation) {
	t.Helper()

	expectedWorkingDir := "../../test/testoutput/1"

	sim.SetWorkingDir()
	actualWorkingDir := sim.JobDir

	if actualWorkingDir != expectedWorkingDir {
		t.Errorf("Got: %s Expected: %s", actualWorkingDir, expectedWorkingDir)
	}
}

func testCreateSimulationDir(t *testing.T) {
	_, err := os.Stat("../../test/testoutput/1")
	if os.IsNotExist(err) {
		t.Errorf("Error: Failed to create simulation directory")
	}
}

func testCopySimulationFiles(t *testing.T, sim *simulations.Simulation) {
	_, err := os.Stat("../../test/testoutput/1/" + sim.JobSubmission.JavaMacro)
	if os.IsNotExist(err) {
		t.Errorf("Error: Failed to copy Java macro: %s", err)
	}
	_, err = os.Stat("../../test/testoutput/1/" + sim.JobSubmission.SimFile)
	if os.IsNotExist(err) {
		t.Errorf("Error: Failed to copy sim file: %s", err)
	}
}

func testCreateSimulationInputFile(t *testing.T) {
	_, err := os.Stat("../../test/testoutput/1/" + "InputParams.csv")
	if os.IsNotExist(err) {
		t.Errorf("Error: Failed to create simulation input file: %s", err)
	}
}

func testCreateJobScript(t *testing.T) {
	_, err := os.Stat("../../test/testoutput/1/sim_1.sh")
	if os.IsNotExist(err) {
		t.Errorf("Error: Failed to create simulation job script: %s", err)
	}
}
