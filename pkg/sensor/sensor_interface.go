package sensor

import "context"

// ISensor gets and validate sensor measurements.
type ISensor interface {
	GetSensorData(ctx context.Context, arguments ...string) ([]string, error)
	Validate(arguments ...string) error
}
