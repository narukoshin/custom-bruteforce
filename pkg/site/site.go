package site

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/structs"
	"errors"
	"fmt"
	"os"
)

var (
	Host 	string	= config.YAMLConfig.S.Host
	Method 	string  = config.YAMLConfig.S.Method
	Fields  []structs.YAMLFields = config.YAMLConfig.F
)

var ErrInvalidMethod = errors.New("please specify a valid request method")

var Methods_Allowed []string = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}


func init(){
	if ok := verify_method(); !ok {
		fmt.Printf("%v\n", ErrInvalidMethod)
		os.Exit(0)
	}
}

func verify_method() bool{
	for _, value := range Methods_Allowed {
		if value == Method {
			return true
		}
	}
	return false
}