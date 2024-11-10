package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

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
	slurmNodeList := flag.String("slurmNodeList", "", "List of slurm nodes allocated")
	nodesPerSim := flag.String("nps", "", "Number of nodes per sim if more than 1")
	testDM := flag.String("test", "", "test with fuction _")
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
	config, err := config.ParseConfigFile(*inputFile)
	if err != nil {
		logger.Fatal("Unable to parse config file", err)
	}
	logger.Log(fmt.Sprintf("Input config file is: %s", *inputFile))

	if *batchSystemFlag == "slurm" {
		batchsystem.GenerateSlurmScript(config.SlurmConfig, *inputFile)
		// dumb and tempory for now till refactor
		logger.Log("Created slurm batch script. Exiting...")
		os.Exit(0)
	}

	if *slurmNodeList != "" {
		logger.Log("Allocating on slurm nodes: " + *slurmNodeList)
		nodes, err := batchsystem.ParseNodeList(*slurmNodeList, config.SlurmConfig.HostName)
		if len(nodes) == 0 {
			if config.SlurmConfig.HostName == "" {
				nodes = append(nodes, *slurmNodeList)
			} else {
				nodes = append(nodes, *slurmNodeList+"."+config.SlurmConfig.HostName)
			}
		}
		if err != nil {
			logger.Fatal("Unable to parse slurm node list", err)
		}
		for _, node := range nodes {
			logger.Log(node)
		}
		simsPerNode := 1
		if *nodesPerSim == "" {
			simsPerNode = (config.SlurmConfig.Ntasks / config.SlurmConfig.Nodes) / config.DesignStudyConfig.NtasksPerSim
		}
		fullNodeList := batchsystem.DuplicateNodes(nodes, simsPerNode)

		config.SlurmConfig.NodeList = fullNodeList
	}

	if *nodesPerSim != "" {
		nps, err := strconv.Atoi(*nodesPerSim)
		if err != nil {
			logger.Fatal("Invalid nodes per sim, needs to be an integer", err)
		}
		logger.Log(fmt.Sprintf("Running each sim with %d nodes", nps))
		fullNodeList := batchsystem.AllocateMultiNodes(config.SlurmConfig.NodeList, nps)

		config.SlurmConfig.NodeList = fullNodeList
		config.DesignStudyConfig.NtasksPerNode = config.DesignStudyConfig.NtasksPerSim / nps
	}

	if *testDM != "" {
		config.Test.Test = true
		config.Test.Function = *testDM
	}

	// Create design manager
	designManager := designmanager.NewDesignManager(*config, logger)
	designManager.Run()
}
