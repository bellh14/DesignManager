// Package: config
// Should parse input config file and return a struct with the config values
package config

import (
	"github.com/bellh14/DesignManager/pkg/types"
	"os"
	"io"
	"encoding/json"
	"fmt"
)

func ParseConfigFile(configFilePath string) types.ConfigFile {

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
