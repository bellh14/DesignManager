package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"math/rand"
	"os"
	"reflect"
	"time"
)

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func WriteBashVariable(file *os.File, name string, value any) {
	file.WriteString(fmt.Sprintf("%s=%v\n", name, value))
}

func WriteStructOfBashVariables(values reflect.Value, file *os.File) {
	for i := 0; i < values.NumField(); i++ {
		value := values.Field(i)
		name := values.Type().Field(i).Name
		WriteBashVariable(file, name, value.Interface())
	}
}

func SeedRand() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func WriteParameterCsv(samples types.ParameterSamples, file *os.File) {
	for i, sample := range samples.Samples {
		file.WriteString(fmt.Sprintf("%v", sample))
		if i < len(samples.Samples)-1 {
			file.WriteString(",")
		}
	}
	file.WriteString("\n")
}

func WriteParameterCsvHeader(designParameters []types.DesignParameter, file *os.File) {
	for i, designParameter := range designParameters {
		file.WriteString(designParameter.Name)
		if i < len(designParameters)-1 {
			file.WriteString(",")
		}
	}
	file.WriteString("\n")
}

func CreateJobSubmission(systemResources types.SystemResourcesType, workingDir string, starCCM types.StarCCM) types.JobSubmissionType {
	return types.JobSubmissionType{
		WorkingDir: workingDir,
		Ntasks:     systemResources.Ntasks,
		Path:       starCCM.Path,
		PodKey:     starCCM.PodKey,
		JavaMacro:  starCCM.JavaMacro,
		SimFile:    starCCM.SimFile,
	}
}
