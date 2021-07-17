package sensor

import "fmt"

// Diveces models
type Diveces struct {
	Devices []Device `json:"devices" yaml:"devices"`
}

// GetDeviceID gets device id.
func (d *Diveces) getDeviceID() (string, error) {
	if devices.Devices[0].ID == "" {
		return "", fmt.Errorf("there is not available device")
	}

	return devices.Devices[0].ID, nil
}

func (d *Diveces) getDeviceSensorsByGroup(sensorGroup string) ([]sensor, error) {
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

	return sensorsIds, nil
}
