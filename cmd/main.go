package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/designmanager"
	"github.com/bellh14/DesignManager/pkg/utils/log"
)

func main() {
	// creat logger
	logger := log.NewLogger(0) // Parse command line arguments

	inputFile := flag.String("config", "", "Input file")
	flag.Parse()

	if *inputFile == "" {
		logger.Fatal("Input file not specified", fmt.Errorf("no config.json file"))
		os.Exit(1)
	}

	// Parse config file
	config := config.ParseConfigFile(*inputFile)
	fmt.Println("Input config file is: ", *inputFile)

	// Create design manager
	designManager := designmanager.NewDesignManager(config, logger)
	designManager.Run()
}
