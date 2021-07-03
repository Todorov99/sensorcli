package sensor

// ISensor gets and validate sensor measurements.
type ISensor interface {
	GetSensorData(arguments ...string) ([]string, error)
	Validate(arguments ...string) error
}
