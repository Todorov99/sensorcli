package sensor

import (
	"github.com/shirou/gopsutil/host"
	"github.com/ttodorov/sensorcli/pkg/util"
)

const (
	celsius     string = "C"
	fahrenheit  string = "F"
	json        string = "JSON"
	yamlStr     string = "YAML"
	cpuTempName string = "cpuTemp"
	unitError   string = "Invalid unit."
	formatError string = "Invalid format."

	tempSensor string = "CPU_TEMP"
)

type cpuTempSensor Measurment

// CreateTempSensor creates instance of temperature sensor.
func CreateTempSensor() ISensor {
	return &cpuTempSensor{}
}

func (tempS *cpuTempSensor) GetSensorData(arguments ...string) ([]string, error) {
	cpuTemp, err := getTempMeasurments(arguments[0], arguments[1])
	if err != nil {
		return nil, err
	}

	return cpuTemp, nil
}

func (tempS *cpuTempSensor) Validate(arguments ...string) error {
	return util.ValidateFormat(arguments[0])
}

func getTempMeasurments(unit string, format string) ([]string, error) {
	var tempData []string

	cpuTempFromSensor, err := getTempFromSensor()
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

func getTempFromSensor() (float64, error) {
	sensorTeperatureInfo, err := host.SensorsTemperatures()
	if err != nil {
		return 0, err
	}

	if len(sensorTeperatureInfo) == 0 {
		return ReadFileSystemFile()
	}

	cpuTemp := sensorTeperatureInfo[0].Temperature

	return cpuTemp, nil
}
