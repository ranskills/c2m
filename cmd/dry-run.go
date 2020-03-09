package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ranskills/c2m/setting"
	"github.com/ranskills/c2m/util"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func CreateDryRunAction(cfg setting.Config) func(args []string, options map[string]string) int {

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
				fmt.Println(string(payload))
			} else {
				fmt.Println(err)
			}
		}
	}

	return func(args []string, options map[string]string) int {

		prettyPrintJson, _ := strconv.ParseBool(options["pretty"])
		watch, _ := strconv.ParseBool(options["watch"])
		jsonfy := util.CreateJsonfy(cfg)

		srcDir := options["src"]
		if !strings.HasSuffix(srcDir, string(os.PathSeparator)) {
			srcDir += string(os.PathSeparator)
		}

		files, err := ioutil.ReadDir(srcDir)

		if err != nil {
			log.Fatal(err)
		}

		processFile := func(filePath string) {
			jsons := jsonfy(filePath)
			processPayloads(jsons, prettyPrintJson)
		}

		for _, file := range files {
			log.Println(file.Name(), file.ModTime())

			processFile(srcDir + file.Name())
		}

		if watch {
			util.WatchDirectory(srcDir, processFile)
		}

		return 0
	}
}
