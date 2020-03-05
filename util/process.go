package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/ranskills/mp/setting"
	"github.com/ranskills/mp/broker"
)

// CreateProcessFileFunc Transforms the line items to JSON and publishes to message broker
func CreateProcessFileFunc(cfg setting.Config, prettyPrintJson bool) func(string) {
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

			fields[mapping.FieldName] = parts[mapping.FieldPosition - 1]
		}
		return fields
	}

	client := broker.GetClient(cfg)

	return func(filePath string) {
		log.Printf("\n\nProcessing: %s\n", filePath)

		stat, err := os.Stat(filePath)

		if err != nil {
			log.Fatal(err)
		}

		notValidFile := stat.IsDir() || path.Ext(filePath) != cfg.File.Extension
		if notValidFile {
			log.Printf("Cannot process the filePath: %s\n", filePath)
			return
		}

		file, err := os.Open(filePath)

		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// return
		nameFields := extractFieldsFromFileName(stat.Name())
		//fmt.Printf("**** %s\n", nameFields)
		headerFields := make(map[string]string)
		line := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line++
			lineText := scanner.Text()
			//log.Println("xx", lineText, strings.TrimSpace(lineText) == "")
			emptyLine := strings.TrimSpace(lineText) == ""

			if emptyLine {
				continue
			}

			if cfg.File.Data.HeaderPrefix != "" {
				headerDetected := strings.HasPrefix(lineText, "test_number, test_name")
				if headerDetected {
					continue
				}
			}

			fields := strings.Split(lineText, ",")

			headingFieldsConfig := headingFieldsMapping[strconv.Itoa(line)]
			headerLineDetected := headingFieldsConfig != (setting.HeadingFieldConfig{})

			if headerLineDetected {
				value := strings.Split(lineText, headingFieldsConfig.Separator)[headingFieldsConfig.FieldPosition-1]
				//fmt.Printf("Configured header settings. %s: %s\n", headingFieldsConfig.FieldName, value)
				headerFields[headingFieldsConfig.FieldName] = value
				continue
			}

			// fmt.Println(len(fields), fields[0])

			j := make(map[string]string)

			for i := 0; i < len(fields); i++ {
				fieldName := getFieldNameByIndex(i + 1)

				j[fieldName] = strings.TrimSpace(fields[i])
			}

			// Should be conditional
			//if true {
				for fieldName, value := range headerFields {
					j[fieldName] = value
				}
				for fieldName, value := range nameFields {
					j[fieldName] = value
				}
			//}


			var payload []byte
			if prettyPrintJson {
				payload, err = json.MarshalIndent(j, "", "\t")
			} else {
				payload, err = json.Marshal(j)
			}

			//

			if err == nil {
				fmt.Println(string(payload))
				//client.Publish()
				client.Publish(cfg.Mqtt.Topic, 0, false, string(payload))
			} else {
				fmt.Println(err)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}

}
