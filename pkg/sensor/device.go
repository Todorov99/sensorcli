package sensor

// Device model
type Device struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Sensors     []Sensor `json:"sensors" yaml:"sensors"`
}
