package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ranskills/c2m/broker"
	"github.com/ranskills/c2m/setting"
	"github.com/ranskills/c2m/util"
)

type ProcessedFileInfo struct {
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	NumRecords int       `json:"numRecords"`
}

// CreateRunAction Creates a run action
func CreateRunAction(cfg setting.Config) func(args []string, options map[string]string) int {
	jsonfy := util.CreateJsonfy(cfg)
	client := broker.GetClient(cfg)

	processPayloads := func(payloads []map[string]string, prettyPrintJson bool) {
		var payload []byte
		var err error

		for _, j := range payloads {

			if prettyPrintJson {
				payload, err = json.MarshalIndent(j, "", "\t")
			} else {
				payload, err = json.Marshal(j)
			}

			if err == nil {
				token := client.Publish(cfg.Mqtt.Topic, 0, false, string(payload))
				token.Wait()
			} else {
				log.Println(err)
			}
		}
	}

	return func(args []string, options map[string]string) int {

		prettyPrintJson, _ := strconv.ParseBool(options["pretty"])
		watch, _ := strconv.ParseBool(options["watch"])

		srcDir := options["src"]
		if !strings.HasSuffix(srcDir, string(os.PathSeparator)) {
			srcDir += string(os.PathSeparator)
		}

		files, err := ioutil.ReadDir(srcDir)

		if err != nil {
			log.Fatal(err)
		}

		processFile := func(filePath string) {
			const trackingFilename string = ".c2m"

			parts := strings.Split(filePath, string(os.PathSeparator))
			filename := parts[len(parts)-1]

			contents, _ := ioutil.ReadFile(trackingFilename)
			trackingMap := make(map[string]ProcessedFileInfo)
			json.Unmarshal(contents, &trackingMap)

			if _, ok := trackingMap[filename]; ok {
				log.Printf("Skipping already processed file %s\n", filename)
				return
			}

			startTime := time.Now()
			jsons := jsonfy(filePath)
			processPayloads(jsons, prettyPrintJson)
			endTime := time.Now()

			if numRecords := len(jsons); numRecords > 0 {
				trackingMap[filename] = ProcessedFileInfo{startTime, endTime, numRecords}
			}

			data, _ := json.MarshalIndent(trackingMap, "", "\t")
			ioutil.WriteFile(trackingFilename, data, os.ModePerm)
		}

		for _, file := range files {
			processFile(srcDir + file.Name())
		}

		if watch {
			util.WatchDirectory(srcDir, processFile)
		}

		time.Sleep(5 * time.Second)

		return 0
	}
}
