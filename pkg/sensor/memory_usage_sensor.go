package sensor

import (
	"context"
	"strconv"

	"github.com/shirou/gopsutil/mem"
	"github.com/ttodorov/sensorcli/pkg/util"
)

const (
	memorySensor string = "MEMORY_USAGE"
)

type cpuMemorySensor Sensor

// CreateMemorySensor creates instance of memory sensor.
func CreateMemorySensor() ISensor {
	return &cpuMemorySensor{}
}

func (memoryS *cpuMemorySensor) GetSensorData(ctx context.Context, arguments ...string) ([]string, error) {
	memoryUsageData, err := getMemoryUsageData(ctx, arguments[1])
	if err != nil {
		return nil, err
	}

	return memoryUsageData, nil
}

func (memoryS *cpuMemorySensor) Validate(arguments ...string) error {
	return util.ValidateFormat(arguments[0])
}

func getMemoryUsageData(ctx context.Context, format string) ([]string, error) {

	var memoryData []string

	memoryUsageValues, err := getMemoryUsageValues(ctx)
	if err != nil {
		return nil, err
	}

	deviceID, err := getDeviceID()
	if err != nil {
		return nil, err
	}

	sensorID, err := GetSensorID(memorySensor)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(memoryUsageValues); i++ {
		memoryMeasurements := SetMeasurementValues(memoryUsageValues[i], sensorID[i], deviceID)
		memoryData = append(memoryData, util.ParseDataAccordingToFormat(format, memoryMeasurements))
	}

	return memoryData, nil
}

func getMemoryUsageValues(ctx context.Context) ([]string, error) {

	memory, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	totalMemory := strconv.FormatUint(memory.Total, 10)
	availableMemory := strconv.FormatUint(memory.Total, 10)
	usedMemory := strconv.FormatUint(memory.Used, 10)
	usedPercentMemory := strconv.FormatFloat(memory.UsedPercent, 'f', 2, 64)

	return []string{totalMemory, availableMemory, usedMemory, usedPercentMemory}, nil
}
