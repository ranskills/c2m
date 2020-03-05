package setting

import (
	"log"
	"strconv"
	"strings"
)

const (
	// HeadingLineIndex The specific line in the file where heading details are located
	HeadingLineIndex              = 0
	HeadingSeparatorIndex         = 1
	HeadingFieldNamePositionIndex = 2
	HeadingFieldNameIndex         = 3

	DataFieldNamePositionIndex = 0
	DataFieldNameIndex         = 1

	NameSeparatorIndex         = 0
	NameFieldNamePositionIndex = 1
	NameFieldNameIndex         = 2
)

// Config The configuration for the tool
type Config struct {
	Mqtt struct {
		Scheme   string `yaml:"scheme"`
		Port     string `yaml:"port"`
		Host     string `yaml:"host"`
		ClientID string `yaml:"clientId"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Topic    string `yaml:"topic"`
	} `yaml:"mqtt"`
	File struct {
		Extension string `yaml:"extension"`
		Heading   struct {
			Fields []string `yaml:"fields"`
		} `yaml:"heading"`
		Data struct {
			HeaderPrefix string   `yaml:"headerPrefix"`
			Fields       []string `yaml:"fields"`
		} `yaml:"data"`
		Name struct {
			Fields []string `yaml:"fields"`
		} `yaml:"name"`
	} `yaml:"file"`
}

// HeadingFieldConfig Heading field configurations
type HeadingFieldConfig struct {
	Separator     string
	FieldPosition int
	FieldName     string
}

// NameFieldConfig Heading field configurations
type NameFieldConfig struct {
	Separator     string
	FieldPosition int
	FieldName     string
}

// BuildDataFieldsMapping Comment
func (c Config) BuildDataFieldsMapping() map[string]string {
	mapping := make(map[string]string)

	for _, f := range c.File.Data.Fields {
		parts := strings.Split(f, ",")

		mapping[parts[DataFieldNamePositionIndex]] = strings.TrimSpace(parts[DataFieldNameIndex])
		// fmt.Println(parts[0], parts[1])
	}
	// fmt.Println(mapping)

	return mapping
}

func (c Config) BuildNameFieldsMapping() map[string]NameFieldConfig {
	mapping := make(map[string]NameFieldConfig)

	for _, f := range c.File.Name.Fields {
		parts := strings.Split(f, ",")
		//log.Println(f)
		//log.Println(parts)
		i, err := strconv.Atoi(parts[NameFieldNamePositionIndex])
		if err != nil {
			log.Fatal(err)
		}

		mapping[parts[HeadingLineIndex]] = NameFieldConfig{
			Separator:     parts[NameSeparatorIndex],
			FieldPosition: i,
			FieldName:     parts[NameFieldNameIndex],
		}
	}

	return mapping
}

// BuildHeadingFieldsMapping a
func (c Config) BuildHeadingFieldsMapping() map[string]HeadingFieldConfig {
	mapping := make(map[string]HeadingFieldConfig)

	for _, f := range c.File.Heading.Fields {
		parts := strings.Split(f, ",")

		i, err := strconv.Atoi(parts[HeadingFieldNamePositionIndex])

		if err != nil {
			log.Fatal(err)
		}

		mapping[parts[HeadingLineIndex]] = HeadingFieldConfig{
			Separator:     parts[HeadingSeparatorIndex],
			FieldPosition: i,
			FieldName:     parts[HeadingFieldNameIndex],
		}
		//fmt.Println(parts[0], parts[1])
	}
	// fmt.Println(mapping)

	return mapping
}
