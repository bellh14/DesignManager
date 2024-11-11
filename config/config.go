// Package: config
// Should parse input config file and return a struct with the config values
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bellh14/DesignManager/pkg/generator/batchsystem"
	"github.com/bellh14/DesignManager/pkg/types"
)

type StarCCM struct {
	StarPath   string `json:"StarPath"`
	PodKey     string `json:"PodKey"`
	JavaMacro  string `json:"JavaMacro"`
	SimFile    string `json:"SimFile"`
	WorkingDir string `json:"WorkingDir"` // dumb bs that needs to be set to start remote star remote servers
}

type DesignParameter struct {
	Name          string  `json:"Name"`
	Units         string  `json:"Units"`
	Min           float64 `json:"Min"`
	Max           float64 `json:"Max"`
	Step          float64 `json:"Step"`
	NumSims       int     `json:"NumSims"`
	ScalingFactor float64 `json:"ScalingFactor"`
	Mean          float64
	StdDev        float64
}

type DesignObjective struct {
	Name   string  `json:"Name"`
	Goal   string  `json:"Goal"`   // minimize or maximize, ex: df would want maximize while drag minimize
	Weight float32 `json:"Weight"` // may no explicitly use this
	Target float32 `json:"Target"`
}

type MOOConfig struct {
	NumGenerations        int     `json:"NumGenerations"`
	NumSimsPerGeneration  int     `json:"NumSimsPerGeneration"`
	OptimizationAlgorithm string  `json:"OptimizationAlgorithm"`
	MutationRate          float32 `json:"MutationRate"`
}

type DesignStudyConfig struct {
	StudyType             string            `json:"StudyType"`
	MOOConfig             MOOConfig         `json:"MOOConfig"`      // optional for temp genetic optimization
	StudyConfigDir        string            `json:"StudyConfigDir"` // optional dir for storing study configs ie sim inputs
	NtasksPerSim          int               `json:"NtasksPerSim"`
	NumSims               int               `json:"NumSims"`
	OptimizationAlgorithm string            `json:"OptimizationAlgorithm"`
	DesignParameters      []DesignParameter `json:"DesignParameters"`
	DesignObjectives      []DesignObjective `json:"DesignObjectives"`
	NtasksPerNode         int
}

type Test struct {
	Test     bool
	Function string
}

type ConfigFile struct {
	UseDM             bool                    `json:"UseDM"` // use dm or just output generated scripts
	OutputDir         string                  `json:"OutputDir"`
	SlurmConfig       batchsystem.SlurmConfig `json:"SlurmConfig"`
	DesignStudyConfig DesignStudyConfig       `json:"DesignStudyConfig"`
	StarCCM           StarCCM                 `json:"Starccm"`
	WorkingDir        string                  `json:"WorkingDir"`
	Test              Test
}

func ParseDesignManagerConfigFile(configFilePath string) types.ConfigFile {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	var config types.ConfigFile

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}

func ParseConfigFile(configFilePath string) (*ConfigFile, error) {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	var config ConfigFile

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
