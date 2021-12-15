package config

import (
	"io/ioutil"
	"log"
	"fmt"
	"os"
	
	"gopkg.in/yaml.v3"
	"custom-bruteforce/pkg/structs"
)

const YAMLFile string = "config.yml"

var YAMLConfig structs.YAMLConfig

func init() {
	// checking if the YAML file is created
	if _, err := os.Stat(YAMLFile); err != nil {
		fmt.Printf("Please create %v file\n", YAMLFile)
		return
	}
	// reading the YAML file contents
	yml, err := ioutil.ReadFile(YAMLFile)
	if err != nil {
		log.Fatal(err)
	}
	// writing YML file contents in the struct
	yaml.Unmarshal(yml, &YAMLConfig)
}