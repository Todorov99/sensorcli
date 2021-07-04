package sensor

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/host"
	"github.com/ttodorov/sensorcli/pkg/util"
)

const (
	tempSensor string = "CPU_TEMP"
)

type cpuTempSensor Measurment

// CreateTempSensor creates instance of temperature sensor.
func CreateTempSensor() ISensor {
	return &cpuTempSensor{}
}

func (tempS *cpuTempSensor) GetSensorData(ctx context.Context, unit, format string) ([]string, error) {
	sensorLogger.Info("Gerring sensor data...")
	cpuTemp, err := getTempMeasurments(ctx, unit, format)
	if err != nil {
		msg := "failed to get temperature measurements"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	return cpuTemp, nil
}

func (tempS *cpuTempSensor) Validate(arguments ...string) error {
	return util.ValidateFormat(arguments[0])
}

func getTempMeasurments(ctx context.Context, unit string, format string) ([]string, error) {
	var tempData []string

	cpuTempFromSensor, err := getTempFromSensor(ctx)
	if err != nil {
		return nil, err
	}

	deviceID, err := devices.getDeviceID()
	if err != nil {
		return nil, err
	}

	sensorID, err := devices.getSensorID(tempSensor)

	temperatureInCurrentUnit := util.ParseTempAccordingToUnit(unit, cpuTempFromSensor)

	measurement := newMeasurement(temperatureInCurrentUnit, sensorID[0], deviceID)
	parsedData := util.ParseDataAccordingToFormat(format, measurement)

	tempData = append(tempData, parsedData)

	return tempData, nil
}

func getTempFromSensor(ctx context.Context) (float64, error) {
	sensorLogger.Info("Getting temperature from sensor")
	sensorTeperatureInfo, err := host.SensorsTemperaturesWithContext(ctx)
	if err != nil {
		return 0, err
	}

	if len(sensorTeperatureInfo) == 0 {
		return ReadFileSystemFile("/sys/class/thermal/cooling_device2/device/status")
	}

	cpuTemp := sensorTeperatureInfo[0].Temperature
	sensorLogger.Info("Temperature from sensor is successfully got")

	return cpuTemp, nil
}
