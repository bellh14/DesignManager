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

type DesignParameter struct {
	Name  string  `json:"name"`
	Units string  `json:"units"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Step  float64 `json:"step"`
}
type DesignStudyConfig struct {
	StudyType        string            `json:"studyType"`
	StudyConfigDir   string            `json:"studyConfigDir"` // optional dir for storing study configs ie sim inputs
	NumSims          int               `json:"numSims"`
	DesignParameters []DesignParameter `json:"designParameters"`
}

type ConfigFile struct {
	UseDM       bool                    `json:"useDM"` // use dm or just output generated scripts
	OutputDir   string                  `json:"outputDir"`
	SlurmConfig batchsystem.SlurmConfig `json:"slurmConfig"`
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

	json.Unmarshal(byteValue, &config)
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

	json.Unmarshal(byteValue, &config)
	return config
}
