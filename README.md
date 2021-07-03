### Sensor CLI

---

 Sensor CLI application, which gets the CPU temperature, CPU usage and memory usage measurements of your local system.

## Getting Started

This application run on local OS and in a Docker container.
The application suports measurements from multiple sensor groups.
Sensor and device information is parsed from previously specified file. 
Sensor measurement data could be exported in JSON and YAML file.

### Installing

Type go install sensor for installing sensor cli application.

## Usage
    sensorcli [flags]

### Flags

    --delta_duration   The period between two sensor measurements you will get. The default duration is 3 seconds.
    --total_duration  Total period of terminating the program. The default total duration is 60 seconds.
    --format   The data could be printed either JSON or YAML (default "JSON").
    --help             help for sensor
    --sensor_group     There are three main sensor groups: CPU_TEMP, CPU_USAGE, MEMORY_USAGE.
    --output_file Writing sensor measurements in CSV file with specific dilimiter.
    --web_hook_url creates post requets to passed url with sensor measurements.

### Flags options

    These are flags options for the sensor cli application. 

    --sensor_group=<CPU_TEMP, CPU_USAGE, MEMORY_USAGE>
    --delta_duration=<seconds>
    --total_duration=<seconds>
    --forma=<JSON, YAML>
    --output_file=<fileName>
    --web_hook_url=<URL>
    --help


### Flags options description

    1. Types
        CPU_TEMP - gets the temperature from sensor at current time and the value in the chosen unit.
        CPU_USAGE - gets the cpu usage data from sensror at current time for cores, frequency and used percent.
        MEMORY_USAGE - gets available, total, used and used percent memory at current time.

    2. Formats
        - JSON - gets data in JSON format
        - YAML - gets data in YAML format

### Run Unit tests with Ginkgo

To run all unit tests successfully first change fileFullPath variable with the full path of model.yaml file from the main sensor directory.

    1. Run all unit tests with coverage of each package:
        -   ginkgo -r -coverprofile coverage.out
        
    2. Run unit test for current package:
        -   ginkgo <package realtive path> - For example: ginkgo ./utils
