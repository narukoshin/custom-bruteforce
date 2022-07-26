package config

import (
	"custom-bruteforce/pkg/structs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"errors"
	"os"
)

// the default config name that we will load.
const YAMLFile string = "config.yml"

var (
	YAMLConfig structs.YAMLConfig

	// Handling errors
	CError error = nil

	// Error message if the config file is not found
	ErrConfigNotFound = errors.New("config file not found, please create it to use this tool")

	// Error message if the config file is empty
	ErrConfigIsEmpty = errors.New("config file is empty")
)

func init(){
	yml := load_file(YAMLFile)
	err := yaml.Unmarshal(yml, &YAMLConfig)
	if err != nil {
		CError = err
		return
	}

	// if `import` option is not empty, then importing the file from the option
	// if the `import` option is empty, we will use config.yml that is loaded above.
	if len(YAMLConfig.Import) != 0 {
		yml = load_file(YAMLConfig.Import)
		err := yaml.Unmarshal(yml, &YAMLConfig)
		if err != nil {
			CError = err
			return
		}
	}
}

func load_file(file_name string) []byte {
	// Checking if the config file exists
	if _, err := os.Stat(file_name); err != nil {
		CError = ErrConfigNotFound
		return nil
	}
	// Reading config file
	yml, err := ioutil.ReadFile(file_name)
	if err != nil {
		CError = err
		return nil
	}
	// Checking if the config file is not empty
	if len(yml) == 0 {
		CError = ErrConfigIsEmpty
		return nil
	}
	return yml
}

// Checking for errors in the app.Run function
func HasError() error {
	return CError
}