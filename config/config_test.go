package config_test

import (
	"testing"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/generator/batchsystem"
)

func compareSlurmConfig(t *testing.T, got, want batchsystem.SlurmConfig) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareStarCCM(t *testing.T, got, want config.StarCCM) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareDesignParameter(t *testing.T, got, want config.DesignParameter) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareDesignObjective(t *testing.T, got, want config.DesignObjective) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareDesignManagerInputParameters(
	t *testing.T,
	got, want config.DesignStudyConfig,
) {
	t.Helper()
	if got.NumSims != want.NumSims {
		t.Errorf("got %v want %v", got.NumSims, want.NumSims)
	}
	if got.NtasksPerSim != want.NtasksPerSim {
		t.Errorf("got %v want %v", got.NtasksPerSim, want.NtasksPerSim)
	}
	if got.StudyType != want.StudyType {
		t.Errorf("got %v want %v", got.StudyType, want.StudyType)
	}
	if got.OptimizationAlgorithm != want.OptimizationAlgorithm {
		t.Errorf("got %v want %v", got.OptimizationAlgorithm, want.OptimizationAlgorithm)
	}
	for i := range got.DesignParameters {
		compareDesignParameter(t, got.DesignParameters[i], want.DesignParameters[i])
	}
	for i := range got.DesignObjectives {
		compareDesignObjective(t, got.DesignObjectives[i], want.DesignObjectives[i])
	}
}

func compareAll(t *testing.T, got, want config.ConfigFile) {
	t.Helper()
	compareSlurmConfig(t, got.SlurmConfig, want.SlurmConfig)
	if got.WorkingDir != want.WorkingDir {
		t.Errorf("got %v want %v", got.WorkingDir, want.WorkingDir)
	}
	compareStarCCM(t, got.StarCCM, want.StarCCM)
	compareDesignManagerInputParameters(
		t,
		got.DesignStudyConfig,
		want.DesignStudyConfig,
	)
}

func TestParseDesignManagerConfigFile(t *testing.T) {
	configFilePath := "../data/inputs/aeromap.json"
	configFile := config.ParseConfigFile(configFilePath)

	expectedSlurmInputs := batchsystem.SlurmConfig{
		WorkingDir: ".",
		JobName:    "2024AeroSweep",
		Nodes:      1,
		Ntasks:     16,
		Partition:  "icx",
		WallTime:   "24:00:00",
		Email:      "test@gmail.com",
		MailType:   "all",
		OutputFile: "output.txt",
		ErrorFile:  "error.txt",
	}

	expectedWorkingDir := "."

	expectedStarCCM := config.StarCCM{
		StarPath:  "/opt/Siemens/19.02.013/STAR-CCM+19.02.013/star/bin/",
		PodKey:    "123456789",
		JavaMacro: "AirfoilAOA.java",
		SimFile:   "S1223.sim",
	}

	expectedDesignParameter1 := config.DesignParameter{
		Name:    "Chassis Angle",
		Units:   "deg",
		Min:     -1.3,
		Max:     1.3,
		NumSims: 9,
	}

	expectedDesignParameter2 := config.DesignParameter{
		Name:    "Chassis Heave",
		Units:   "inches",
		Min:     -1.69,
		Max:     0.31,
		NumSims: 9,
	}

	expectedDesignParameters := []config.DesignParameter{
		expectedDesignParameter1,
		expectedDesignParameter2,
	}

	expectedDesignObjective1 := config.DesignObjective{
		Name:   "Design Objective 1",
		Weight: 1.0,
		Goal:   "Maximize",
	}

	expectedDesignObjective2 := config.DesignObjective{
		Name:   "Design Objective 2",
		Weight: 0.75,
		Goal:   "Minimize",
	}

	expectedDesignObjectives := []config.DesignObjective{
		expectedDesignObjective1,
		expectedDesignObjective2,
	}

	expectedDesignStudyConfig := config.DesignStudyConfig{
		NumSims:               81,
		NtasksPerSim:          16,
		StudyType:             "AeroMap",
		OptimizationAlgorithm: "NSGA-II",
		DesignParameters:      expectedDesignParameters,
		DesignObjectives:      expectedDesignObjectives,
	}

	expectedConfigFile := config.ConfigFile{
		SlurmConfig:       expectedSlurmInputs,
		WorkingDir:        expectedWorkingDir,
		StarCCM:           expectedStarCCM,
		DesignStudyConfig: expectedDesignStudyConfig,
	}

	compareAll(t, configFile, expectedConfigFile)
}
