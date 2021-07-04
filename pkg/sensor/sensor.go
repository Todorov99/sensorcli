package sensor

import (
	"fmt"
	"sort"
)

// Sensor model
type Sensor struct {
	ID           string       `json:"id" yaml:"id"`
	Name         string       `json:"name" yaml:"name"`
	Description  string       `json:"description" yaml:"description"`
	Unit         string       `json:"unit" yaml:"unit"`
	SensorGroups []string     `json:"sensorGroups" yaml:"sensorGroups"`
	Measurments  []Measurment `json:"measurements" yaml:"measurements"`
}

// GetSensorID of specified sensor group.
func (d *Diveces) getSensorID(sensorType string) ([]string, error) {
	var sensorsIds []string

	for _, currentSensor := range d.Devices[0].Sensors {

		if currentSensor.SensorGroups[0] == sensorType {

			currentSensorID, err := getSensorIDAccordingToSensorName(currentSensor.Name, currentSensor.ID)
			if err != nil {
				return nil, err
			}

			sensorsIds = append(sensorsIds, currentSensorID)
		}
	}

	sort.Sort(sort.StringSlice(sensorsIds))
	return sensorsIds, nil
}

// GetTempSensorUnit gets current unit for temperature sensor measurment.
func GetSensorUnit(sensorGroup string) (string, error) {
	sensorLogger.Info("Getting current unit for temperature sensor")

	for _, currentsensor := range devices.Devices[0].Sensors {
		for _, sgr := range currentsensor.SensorGroups {

			if sgr == sensorGroup {
				return currentsensor.Unit, nil
			}
		}
	}

	return "", fmt.Errorf("failed to get sensor unit")
}

func getSensorIDAccordingToSensorName(sensorName string, currentSensorID string) (string, error) {

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
	case memoryAvailableBytes:
		return currentSensorID, nil
	case memoryUsedBytes:
		return currentSensorID, nil
	case memoryUsedPercent:
		return currentSensorID, nil

	}

	return "", fmt.Errorf("there is not such sensor name")
}
