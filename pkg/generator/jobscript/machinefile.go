package jobscript

import (
	"fmt"
	"os"
	"strings"
)

func CreateMachineFile(fileName string, nodes string, ntasks int) error {
	machinefile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer machinefile.Close()
	splitNodes := strings.Split(nodes, ",")

	for _, node := range splitNodes {
		formattedHostNode := fmt.Sprintf("%s:%d", node, ntasks)
		_, err := machinefile.WriteString(formattedHostNode + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
