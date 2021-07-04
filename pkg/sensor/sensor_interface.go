package sensor

import "context"

// ISensor gets and validate sensor measurements.
type ISensor interface {
	GetSensorData(ctx context.Context, unit, format string) ([]string, error)
	Validate(arguments ...string) error
}
