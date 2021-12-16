package bruteforce

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/structs"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var Types_Available []string = []string{"list", "file"}

var (
	Field 	string 		= config.YAMLConfig.B.Field
	Type  	string 		= config.YAMLConfig.B.Type
	Source	string 		= config.YAMLConfig.B.Source
	List	[]string	= config.YAMLConfig.B.List
	Fail	structs.YAMLOn_fail = config.YAMLConfig.OF
)

func Dictionary() []string {
	if Type == "list" {
		if len(List) == 0 {
			fmt.Println("Please add passwords to the list")
			os.Exit(0)
		}
		return List
	} else if Type == "file" {
		if _, err := os.Stat(Source); os.IsNotExist(err) {
			fmt.Printf("File %s does not exist\n", Source)
			os.Exit(0)
		}
		contents, err := ioutil.ReadFile(Source)
		if err != nil {
			log.Fatal(err)
		}
		return strings.Split(string(contents), "\n")
	}
	return []string{}
}