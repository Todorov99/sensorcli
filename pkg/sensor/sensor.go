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
	Measurments  []measurment `json:"measurements" yaml:"measurements"`
}

func (d *Diveces) getSensorsByGroup(sensorGroup string) ([]sensor, error) {
	var sensorsIds []sensor

	for _, currentSensor := range d.Devices[0].Sensors {

		if currentSensor.SensorGroups[0] == sensorGroup {

			_, err := currentSensor.getSensorIDAccordingToSensorName(currentSensor.Name, currentSensor.ID)
			if err != nil {
				return nil, err
			}

			sensorsIds = append(sensorsIds, currentSensor)
		}
	}

	//sort.Sort(sort.StringSlice(sensorsIds))
	return sensorsIds, nil
}

// GetTempSensorUnit gets current unit for temperature sensor measurment.
func GetSensorUnits(sensorGroup string) ([]string, error) {
	sensorLogger.Info("Getting current unit for temperature sensor")
	units := []string{}

	for _, currentsensor := range devices.Devices[0].Sensors {
		for _, sgr := range currentsensor.SensorGroups {

			if sgr == sensorGroup {
				units = append(units, currentsensor.Unit)
			}
		}
	}

	return units, nil
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
	case memoryAvailableBytes:
		return currentSensorID, nil
	case memoryUsedBytes:
		return currentSensorID, nil
	case memoryUsedPercent:
		return currentSensorID, nil

	}

	return "", fmt.Errorf("there is not such sensor name")
}
