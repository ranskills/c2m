// Package util provides ...
package util

import (
	"fmt"
	"strings"
)

// GetConfigOptionValue Returns the config option value or defaults to config.yml in the current directory
func GetOptionValue(args []string, shortOption, longOption string) (string, error) {

	for k, v := range args {

		usesLongOption := strings.HasPrefix(v, longOption)
		usesShortOption := v == shortOption

		if usesLongOption {
			parts := strings.Split(v, "=")

			if len(parts) == 2 {
				return parts[1], nil
			} else {
				return "", fmt.Errorf("GetOptionValue: No value provided for the option %s", longOption)
			}
		} else if usesShortOption {
			if len(args) > k+1 {
				return args[k+1], nil
			} else {
				return "", fmt.Errorf("GetOptionValue: No value provided for the option %s", shortOption)
			}
		}
	}

	return "", fmt.Errorf("GetOptionValue: None of the options %s|%s are present", shortOption, longOption)
}

// GetConfigOptionValue Returns the config option value or defaults to config.yml in the current directory
func GetConfigOptionValue(args []string) string {
	v, err := GetOptionValue(args, "-c", "--config")

	if err == nil {
		return v
	}

	return "./config.yml"
}
