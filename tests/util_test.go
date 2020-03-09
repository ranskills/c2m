package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ranskills/c2m/setting"
	"github.com/ranskills/c2m/util"
	"github.com/stretchr/testify/assert"
)

func TestGetOptionValue(t *testing.T) {

	//var expected string
	var (
		expected string
		value    string
		err      error
	)

	expected = "config-test.yml"
	value, err = util.GetOptionValue([]string{"-c", expected}, "-c", "--config")
	assert.Equal(t, expected, value, fmt.Sprintf("Configuration file not defaulting to to %s", expected))

	value, err = util.GetOptionValue([]string{"-c"}, "-c", "--config")
	assert.Error(t, err, "Error not returned when no value is provided for the option")

	// Testing long option
	value, err = util.GetOptionValue([]string{"--config=" + expected}, "-c", "--config")
	assert.Equal(t, expected, value, fmt.Sprintf("Configuration file not defaulting to to %s", expected))

	value, err = util.GetOptionValue([]string{"--config"}, "-c", "--config")
	assert.Error(t, err, "Error not returned when no value is provided for the option")

	value, err = util.GetOptionValue([]string{}, "-c", "--config")
	assert.Error(t, err, "Option not present")

	value, err = util.GetOptionValue([]string{"--pretty", "--watch"}, "-c", "--config")
	assert.Error(t, err, "Option not present")
}

func TestGetConfigOptionValue(t *testing.T) {

	const defaultValue = "./config.yml"
	var providedConfigFile string

	value := util.GetConfigOptionValue([]string{})
	assert.Equal(t, defaultValue, value, fmt.Sprintf("Configuration file not defaulting to to %s", defaultValue))

	providedConfigFile = "config-labview.yml"
	value = util.GetConfigOptionValue([]string{"-c", providedConfigFile})
	assert.Equal(t, providedConfigFile, value, fmt.Sprintf("The passed configuration file %s not returned, but instead gave %s", providedConfigFile, value))

	providedConfigFile = "config-labview.yml"
	value = util.GetConfigOptionValue([]string{"--config=" + providedConfigFile})
	assert.Equal(t, providedConfigFile, value, fmt.Sprintf("The passed configuration file %s not returned, but instead gave %s", providedConfigFile, value))
}

func TestJsonfy(t *testing.T) {
	var config setting.Config

	cleanenv.ReadConfig("config.yml", &config)

	expected := []map[string]string{
		map[string]string{
			"lowerLimit":    "0.00",
			"measuredValue": "1.00",
			"partNumber":    "01-001371e9bs",
			"serialNumber":  "72582",
			"status":        "PASS",
			"testName":      "Power Up",
			"testNumber":    "1",
			"timeStamp":     "2020-2-5_4-49-30",
			"upperLimit":    "1.00",
		},
		map[string]string{
			"lowerLimit":    "173.00",
			"measuredValue": "234.65",
			"partNumber":    "01-001371e9bs",
			"serialNumber":  "72582",
			"status":        "PASS",
			"testName":      "Main Blower at 100 % L1-L2",
			"testNumber":    "2",
			"timeStamp":     "2020-2-5_4-49-30",
			"upperLimit":    "260.00",
		},
		map[string]string{
			"serialNumber":  "72582",
			"partNumber":    "01-001371e9bs",
			"testNumber":    "11",
			"testName":      "Drum Agitator Clockwise at 100% Speed L1-L2",
			"lowerLimit":    "173.00",
			"upperLimit":    "260.00",
			"measuredValue": "235.54",
			"status":        "PASS",
			"timeStamp":     "2020-2-5_4-49-30",
		},
		map[string]string{
			"serialNumber":  "72582",
			"partNumber":    "01-001371e9bs",
			"testNumber":    "61",
			"testName":      "End of Tests",
			"lowerLimit":    "0.00",
			"upperLimit":    "1.00",
			"measuredValue": "1.00",
			"status":        "PASS",
			"timeStamp":     "2020-2-5_4-49-30",
		},
	}

	jsonfy := util.CreateJsonfy(config)
	jsons := jsonfy("./files/d/2020-2-5_4-49-30 PM_72582.csv")
	log.Println(jsons)
	assert.Equal(t, expected, jsons)

	jsons = jsonfy("./files/d/non-existent-file.csv")
	assert.Nil(t, jsons, "No JSON data should be returned for a non-existent file")

	jsons = jsonfy("./files/d/hello.txt")
	assert.Nil(t, jsons, "No JSON data should be returned for a non-csv file")

}
