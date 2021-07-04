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
