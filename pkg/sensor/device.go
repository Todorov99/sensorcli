package sensor

import "fmt"

// Device model
type Device struct {
	ID          string   `json:"id" yaml:"id" mapstructure:"id,omitempty"`
	Name        string   `json:"name" yaml:"name" mapstructure:"name,omitempty"`
	Description string   `json:"description" yaml:"description" mapstructure:"description,omitempty"`
	Sensors     []Sensor `json:"sensors" yaml:"sensors" mapstructure:"sensors,omitempty"`
}

// GetDeviceID gets device id.
func (d *Device) GetDeviceID() (string, error) {
	sensorLogger.Info("Getting device id...")
	if device.ID == "" {
		return "", fmt.Errorf("there is not available device")
	}

	return device.ID, nil
}

func (d *Device) GetDeviceSensorsByGroup(sensorGroup string) ([]Sensor, error) {
	sensorLogger.Info("Getting device %q sensors by %q group...", d.Name, sensorGroup)
	var sensorsIds []Sensor

	for _, currentSensor := range d.Sensors {
		if currentSensor.SensorGroups == sensorGroup {
			_, err := currentSensor.getSensorIDAccordingToSensorName(currentSensor.Name, currentSensor.ID)
			if err != nil {
				sensorLogger.Error(err)
				return nil, fmt.Errorf("failed to get device sensors by %q group: %w", sensorGroup, err)
			}

			sensorsIds = append(sensorsIds, currentSensor)
		}
	}

	return sensorsIds, nil
}
