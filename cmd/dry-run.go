package cmd

import (
	"fmt"
	"github.com/ranskills/mp/setting"
	"github.com/ranskills/mp/util"
	"io/ioutil"
	"log"
	"strconv"
)

// CreateDryRunAction Creates a dry run action
func CreateDryRunAction(cfg setting.Config) func(args []string, options map[string]string) int {

	return func(args []string, options map[string]string) int {

		prettyPrintJson, _ := strconv.ParseBool(options["pretty"])

		processFile := util.CreateProcessFileFunc(cfg, prettyPrintJson)

		srcDir := options["src"]

		files, err := ioutil.ReadDir(srcDir)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name(), file.ModTime())

			// util.ProcessFile("./files/" + file.Name())
			processFile(srcDir + file.Name())
		}

		// watchDirectory()

		return 0
	}
}

//func dryRunAction(args []string, options map[string]string) int {
//
//	// processFile := util.CreateProcessFileFunc(cfg)
//	fmt.Println(args)
//	fmt.Println("config", options["config"])
//	fmt.Println("src", options["src"])
//
//	// return 0
//
//	files, err := ioutil.ReadDir(options["src"])
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, file := range files {
//		fmt.Println(file.Name(), file.ModTime())
//
//		// util.ProcessFile("./files/" + file.Name())
//		// processFile("./files/" + file.Name())
//	}
//
//	// watchDirectory()
//
//	return 0
//}
