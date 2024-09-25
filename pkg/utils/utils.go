package utils

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
)

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func WriteBashVariable(file *os.File, name string, value any) {
	_, err := fmt.Fprintf(file, "%s=%v\n", name, value)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func Contains(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func WriteStructOfBashVariables(values reflect.Value, file *os.File, excluded []string) {
	for i := 0; i < values.NumField(); i++ {
		if Contains(excluded, values.Type().Field(i).Name) {
			continue
		}
		value := values.Field(i)
		name := values.Type().Field(i).Name
		WriteBashVariable(file, name, value.Interface())
	}
}

func ReadLineByNumber(filePath string, lineNumber int) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)
	currentLine := 1

	// Read the file line by line
	for scanner.Scan() {
		if currentLine == lineNumber {
			// Found the line, return its content
			return scanner.Text(), nil
		}
		currentLine++
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Line number was not found
	return "", fmt.Errorf("line number %d out of range", lineNumber)
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return out.Close()
}

func SeedRand() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func WriteParameterCsv(samples [][]float64, file *os.File) {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range samples {
		strRow := make([]string, len(row))
		for i, value := range row {
			strRow[i] = strconv.FormatFloat(value, 'f', 4, 64)
		}

		// Write the string slice to the CSV file
		if err := writer.Write(strRow); err != nil {
			log.Fatalf("failed to write row: %s", err)
		}
	}
}

func WriteParameterCsvHeader(designParameters []string, file *os.File) {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(designParameters); err != nil {
		fmt.Println("Error: ", err)
	}
}

func WriteSimulationInputCSV(values []float64, file *os.File) {
	writer := csv.NewWriter(file)
	defer writer.Flush()

	strRow := make([]string, len(values))
	for i, value := range values {
		strRow[i] = strconv.FormatFloat(value, 'f', 4, 64)
	}
	if err := writer.Write(strRow); err != nil {
		log.Fatalf("Failed to write Input csv: %s", err)
	}
}

func ConvertStringSliceToFloat(strValues []string) ([]float64, error) {
	floatValues := make([]float64, len(strValues))
	for i, value := range strValues {
		if value == "" {
			continue
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		floatValues[i] = floatValue
	}
	return floatValues, nil
}
