package cmd

import "fmt"

// VERSION The current version of this tool
const VERSION = "1.0.0-dev"

// VersionHandler Prints the current version of the CLI
func VersionHandler(args []string, options map[string]string) int {
	fmt.Printf("c2m %s\n", VERSION)
	return 0
}
