package util

import "strconv"

const (
	unitTypeFahrenheit string = "F"
	unitTypeCelsius    string = "C"
	celsiusName        string = "Celsius"
	fahrenheitName     string = "Fahrenheit"
)

func fahrenheit(temperature float64) float64 {
	return temperature*(9/5) + 32
}

// ParseTempAccordingToUnit returns correct parsed temperature.
func ParseTempAccordingToUnit(unit string, temperature float64) string {
	if unit == unitTypeFahrenheit {
		return strconv.FormatFloat(fahrenheit(temperature), 'f', 1, 64)
	}

	return strconv.FormatFloat(temperature, 'f', 1, 64)
}
