package sensor

import (
	"fmt"
)

// Sensor model
type sensor struct {
	ID           string       `json:"id" yaml:"id"`
	Name         string       `json:"name" yaml:"name"`
	Description  string       `json:"description" yaml:"description"`
	Unit         string       `json:"unit" yaml:"unit"`
	SensorGroups []string     `json:"sensorGroups" yaml:"sensorGroups"`
	Measurments  []Measurment `json:"measurements" yaml:"measurements"`
}

func (s *sensor) getSensorIDAccordingToSensorName(sensorName string, currentSensorID string) (string, error) {

	switch sensorName {

	case cpuTempCelsius:
		return currentSensorID, nil
	case cpuUsagePercent:
		return currentSensorID, nil
	case cpuCoresCount:
		return currentSensorID, nil
	case cpuFrequency:
		return currentSensorID, nil
	case memoryTotal:
		return currentSensorID, nil
	case memoryAvailable:
		return currentSensorID, nil
	case memoryUsed:
		return currentSensorID, nil
	case memoryUsedPercent:
		return currentSensorID, nil

	}

	return "", fmt.Errorf("there is not such sensor name")
}
