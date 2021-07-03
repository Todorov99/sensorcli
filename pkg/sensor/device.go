package sensor

import "fmt"

// Device model
type Device struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Sensors     []Sensor `json:"sensors" yaml:"sensors"`
}

// GetDeviceID gets device id.
func getDeviceID() (string, error) {
	devices, err := readYamlFile(fileFullPath)
	if err != nil {
		return "", err
	}

	if devices.Devices[0].ID == "" {
		return "", fmt.Errorf("There is no available device.")
	}

	return devices.Devices[0].ID, nil
}
