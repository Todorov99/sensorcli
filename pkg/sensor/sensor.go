package sensor

import (
	"fmt"
)

type Sensor struct {
	ID           int32        `json:"id" yaml:"id" mapstructure:"id,omitempty"`
	Name         string       `json:"name" yaml:"name" mapstructure:"name,omitempty"`
	Description  string       `json:"description" yaml:"description" mapstructure:"description,omitempty"`
	Unit         string       `json:"unit" yaml:"unit" mapstructure:"unit,omitempty"`
	SensorGroups string       `json:"sensorGroups" yaml:"sensorGroups" mapstructure:"sensorGroups,omitempty"`
	Measurments  []Measurment `json:"measurements" yaml:"measurements"`
}

func (s *Sensor) getSensorIDAccordingToSensorName(sensorName string, currentSensorID string) (string, error) {
	switch sensorName {
	case cpuTemp:
		return currentSensorID, nil
	case cpuUsage:
		return currentSensorID, nil
	case cpuCores:
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

	return "", fmt.Errorf("there is not such sensor name: %q", sensorName)
}
