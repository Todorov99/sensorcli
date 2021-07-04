package sensor

import "time"

// Measurment model
type Measurment struct {
	MeasuredAt time.Time `json:"measuredAt" yaml:"measuredAt"`
	Value      string    `json:"value" yaml:"value"`
	SensorID   string    `json:"sensorId" yaml:"sensorId"`
	DeviceID   string    `json:"deviceId" yaml:"deviceId"`
}

// SetMeasurementValues sets property fields of measurement model.
func newMeasurement(value string, sensorID string, deviceID string) Measurment {
	return Measurment{time.Now(), value, sensorID, deviceID}
}
