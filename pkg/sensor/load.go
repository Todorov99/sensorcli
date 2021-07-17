package sensor

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	fileFullPath    string = "./model.yaml"
	csvFileFullPath string = "./"

	cpuTempCelsius    string = "cpuTemp"
	cpuCoresCount     string = "cpuCores"
	cpuUsagePercent   string = "cpuUsage"
	cpuFrequency      string = "cpuFrequency"
	memoryTotal       string = "memoryTotal"
	memoryAvailable   string = "memoryAvailable"
	memoryUsed        string = "memoryUsed"
	memoryUsedPercent string = "memoryUsedPercent"
)

var devices *Diveces

// load the content of the model.yaml file
func init() {
	fileName, err := filepath.Abs("./model.yaml")
	if err != nil {
		sensorLogger.Panic(err)
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		sensorLogger.Panic(err)
	}

	fileErr := yaml.Unmarshal(yamlFile, &devices)
	if fileErr != nil {
		sensorLogger.Panic(err)
	}

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
