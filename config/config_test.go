package config_test

import (
	"testing"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/types"
)

func compareSystemResourcesType(t *testing.T, got, want types.SystemResourcesType) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareStarCCM(t *testing.T, got, want types.StarCCM) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareDesignParameter(t *testing.T, got, want types.DesignParameter) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareDesignObjective(t *testing.T, got, want types.DesignObjective) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareDesignManagerInputParameters(t *testing.T, got, want types.DesignManagerInputParameters) {
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

func compareAll(t *testing.T, got, want types.ConfigFile) {
	t.Helper()
	compareSystemResourcesType(t, got.SystemResources, want.SystemResources)
	if got.WorkingDir != want.WorkingDir {
		t.Errorf("got %v want %v", got.WorkingDir, want.WorkingDir)
	}
	compareStarCCM(t, got.StarCCM, want.StarCCM)
	compareDesignManagerInputParameters(t, got.DesignManagerInputParameters, want.DesignManagerInputParameters)
}

func TestParseDesignManagerConfigFile(t *testing.T) {
	configFilePath := "../data/inputs/DesignManagerConfig.json"
	configFile := config.ParseDesignManagerConfigFile(configFilePath)

	expectedSystemResources := types.SystemResourcesType{
		Partition: "normal",
		Nodes:     16,
		Ntasks:    16,
	}

	expectedWorkingDir := "/scratch/ganymede/<user>/DM/"

	expectedStarCCM := types.StarCCM{
		StarPath:  "/opt/Siemens/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin/",
		PodKey:    "<podkey>",
		JavaMacro: "macro.java",
		SimFile:   "simfile.sim",
	}

	expectedDesignParameter1 := types.DesignParameter{
		Name:    "Design Parameter 1",
		Min:     0.0,
		Max:     1.0,
		Step:    0.1,
		NumSims: 8,
	}

	expectedDesignParameter2 := types.DesignParameter{
		Name:    "Design Parameter 2",
		Min:     0.0,
		Max:     1.0,
		Step:    0.1,
		NumSims: 8,
	}

	expectedDesignParameters := []types.DesignParameter{
		expectedDesignParameter1,
		expectedDesignParameter2,
	}

	expectedDesignObjective1 := types.DesignObjective{
		Name:   "Design Objective 1",
		Weight: 1.0,
		Goal:   "Maximize",
	}

	expectedDesignObjective2 := types.DesignObjective{
		Name:   "Design Objective 2",
		Weight: 0.75,
		Goal:   "Minimize",
	}

	expectedDesignObjectives := []types.DesignObjective{
		expectedDesignObjective1,
		expectedDesignObjective2,
	}

	expectedDesignManagerParameters := types.DesignManagerInputParameters{
		NumSims:               100,
		NtasksPerSim:          -1,
		StudyType:             "Pareto",
		OptimizationAlgorithm: "NSGA-II",
		DesignParameters:      expectedDesignParameters,
		DesignObjectives:      expectedDesignObjectives,
	}

	expectedConfigFile := types.ConfigFile{
		SystemResources:              expectedSystemResources,
		WorkingDir:                   expectedWorkingDir,
		StarCCM:                      expectedStarCCM,
		DesignManagerInputParameters: expectedDesignManagerParameters,
	}

	compareAll(t, configFile, expectedConfigFile)
}
