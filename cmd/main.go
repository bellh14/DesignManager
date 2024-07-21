package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/designmanager"
	"github.com/bellh14/DesignManager/pkg/generator/batchsystem"
	"github.com/bellh14/DesignManager/pkg/utils/log"
)

func main() {
	// creat logger
	logger := log.NewLogger(0, "DM", "#941ff4") // Parse command line arguments

	inputFile := flag.String("config", "", "Input file")
	batchSystemFlag := flag.String("bs", "", "batch system (only supports slurm right now)")
	flag.Parse()

	if *inputFile == "" {
		logger.Fatal("Input file not specified", fmt.Errorf("no config.json file"))
		os.Exit(1)
	}

	_, err := os.Stat(*inputFile)
	if os.IsNotExist(err) {
		logger.Fatal("Input File does not exist", err)
		os.Exit(1)
	}

	// Parse config file
	config := config.ParseConfigFile(*inputFile)
	logger.Log(fmt.Sprintf("Input config file is: %s", *inputFile))

	if *batchSystemFlag == "slurm" {
		batchsystem.GenerateSlurmScript(config.SlurmConfig, *inputFile)
	}

	// Create design manager
	designManager := designmanager.NewDesignManager(config, logger)
	designManager.Run()
}
