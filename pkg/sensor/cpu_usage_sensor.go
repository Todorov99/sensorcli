package sensor

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/ttodorov/sensorcli/pkg/util"
)

const (
	cpuUsageName        string = "cpuUsagePercent"
	cpuCoresCountName   string = "cpuCoresCountName"
	cpuFrequencyMHzName string = "cpuFrequencyMHz"

	cpuUsageUnit     string = "%"
	cpuCoresUnit     string = "count"
	cpuFrequencyUnit string = "MHz"

	coresDescription     string = "CPU cores count"
	usageDescription     string = "CPU usage percent"
	frequencyDescription string = "CPU frequency MHz"

	invalidFlagError string = "Invalid flag."

	usageSensor string = "CPU_USAGE"
)

type cpuUsageSensor Measurment

// CreateUsageSensor creates instance of usage sensor.
func CreateUsageSensor() ISensor {
	return &cpuUsageSensor{}
}

func (usageS *cpuUsageSensor) GetSensorData(arguments ...string) ([]string, error) {

	cpuUsage, err := getUsageMeasurements(arguments[1])

	if err != nil {
		return nil, err
	}

	return cpuUsage, nil
}

func (usageS *cpuUsageSensor) Validate(arguments ...string) error {
	return util.ValidateFormat(arguments[0])
}

func getUsageMeasurements(format string) ([]string, error) {
	var usageData []string

	deviceID, err := getDeviceID()
	if err != nil {
		return nil, err
	}

	sensorID, err := GetSensorID(usageSensor)
	if err != nil {
		return nil, err
	}

	cpuInfo, err := getCPUUsageInfo()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cpuInfo); i++ {
		usageMeasurements := SetMeasurementValues(cpuInfo[i], sensorID[i], deviceID)
		usageData = append(usageData, util.ParseDataAccordingToFormat(format, usageMeasurements))
	}

	return usageData, nil
}

func getCPUCoresAndFrequency() (string, string, error) {
	var cpuInfo []cpu.InfoStat
	var err error

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		cpuInfo, err = darwinArm64Info(context.Background())
		if err != nil {
			return "", "", fmt.Errorf("error getting cpu cores and frequency on darwin arm64: %w", err)
		}

		return strconv.FormatInt(int64(cpuInfo[0].Cores), 10), strconv.FormatFloat(cpuInfo[0].Mhz, 'f', 2, 64), nil
	}

	cpuInfo, err = cpu.Info()
	if err != nil {
		return "", "", fmt.Errorf("error with getting cpu cores and frequency: %w", err)
	}

	return strconv.FormatInt(int64(cpuInfo[0].Cores), 10), strconv.FormatFloat(cpuInfo[0].Mhz, 'f', 2, 64), nil
}

func getUsedPercent() (string, error) {
	cpuUsedPercentage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", fmt.Errorf("error in getting used cpu percent")
	}

	usedPercentage := cpuUsedPercentage[0] / 100

	return strconv.FormatFloat(usedPercentage, 'f', 2, 64), nil
}

func getCPUUsageInfo() ([]string, error) {
	cpuCores, cpuFrequency, err := getCPUCoresAndFrequency()
	if err != nil {
		return nil, err
	}

	cpuUsage, err := getUsedPercent()
	if err != nil {
		return nil, err
	}

	cpuInfo := []string{cpuUsage, cpuCores, cpuFrequency}

	return cpuInfo, nil
}
