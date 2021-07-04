/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/ttodorov/sensorcli/pkg/sensor"
)

var (
	format        string
	deltaDuration int64
	sensorGroup   []string
	totalDuration float64
	file          string
	webHook       string
)

// cpuCmd represents the cpu command
var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		interuptSignal := make(chan os.Signal, 1)

		signal.Notify(interuptSignal, os.Interrupt)
		defer func() {
			signal.Stop(interuptSignal)
			cancel()
		}()

		go func() {
			select {
			case <-interuptSignal:
				cancel()
			case <-ctx.Done():
			}
		}()

		err := terminateForTotalDuration(ctx)
		if err != nil {
			cmdLogger.Error(err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpuCmd)

	cpuCmd.Flags().StringVar(&format, "format", "JSON", "The data could be printed either JSON or YAML format.")
	cpuCmd.Flags().Int64Var(&deltaDuration, "delta_duration", 3, "The period of which you will get your sensor data.")
	cpuCmd.Flags().StringVar(&file, "output_file", "", "Writing the output into CSV file.")
	cpuCmd.Flags().Float64Var(&totalDuration, "total_duration", 60.0, "Terminating the whole program after specified duration")
	cpuCmd.Flags().StringVar(&webHook, "web_hook_url", "", "Expose to current port.")

	cpuCmd.Flags().StringSliceVar(&sensorGroup, "sensor_group", []string{""}, "There are three main sensor groups: CPU_TEMP, CPU_USAGE and MEMORY_USAGE.")
}

func getDeltaDurationInSeconds() time.Duration {
	return time.Duration(deltaDuration) * time.Second
}

func getTotalDurationInSeconds() time.Duration {
	return time.Duration(totalDuration) * time.Second
}

func getSensorInfo(ctx context.Context, sensorGroup string) ([]string, error) {
	if sensorGroup == "" {
		cmdLogger.Errorf("invalid sensor group")
		return nil, fmt.Errorf("invalid sensor group")
	}

	sensorType, err := sensor.NewSensor(sensorGroup)
	if err != nil {
		return nil, err
	}

	err = sensorType.Validate(format)
	if err != nil {
		return nil, err
	}

	unit, err := sensor.GetSensorUnits(sensorGroup)
	if err != nil {
		cmdLogger.Errorf(err.Error())
		return nil, err
	}

	fmt.Println(unit)

	sensorInfo, err := sensorType.GetSensorData(ctx, unit, format)
	if err != nil {
		cmdLogger.Errorf(err.Error())
		return nil, err
	}

	return sensorInfo, nil
}

func getMultipleSensorsMeasurements(ctx context.Context) ([]string, error) {
	var multipleSensorsData []string

	for i := 0; i < len(sensorGroup); i++ {

		var currentSensorGroupData []string

		currentSensorGroupData, err := getSensorInfo(ctx, sensorGroup[i])
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(currentSensorGroupData); j++ {
			multipleSensorsData = append(multipleSensorsData, currentSensorGroupData[j])
		}

	}

	return multipleSensorsData, nil
}

func terminateForTotalDuration(ctx context.Context) error {

	appTerminaitingDuration := time.After(getTotalDurationInSeconds())

	for {
		select {
		case <-ctx.Done():
			cmdLogger.Error(ctx.Err())
			return ctx.Err()
		case <-appTerminaitingDuration:
			return nil
		default:

			multipleSensorsData, err := getMultipleSensorsMeasurements(ctx)
			if err != nil {
				cmdLogger.Error(err)
				return err
			}

			if file != "" {
				go sensor.WriteOutputToCSV(multipleSensorsData, file)
				cmdLogger.Info("Writing sensor measurements in CSV file.")
			}

			err = getMeasurementsInDeltaDuration(ctx, multipleSensorsData)
			if err != nil {
				return err
			}
		}
	}

}

func getMeasurementsInDeltaDuration(ctx context.Context, sensorData []string) error {
	measurementDuration := time.After(getDeltaDurationInSeconds())
	done := make(chan bool)
	sensorsData := sendSensorData(sensorData, done)

	for {
		select {
		case data := <-sensorsData:

			if webHook != "" {
				webHookURL(webHook, data)
			}

			fmt.Println(data)
		case <-measurementDuration:
			done <- true
			return nil
		case <-ctx.Done():
			done <- true
			cmdLogger.Error(ctx.Err())
			return ctx.Err()
		}
	}

}

//TODO integrate with http server
func webHookURL(url string, data string) {
	var json = []byte(data)
	http.Post(url, "application/json", bytes.NewBuffer(json))
}

func sendSensorData(sensorsInfo []string, done chan bool) <-chan string {
	out := make(chan string)

	go func() {
		for _, currentSensorInfo := range sensorsInfo {
			out <- currentSensorInfo
		}

		if <-done {
			close(out)
			return
		}
	}()

	return out
}
