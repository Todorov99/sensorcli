package sensor

import "context"

// ISensor gets and validate sensor measurements.
type ISensor interface {
	GetSensorData(ctx context.Context, format string) ([]string, error)
	ValidateFormat(format string) error
	ValidateUnit() error
}
