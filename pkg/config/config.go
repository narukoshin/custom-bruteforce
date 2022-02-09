package config

import (
	"custom-bruteforce/pkg/structs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"errors"
	"os"
)

const YAMLFile string = "config.yml"

var YAMLConfig structs.YAMLConfig

// Handling errors
var CError error = nil

// Error message if the config file is not found
var ErrConfigNotFound = errors.New("config file not found, please create it to use this tool")

// Error message if the config file is empty
var ErrConfigIsEmpty = errors.New("config file is empty")

func init() {
	// Checking if the config file exists
	if _, err := os.Stat(YAMLFile); err != nil {
		CError = ErrConfigNotFound
		return
	}
	// Reading config file
	yml, err := ioutil.ReadFile(YAMLFile)
	if err != nil {
		CError = err
		return
	}
	// Checking if the config file is not empty
	if len(yml) == 0 {
		CError = ErrConfigIsEmpty
		return
	}
	// writing YML file contents in the struct
	err = yaml.Unmarshal(yml, &YAMLConfig)
	if err != nil {
		CError = err
		return
	}
}

// Checking for errors in the app.Run function
func HasError() error {
	return CError
}