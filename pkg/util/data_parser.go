package util

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

const (
	jsonType string = "JSON"
	yamlType string = "YAML"
)

// ParseDataAccordingToFormat returns corect parsed data.
func ParseDataAccordingToFormat(format string, data interface{}) string {
	if format == yamlType {
		return yamlParser(data)
	}

	return jsonParser(data)
}

func jsonParser(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func yamlParser(data interface{}) string {
	yamlData, _ := yaml.Marshal(data)
	return string(yamlData)
}
