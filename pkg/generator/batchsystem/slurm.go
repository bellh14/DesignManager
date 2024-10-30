package batchsystem

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type SlurmConfig struct {
	HostName   string   `json:"HostName"`
	JobName    string   `json:"JobName"`
	Partition  string   `json:"Partition"`
	Nodes      int      `json:"Nodes"`
	Ntasks     int      `json:"Ntasks"`
	WallTime   string   `json:"WallTime"` // "hh:mm:ss"
	Email      string   `json:"Email"`
	MailType   string   `json:"MailType"` // "begin", "end", "fail", "all"
	OutputFile string   `json:"OutputFile"`
	ErrorFile  string   `json:"ErrorFile"`
	WorkingDir string   `json:"WorkingDir"`
	NodeList   []string `json:"NodeList"`
}

func WriteSlurmVariable(file *os.File, name string, value any) {
	switch name {

	case "JobName":
		fmt.Fprintf(file, "#SBATCH\t\t-J \"%s\"\n", value)
	case "Partition":
		fmt.Fprintf(file, "#SBATCH\t\t-p %s\n", value)
	case "Nodes":
		fmt.Fprintf(file, "#SBATCH\t\t-N %d\n", value)
	case "Ntasks":
		fmt.Fprintf(file, "#SBATCH\t\t-n %d\n", value)
	case "WallTime":
		fmt.Fprintf(file, "#SBATCH\t\t-t %s\n", value)
	// case "Email":
	// 	fmt.Fprintf(file, "#SBATCH\t\t-mail-user=%s\n", value)
	// case "MailType":
	// 	fmt.Fprintf(file, "#SBATCH\t\t-mail-type=%s\n", value)
	case "OutputFile":
		fmt.Fprintf(file, "#SBATCH\t\t-o \"%s\"\n", value)
	case "ErrorFile":
		fmt.Fprintf(file, "#SBATCH\t\t-e \"%s\"\n\n", value)
	case "WorkingDir":
		return
	default:
		return
	}
}

func WriteStructOfSlurmVariables(values reflect.Value, file *os.File) {
	for i := 0; i < values.NumField(); i++ {
		value := values.Field(i)
		name := values.Type().Field(i).Name
		WriteSlurmVariable(file, name, value.Interface())
	}
}

func GenerateSlurmScript(slurmConfig SlurmConfig, configFile string) {
	// TODO: make this less painful to read

	slurmScript, err := os.Create(
		fmt.Sprintf("%s/%s.sh", slurmConfig.WorkingDir, slurmConfig.JobName),
	)
	if err != nil {
		// TODO: handle error
		fmt.Println(err)
	}
	defer slurmScript.Close()

	slurmScript.WriteString("#!/bin/bash\n\n")

	slurmConfigValues := reflect.ValueOf(slurmConfig)

	WriteStructOfSlurmVariables(slurmConfigValues, slurmScript)

	fmt.Fprintf(slurmScript, "./DesignManager -config %s", configFile)

	err = os.Chmod(fmt.Sprintf("%s/%s.sh", slurmConfig.WorkingDir, slurmConfig.JobName), 0o777)
	if err != nil {
		log.Fatal(err)
	}
}

func parseNodeRange(nodePrefix, rangePart string, hostName string) ([]string, error) {
	if !strings.Contains(rangePart, "-") {
		num, err := strconv.Atoi(rangePart)
		if err != nil {
			return nil, err
		}
		node := ""
		if hostName != "" {
			node = fmt.Sprintf("%s%d.%s", nodePrefix, num, hostName)
		} else {
			node = fmt.Sprintf("%s%d", nodePrefix, num)
		}
		return []string{node}, nil
	}

	rangeParts := strings.Split(rangePart, "-")

	start, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid node start number: %v", err)
	}
	end, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid end number: %v", err)
	}
	var nodes []string
	for i := start; i <= end; i++ {
		node := ""
		if hostName != "" {
			node = fmt.Sprintf("%s%d.%s", nodePrefix, i, hostName)
		} else {
			node = fmt.Sprintf("%s%d", nodePrefix, i)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func ParseNodeList(slurmNodeList string, hostName string) ([]string, error) {
	var allNodes []string

	// "c519-[051-054,061-064,071-074,081-084]"
	re := regexp.MustCompile(`([a-zA-Z0-9\-]+)\[(.*?)\]`)

	matches := re.FindAllStringSubmatch(slurmNodeList, -1)

	for _, match := range matches {
		nodePrefix := match[1]
		rangePart := match[2]

		ranges := strings.Split(rangePart, ",")

		for _, r := range ranges {
			nodes, err := parseNodeRange(nodePrefix, r, hostName)
			if err != nil {
				return nil, err
			}
			allNodes = append(allNodes, nodes...)
		}
	}

	return allNodes, nil
}

func DuplicateNodes(nodes []string, simsPerNode int) []string {
	var duplicateNodes []string

	for _, node := range nodes {
		for range simsPerNode {
			duplicateNodes = append(duplicateNodes, node)
		}
	}

	return duplicateNodes
}

func AllocateMultiNodes(nodes []string, nodesPerSim int) []string {
	var multiNodes []string

	for i := 0; i < len(nodes); i += nodesPerSim {

		end := i + nodesPerSim

		if end >= len(nodes) {
			end = len(nodes)
		}
		multiNode := ""
		for j := i; j < end; j += 1 {
			fmt.Println(nodes[j] + ",")
			if j != end-1 {
				multiNode += nodes[j] + ","
			} else {
				multiNode += nodes[j]
			}
		}

		multiNodes = append(multiNodes, multiNode)
		fmt.Println(multiNode)
	}
	return multiNodes
}
