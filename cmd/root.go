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
	"github.com/ttodorov/sensorcli/pkg/logger"
	"github.com/ttodorov/sensorcli/pkg/sensor"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var rootLogger logger.Logger = logger.NewLogger("./root")

var (
	cfgFile       string
	format        string
	deltaDuration int64
	sensorGroup   []string
	totalDuration float64
	file          string
	webHook       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sensorcli",
	Short: "Cli app which gets data from the sensors.",
	Long: `Cli app which gets cpu temperature data, 
	cpu usage data and memory usage data from the sensors of your local PC.`,

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

		err := sensor.ReadYamlFile("./model.yaml")
		if err != nil {
			rootLogger.Error(err)
			return err
		}

		err = terminateForTotalDuration(ctx)
		if err != nil {
			rootLogger.Error(err)
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&format, "format", "JSON", "The data could be printed either JSON or YAML format.")
	rootCmd.Flags().Int64Var(&deltaDuration, "delta_duration", 3, "The period of which you will get your sensor data.")
	rootCmd.Flags().StringVar(&file, "output_file", "", "Writing the output into CSV file.")
	rootCmd.Flags().Float64Var(&totalDuration, "total_duration", 60.0, "Terminating the whole program after specified duration")
	rootCmd.Flags().StringVar(&webHook, "web_hook_url", "", "Expose to current port.")

	rootCmd.Flags().StringSliceVar(&sensorGroup, "sensor_group", []string{""}, "There are three main sensor groups: CPU_TEMP, CPU_USAGE and MEMORY_USAGE.")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".sensorcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".sensorcli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func getDeltaDurationInSeconds() time.Duration {
	return time.Duration(deltaDuration) * time.Second
}

func getTotalDurationInSeconds() time.Duration {
	return time.Duration(totalDuration) * time.Second
}

func getSensorInfo(ctx context.Context, sensorGroup string) ([]string, error) {
	if sensorGroup == "" {
		rootLogger.Errorf("invalid sensor group")
		return nil, fmt.Errorf("invalid sensor group")
	}

	sensorType, err := sensor.CreateSensor(sensorGroup)
	if err != nil {
		return nil, err
	}

	err = sensorType.Validate(format)
	if err != nil {
		return nil, err
	}

	unit, err := sensor.GetSensorUnit(sensorGroup)
	if err != nil {
		rootLogger.Errorf(err.Error())
		return nil, err
	}

	sensorInfo, err := sensorType.GetSensorData(ctx, unit, format)
	if err != nil {
		rootLogger.Errorf(err.Error())
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
			rootLogger.Error(ctx.Err())
			return ctx.Err()
		case <-appTerminaitingDuration:
			return nil
		default:

			multipleSensorsData, err := getMultipleSensorsMeasurements(ctx)
			if err != nil {
				rootLogger.Error(err)
				return err
			}

			if file != "" {
				go sensor.WriteOutputToCSV(multipleSensorsData, file)
				rootLogger.Info("Writing sensor measurements in CSV file.")
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
	sensorsData := SendSensorData(sensorData, done)

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
			rootLogger.Error(ctx.Err())
			return ctx.Err()
		}
	}

}

func webHookURL(url string, data string) {
	var json = []byte(data)
	http.Post(url, "application/json", bytes.NewBuffer(json))
}

// SendSensorData ...
func SendSensorData(sensorsInfo []string, done chan bool) <-chan string {
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
