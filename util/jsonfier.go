package util

import (
	"bufio"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/ranskills/c2m/setting"
)

// CreateJsonfy Transforms the line items to JSON and publishes to message broker
// TODO:
// 1. Consider renaming this to CreateJsonBuilder
// 2. Return an array of json string i.e. []string
func CreateJsonfy(cfg setting.Config) func(string) []map[string]string {
	dataFieldsMapping := cfg.BuildDataFieldsMapping()
	headingFieldsMapping := cfg.BuildHeadingFieldsMapping()
	nameFieldsMapping := cfg.BuildNameFieldsMapping()

	getFieldNameByIndex := func(index int) string {
		return dataFieldsMapping[strconv.Itoa(index)]
	}

	extractFieldsFromFileName := func(name string) map[string]string {
		fields := make(map[string]string)

		for _, mapping := range nameFieldsMapping {
			parts := strings.Split(name, mapping.Separator)

			fields[mapping.FieldName] = parts[mapping.FieldPosition-1]
		}
		return fields
	}

	copyToMap := func(from map[string]string, to map[string]string) {
		for k, v := range from {
			to[k] = v
		}
	}

	return func(filePath string) []map[string]string {

		stat, err := os.Stat(filePath)

		if err != nil {
			log.Println(err)
			return nil
		}

		notValidFile := stat.IsDir() || path.Ext(filePath) != cfg.File.Extension
		if notValidFile == true {
			log.Printf("Skipping file %s. Reason: Not a CSV file\n", stat.Name())
			return nil
		}

		log.Printf("Processing: %s\n", stat.Name())
		file, err := os.Open(filePath)

		//if err != nil {
		//	log.Println(err)
		//	return nil
		//}
		defer file.Close()

		var ret []map[string]string

		nameFields := extractFieldsFromFileName(stat.Name())
		headerFields := make(map[string]string)
		line := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line++
			lineText := scanner.Text()
			emptyLine := strings.TrimSpace(lineText) == ""

			if emptyLine {
				continue
			}

			if cfg.File.Data.HeaderPrefix != "" {
				headerDetected := strings.HasPrefix(lineText, cfg.File.Data.HeaderPrefix)
				if headerDetected {
					continue
				}
			}

			fields := strings.Split(lineText, ",")

			headingFieldsConfig := headingFieldsMapping[strconv.Itoa(line)]
			headerLineDetected := headingFieldsConfig != (setting.HeadingFieldConfig{})

			if headerLineDetected {
				value := strings.Split(lineText, headingFieldsConfig.Separator)[headingFieldsConfig.FieldPosition-1]
				headerFields[headingFieldsConfig.FieldName] = value
				continue
			}

			m := make(map[string]string)

			for i := 0; i < len(fields); i++ {
				fieldName := getFieldNameByIndex(i + 1)

				m[fieldName] = strings.TrimSpace(fields[i])
			}

			copyToMap(headerFields, m)
			copyToMap(nameFields, m)

			ret = append(ret, m)
		}

		return ret
	}
}
