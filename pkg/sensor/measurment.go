package sensor

import (
	"time"
)

// Measurment model
type measurment struct {
	MeasuredAt time.Time `json:"measuredAt" yaml:"measuredAt"`
	Value      string    `json:"value" yaml:"value"`
	SensorID   string    `json:"sensorId" yaml:"sensorId"`
	DeviceID   string    `json:"deviceId" yaml:"deviceId"`
}

// SetMeasurementValues sets property fields of measurement model.
func newMeasurement(value string, sensorID string, deviceID string) measurment {
	return measurment{time.Now(), value, sensorID, deviceID}
}

func newMeasurements(info interface{}) []measurment {
	var m = []measurment{}
	switch v := info.(type) {
	case cpuUsageSensor:
		for _, s := range v.sensors {
			switch s.Name {
			case cpuCoresCount:
				m = append(m, newMeasurement(v.cpuCores, s.ID, v.deviceID))
			case cpuFrequency:
				m = append(m, newMeasurement(v.cpuFrequency, s.ID, v.deviceID))
			case cpuUsagePercent:
				m = append(m, newMeasurement(v.cpuUsage, s.ID, v.deviceID))
			}

		}
	case cpuMemorySensor:
		for _, s := range v.sensors {
			switch s.Name {
			case memoryAvailableBytes:
				m = append(m, newMeasurement(v.availableMemory, s.ID, v.deviceID))
			case memoryUsedBytes:
				m = append(m, newMeasurement(v.usedMemory, s.ID, v.deviceID))
			case memoryTotal:
				m = append(m, newMeasurement(v.totalMemory, s.ID, v.deviceID))
			case memoryUsedPercent:
				m = append(m, newMeasurement(v.usedPercentMemory, s.ID, v.deviceID))
			}

		}
	case cpuTempSensor:
		for _, s := range v.sensors {
			switch s.Name {
			case cpuTempCelsius:
				m = append(m, newMeasurement(v.cpuTemp, s.ID, v.deviceID))
			}

		}
	}

	return m
}
