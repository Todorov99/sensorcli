package sensor

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	fileFullPath    string = "./model.yaml"
	csvFileFullPath string = "./"

	cpuTempCelsius       string = "cpuTempCelsius"
	cpuCoresCount        string = "cpuCoresCount"
	cpuUsagePercent      string = "cpuUsagePercent"
	cpuFrequency         string = "cpuFrequency"
	memoryTotal          string = "memoryTotal"
	memoryAvailableBytes string = "memoryAvailableBytes"
	memoryUsedBytes      string = "memoryUsedBytes"
	memoryUsedPercent    string = "memoryUsedPercent"
)

var devices *Diveces

// ReadYamlFile deserializing model from yaml file.
func ReadYamlFile(filePath string) error {
	fileName, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("error with getting yaml file name")
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	fileErr := yaml.Unmarshal(yamlFile, &devices)
	if fileErr != nil {
		return err
	}

	return nil
}

// WriteOutputToCSV measurement output to CSV file.
func WriteOutputToCSV(data []string, csvFileName string) error {

	fileName := csvFileFullPath + csvFileName + ".csv"

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Comma = '|'

	writingErr := writer.Write(data)
	if writingErr != nil {
		return nil
	}

	return nil
}

//ReadFileSystemFile reads server temperature from filesystem file.
func ReadFileSystemFile(fileSystemPath string) (float64, error) {

	fileName, err := filepath.Abs(fileSystemPath)

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	var temp float64

	err = yaml.Unmarshal(fileContent, &temp)
	if err != nil {
		return 0, err
	}

	return temp, nil
}
