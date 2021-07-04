package sensor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/mem"
	"github.com/ttodorov/sensorcli/pkg/util"
)

const (
	memorySensor string = "MEMORY_USAGE"
)

type cpuMemorySensor struct {
	totalMemory       string
	availableMemory   string
	usedMemory        string
	usedPercentMemory string
	deviceID          string
	sensors           []sensor
}

// CreateMemorySensor creates instance of memory sensor.
func CreateMemorySensor() ISensor {
	return &cpuMemorySensor{}
}

func (memoryS *cpuMemorySensor) GetSensorData(ctx context.Context, unit []string, format string) ([]string, error) {
	sensorLogger.Info("Gerring sensor data...")
	memoryUsageData, err := getMemoryUsageData(ctx, format)
	if err != nil {
		msg := "failed to get memory usage data: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	return memoryUsageData, nil
}

func (memoryS *cpuMemorySensor) Validate(arguments ...string) error {
	return util.ValidateFormat(arguments[0])
}

func getMemoryUsageData(ctx context.Context, format string) ([]string, error) {
	sensorLogger.Info("Getting memory usage data...")
	var memoryData []string

	deviceID, err := devices.getDeviceID()
	if err != nil {
		msg := "failed to get deviceID: %w"
		return nil, fmt.Errorf(msg, err)
	}

	sensors, err := devices.getSensorsByGroup(memorySensor)
	if err != nil {
		msg := "failed to get sensorID: %w"
		return nil, fmt.Errorf(msg, err)
	}

	memoryUsageValues, err := getMemoryUsageValues(ctx)
	if err != nil {
		msg := "failed to get memory usage data: %w"
		sensorLogger.Errorf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	memoryUsageValues.sensors = sensors
	memoryUsageValues.deviceID = deviceID

	measurements := newMeasurements(memoryUsageValues)
	for _, m := range measurements {
		memoryData = append(memoryData, util.ParseDataAccordingToFormat(format, m))
	}

	return memoryData, nil
}

func getMemoryUsageValues(ctx context.Context) (cpuMemorySensor, error) {
	sensorLogger.Info("Getting memory usage data...")
	memory, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		msg := "failed to get virtual memory: %w"
		sensorLogger.Errorf(msg, err)
		return cpuMemorySensor{}, fmt.Errorf(msg, err)
	}

	return cpuMemorySensor{
		totalMemory:       strconv.FormatUint(memory.Total, 10),
		availableMemory:   strconv.FormatUint(memory.Total, 10),
		usedMemory:        strconv.FormatUint(memory.Used, 10),
		usedPercentMemory: strconv.FormatFloat(memory.UsedPercent, 'f', 2, 64),
	}, nil
}
