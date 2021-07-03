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
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/ttodorov/sensorcli/pkg/logger"
	"github.com/ttodorov/sensorcli/pkg/sensor"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var rootLogger logger.Logger = logger.NewLogger("./cmd/root")

var (
	cfgFile       string
	format        string
	deltaDuration int64
	sensorGroup   []string
	totalDuration float64
	file          string
	webHook       string
)

const (
	cpuTemp        string = "CPU_TEMP"
	cpuUsage       string = "CPU_USAGE"
	memoryUsage    string = "MEMORY_USAGE"
	outputToFile   string = "output_file"
	invalidCommand string = "Invalid command."
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sensorcli",
	Short: "Cli app which gets data from the sensors.",
	Long:  `Cli app which gets cpu temperature data, cpu usage data and memory usage data from the sensors of your local PC.`,

	Run: func(cmd *cobra.Command, args []string) {
		terminateForTotalDuration()
		for {

		}
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

func getSensorInfo(sensorGroup string) ([]string, error) {
	if sensorGroup == "" {
		rootLogger.Errorf("invalid sensor group")
		return nil, fmt.Errorf("invalid sensor group")
	}

	sensorType, err := sensor.CreateSensor(sensorGroup)

	err = sensorType.Validate(format)

	unit, err := sensor.GetTempSensorUnit(sensorGroup)
	if err != nil {
		rootLogger.Errorf(err.Error())
		return nil, err
	}

	sensorInfo, err := sensorType.GetSensorData(unit, format)
	if err != nil {
		rootLogger.Errorf(err.Error())
		return nil, err
	}

	return sensorInfo, nil
}

func getMultipleSensorsMeasurements() ([]string, error) {
	var multipleSensorsData []string

	for i := 0; i < len(sensorGroup); i++ {

		var currentSensorGroupData []string

		currentSensorGroupData, err := getSensorInfo(sensorGroup[i])
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(currentSensorGroupData); j++ {
			multipleSensorsData = append(multipleSensorsData, currentSensorGroupData[j])
		}

	}

	return multipleSensorsData, nil
}

func terminateForTotalDuration() {

	appTerminaitingDuration := time.After(getTotalDurationInSeconds())

	for {
		select {
		case <-appTerminaitingDuration:
			return
		default:

			multipleSensorsData, err := getMultipleSensorsMeasurements()
			if err != nil {
				rootLogger.Error(err)
				panic(err)
			}

			if file != "" {
				go sensor.WriteOutputToCSV(multipleSensorsData, file)
				rootLogger.Info("Writing sensor measurements in CSV file.")
			}

			getMeasurementsInDeltaDuration(multipleSensorsData)
		}
	}

}

func getMeasurementsInDeltaDuration(sensorData []string) {

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
			return
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
		for _, i := range sensorsInfo {
			out <- i
		}

		if <-done {
			close(out)
			return
		}
	}()

	return out
}
