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
	usageSensor string = "CPU_USAGE"
)

type cpuUsageSensor struct {
	cpuUsage     string
	cpuCores     string
	cpuFrequency string
	deviceID     string
	sensors      []sensor
}

// CreateUsageSensor creates instance of usage sensor.
func CreateUsageSensor() ISensor {
	return &cpuUsageSensor{}
}

func (usageS *cpuUsageSensor) GetSensorData(ctx context.Context, unit []string, format string) ([]string, error) {
	cpuUsage, err := getUsageMeasurements(ctx, format)

	if err != nil {
		return nil, err
	}

	return cpuUsage, nil
}

func (usageS *cpuUsageSensor) Validate(arguments ...string) error {
	return util.ValidateFormat(arguments[0])
}

func getUsageMeasurements(ctx context.Context, format string) ([]string, error) {
	var usageData []string

	deviceID, err := devices.getDeviceID()
	if err != nil {
		return nil, err
	}

	sensors, err := devices.getSensorsByGroup(usageSensor)
	if err != nil {
		return nil, err
	}

	cpuInfo, err := getCPUUsageInfo(ctx)
	if err != nil {
		return nil, err
	}

	cpuInfo.deviceID = deviceID
	cpuInfo.sensors = sensors

	measurements := newMeasurements(cpuInfo)

	for _, m := range measurements {
		usageData = append(usageData, util.ParseDataAccordingToFormat(format, m))
	}

	return usageData, nil
}

func getCPUCoresAndFrequency(ctx context.Context) (string, string, error) {
	var cpuInfo []cpu.InfoStat
	var err error

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		cpuInfo, err = darwinArm64InfoWithContext(ctx)
		if err != nil {
			return "", "", fmt.Errorf("error getting cpu cores and frequency on darwin arm64: %w", err)
		}

		return strconv.FormatInt(int64(cpuInfo[0].Cores), 10), strconv.FormatFloat(cpuInfo[0].Mhz, 'f', 2, 64), nil
	}

	cpuInfo, err = cpu.InfoWithContext(ctx)
	if err != nil {
		return "", "", fmt.Errorf("error with getting cpu cores and frequency: %w", err)
	}

	return strconv.FormatInt(int64(cpuInfo[0].Cores), 10), strconv.FormatFloat(cpuInfo[0].Mhz, 'f', 2, 64), nil
}

func getUsedPercent(ctx context.Context) (string, error) {
	cpuUsedPercentage, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return "", fmt.Errorf("error in getting used cpu percent")
	}

	usedPercentage := cpuUsedPercentage[0] / 100

	return strconv.FormatFloat(usedPercentage, 'f', 2, 64), nil
}

func getCPUUsageInfo(ctx context.Context) (cpuUsageSensor, error) {
	cores, frequency, err := getCPUCoresAndFrequency(ctx)
	if err != nil {
		return cpuUsageSensor{}, err
	}

	usage, err := getUsedPercent(ctx)
	if err != nil {
		return cpuUsageSensor{}, err
	}

	return cpuUsageSensor{
		cpuCores:     cores,
		cpuFrequency: frequency,
		cpuUsage:     usage,
	}, nil
}
