package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bellh14/DFRDesignManager/pkg/types"
	"io"
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

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func SeedRand() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func WriteParameterCsv(samples []types.SimInput, file *os.File) {
	for i, sample := range samples {
		file.WriteString(fmt.Sprintf("%v", sample.Value))
		if i < len(samples)-1 {
			file.WriteString(",")
		}
	}
	file.WriteString("\n")
}

func WriteParameterCsvHeader(designParameters []types.SimInput, file *os.File) {
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
