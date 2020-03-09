package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ranskills/c2m/setting"
)

func getConfig(t *testing.T) setting.Config {
	var cfg setting.Config

	err := cleanenv.ReadConfig("config.yml", &cfg)

	if err != nil {
		t.Fatal("Failed loading config file")
	}
	return cfg
}

func TestBuildHeadingFieldsMapping(t *testing.T) {
	cfg := getConfig(t)

	expected := map[string]setting.HeadingFieldConfig{
		"1": setting.HeadingFieldConfig{
			Separator:     " ",
			FieldPosition: 3,
			FieldName:     "serialNumber",
		},
		"2": setting.HeadingFieldConfig{
			Separator:     " ",
			FieldPosition: 3,
			FieldName:     "partNumber",
		},
	}

	actual := cfg.BuildHeadingFieldsMapping()
	t.Log(actual)

	assert.Equal(t, expected, actual)
}

func TestBuildNameFieldsMapping(t *testing.T) {
	cfg := getConfig(t)

	// TODO: Change key either be an auto number or the field name (I'll favour the later)
	expected := map[string]setting.NameFieldConfig{
		" ": setting.NameFieldConfig{
			Separator:     " ",
			FieldPosition: 1,
			FieldName:     "timeStamp",
		},
	}

	actual := cfg.BuildNameFieldsMapping()
	t.Log(actual)

	assert.Equal(t, expected, actual)
}

func TestBuildDataFieldsMapping(t *testing.T) {
	cfg := getConfig(t)

	expected := map[string]string{
		"1": "testNumber",
		"2": "testName",
		"3": "lowerLimit",
		"4": "upperLimit",
		"5": "measuredValue",
		"6": "status",
	}

	actual := cfg.BuildDataFieldsMapping()
	t.Log(actual)

	assert.Equal(t, expected, actual)
}
