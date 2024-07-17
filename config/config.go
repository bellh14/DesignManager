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
	StarPath  string `json:"StarPath"`
	PodKey    string `json:"PodKey"`
	JavaMacro string `json:"JavaMacro"`
	SimFile   string `json:"SimFile"`
}

type DesignParameter struct {
	Name    string  `json:"Name"`
	Units   string  `json:"Units"`
	Min     float64 `json:"Min"`
	Max     float64 `json:"Max"`
	Step    float64 `json:"Step"`
	NumSims int     `json:"NumSims"`
}
type DesignStudyConfig struct {
	StudyType        string            `json:"StudyType"`
	StudyConfigDir   string            `json:"StudyConfigDir"` // optional dir for storing study configs ie sim inputs
	NtasksPerSim     int               `json:"NtasksPerSim"`
	NumSims          int               `json:"NumSims"`
	DesignParameters []DesignParameter `json:"DesignParameters"`
}

type ConfigFile struct {
	UseDM             bool                    `json:"UseDM"` // use dm or just output generated scripts
	OutputDir         string                  `json:"OutputDir"`
	SlurmConfig       batchsystem.SlurmConfig `json:"SlurmConfig"`
	DesignStudyConfig DesignStudyConfig       `json:"DesignStudyConfig"`
	StarCCM           StarCCM                 `json:"Starccm"`
	WorkingDir        string                  `json:"WorkingDir"`
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

func ParseConfigFile(configFilePath string) ConfigFile {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	var config ConfigFile

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}
