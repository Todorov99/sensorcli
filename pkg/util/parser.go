package util

import (
	"encoding/json"
	"strconv"

	"github.com/ttodorov/sensorcli/pkg/logger"
	yaml "gopkg.in/yaml.v2"
)

var utilLogger logger.Logger = logger.NewLogger("./util")

const (
	unitTypeFahrenheit string = "F"
	unitTypeCelsius    string = "C"
	celsiusName        string = "Celsius"
	fahrenheitName     string = "Fahrenheit"

	jsonType string = "JSON"
	yamlType string = "YAML"
)

// ParseTempAccordingToUnit returns correct parsed temperature.
func ParseTempAccordingToUnit(unit string, temperature float64) string {
	if unit == unitTypeFahrenheit {
		utilLogger.Info("Parsing temperatur to fahrenheit unit")
		return strconv.FormatFloat(fahrenheit(temperature), 'f', 1, 64)
	}

	utilLogger.Info("Parsing temperatur to celsius unit")
	return strconv.FormatFloat(temperature, 'f', 1, 64)
}

// ParseDataAccordingToFormat returns corect parsed data.
func ParseDataAccordingToFormat(format string, data interface{}) string {
	if format == yamlType {
		utilLogger.Info("Parsing data to yaml")
		return yamlParser(data)
	}

	utilLogger.Info("Parsing data to json")
	return jsonParser(data)
}

func fahrenheit(temperature float64) float64 {
	return temperature*(9/5) + 32
}

func jsonParser(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func yamlParser(data interface{}) string {
	yamlData, _ := yaml.Marshal(data)
	return string(yamlData)
}
