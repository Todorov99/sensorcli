package sensor

import (
	"context"

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

func (tempS *cpuTempSensor) GetSensorData(ctx context.Context, arguments ...string) ([]string, error) {
	cpuTemp, err := getTempMeasurments(ctx, arguments[0], arguments[1])
	if err != nil {
		return nil, err
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

	deviceID, err := getDeviceID()
	if err != nil {
		return nil, err
	}

	sensorID, err := GetSensorID(tempSensor)

	temperatureInCurrentUnit := util.ParseTempAccordingToUnit(unit, cpuTempFromSensor)

	measurement := SetMeasurementValues(temperatureInCurrentUnit, sensorID[0], deviceID)
	parsedData := util.ParseDataAccordingToFormat(format, measurement)

	tempData = append(tempData, parsedData)

	return tempData, nil
}

func getTempFromSensor(ctx context.Context) (float64, error) {
	sensorTeperatureInfo, err := host.SensorsTemperaturesWithContext(ctx)
	if err != nil {
		return 0, err
	}

	if len(sensorTeperatureInfo) == 0 {
		return ReadFileSystemFile("/sys/class/thermal/cooling_device2/device/status")
	}

	cpuTemp := sensorTeperatureInfo[0].Temperature

	return cpuTemp, nil
}
