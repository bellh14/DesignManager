package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bellh14/DesignManager/config"
	"github.com/bellh14/DesignManager/pkg/designmanager"

)

func main() {
	// Parse command line arguments
	inputFile := flag.String("config", "", "Input file")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: Input file not specified")
		os.Exit(1)
	}

	// Parse config file
	config := config.ParseConfigFile(*inputFile)
	fmt.Println("Input config file is: ", *inputFile)

	// Create design manager
	designManager := designmanager.NewDesignManager(config)
	designManager.Run()
}
