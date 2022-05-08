# Sensor CLI

 Sensor CLI is an application that gets hardware measurements for a specific predifined group of sensors

## Description

This application run on local OS and could use flags for sending the collected measurements to Web Server.
The application suports measurements from multiple sensor groups.

- Available groups:
   - CPU_TEMP
   - CPU_USAGE
   - MEMORY_USAGE

- Base CLI usage description:
```
CLI application which gets data from the predefined sensor groups from the hardware where it operates.

Usage:
  sensorcli [flags]
  sensorcli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  cpu         Start collecting hardware metrics
  help        Help about any command

Flags:
  -h, --help   help for sensorcli

Use "sensorcli [command] --help" for more information about a command.
```

- Detailed `cpu` command flag information:
```
Start collecting hardware metrics

Usage:
  sensorcli cpu [flags]

Flags:
      --configDirPath string   The path to the configuration file for the device for which the hardware measurements will be collected
      --delta_duration int     The period between each measurement (default 3)
      --email string           The email to which the final result should be send
      --format string          The data could be printed either JSON or YAML format (default "JSON")
      --generateReport         Flag that shows whether to generate report file
  -h, --help                   help for cpu
      --mail_hook_url string   Base URL to the mailsender API
      --password string        The password of the user used for remote authentication to the sensor API
      --reportType string      The type of the report file that has to be generated. Possible values xlsx, csv. (default "xlsx")
      --rootCAPath string      The path to the root Certificate Authoritate (CA) used for the TLS client config. If no CA is provided the verification of the certificates is skipped
      --sensor_group strings   There are three main sensor groups: CPU_TEMP, CPU_USAGE and MEMORY_USAGE. Each senosor group could have system file that will hold specific information
      --total_duration float   The period in which the hardware measurements will be collected (default 60)
      --username string        The username of the user used for remote authentication to the sensor API
      --web_hook_url string    Base URL to the sensor API
```

### Example device config

- This config could be used in case you want to collect measurements on your machine without sending them the the Web Server. To be able to specify the directory where the config is use configDirPath. Make sure that the file is called `device_cfg.yaml`:

```
id: 1
name: device_name
description: my laptop device
sensors:
- id: 1
  name: cpuTemperature
  description: Measures CPU temperature in provided unit
  unit: C
  sensorGroups: CPU_TEMP
- id: 2
  name: cpuUsagePercentage
  description: Measures CPU usage in percentages
  unit: '%'
  sensorGroups: CPU_USAGE
- id: 3
  name: cpuCores
  description: Gets the number of CPU cores
  unit: count
  sensorGroups: CPU_USAGE
- id: 4
  name: cpuFrequency
  description: Measures CPU frequency in a provided unit
  unit: GHz
  sensorGroups: CPU_USAGE
- id: 5
  name: memoryTotal
  description: Measures memory total RAM
  unit: GigaBytes
  sensorGroups: MEMORY_USAGE
- id: 6
  name: memoryAvailable
  description: Gets the available RAM in a provided unit
  unit: GigaBytes
  sensorGroups: MEMORY_USAGE
- id: 7
  name: memoryUsed
  description: Gets the used RAM from the programs in a provided unit
  unit: GigaBytes
  sensorGroups: MEMORY_USAGE
- id: 8
  name: memoryUsedPercentage
  description: Used percentage RAM from the programs
  unit: '%'
  sensorGroups: MEMORY_USAGE

```

### Example commands

- Binaries for different OS could be find in the `Releases` section in Github.

- Start measurements with report and saving in Web Server database with configuration fetched from the server. When you are inside the directory got after unpacking the archive received from the API run:

Example:

```
./sensorcli_darwin cpu --sensor_group CPU_TEMP --sensor_group=CPU_USAGE --sensor_group=MEMORY_USAGE --total_duration=20 --web_hook_url='https://localhost:8081' --username='ttodorov' --password='Abcd123!@' --generateReport=true --reportType=xlsx --mail_hook_url='https://localhost:8083' --email='todor.mtodorov01@gmail.com' --configDirPath=.
```

- Start separatelly from the Web server with report generation and mail notification. Create your device configuration described in the section above and then run:

Example:

```
./sensorcli_darwin cpu --sensor_group CPU_TEMP --sensor_group=CPU_USAGE --sensor_group=MEMORY_USAGE --total_duration=20 --generateReport=true --reportType=xlsx --mail_hook_url='https://localhost:8083' --email='todor.mtodorov01@gmail.com' --configDirPath=.
```

- Start basic measurement colleaction with additional file for getting the temperature in case there is not any installed driver on the machine:

Example:

```
./sensorcli_darwin cpu --sensor_group CPU_TEMP=thermal --sensor_group=CPU_USAGE --sensor_group=MEMORY_USAGE --total_duration=20 --generateReport=true --reportType=xlsx --configDirPath=.
```

### Installing

```
go build -o sensorcli ./

```

## Use as go mod
```
go get -u github.com/Todorov99/sensorcli
```