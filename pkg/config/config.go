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
	// checking if the file is not empty
	if len(yml) == 0 {
		fmt.Println("The YAML file is empty")
		os.Exit(0)
	}
	// writing YML file contents in the struct
	err = yaml.Unmarshal(yml, &YAMLConfig)
	if err != nil {
		log.Fatal(err)
	}
}