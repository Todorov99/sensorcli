package sensor

import (
	"fmt"

	"github.com/ttodorov/sensorcli/pkg/logger"
)

var sensorLogger logger.Logger = logger.NewLogger("./pkg/sensor")

// CreateSensor sensor type instance.
func CreateSensor(sensorType string) (ISensor, error) {

	switch sensorType {
	case "CPU_TEMP":
		sensorLogger.Info("Getting temp sensor measurements.")
		return CreateTempSensor(), nil
	case "CPU_USAGE":
		sensorLogger.Info("Getting usage sensor measurements.")
		return CreateUsageSensor(), nil
	case "MEMORY_USAGE":
		sensorLogger.Info("Getting memory sensor measurements.")
		return CreateMemorySensor(), nil
	}

	return nil, fmt.Errorf("Error in getting sensor type.")
}
