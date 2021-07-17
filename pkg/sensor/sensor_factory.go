package sensor

import (
	"fmt"

	"github.com/ttodorov/sensorcli/pkg/logger"
)

var sensorLogger logger.Logger = logger.NewLogger("./sensor")

// CreateSensor sensor type instance.
func NewSensor(sensorType string) (ISensor, error) {

	switch sensorType {
	case tempSensor:
		sensorLogger.Info("Getting temp sensor measurements.")
		return CreateTempSensor(), nil
	case usageSensor:
		sensorLogger.Info("Getting usage sensor measurements.")
		return CreateUsageSensor(), nil
	case memorySensor:
		sensorLogger.Info("Getting memory sensor measurements.")
		return CreateMemorySensor(), nil
	}

	return nil, fmt.Errorf("error in getting sensor type")
}
