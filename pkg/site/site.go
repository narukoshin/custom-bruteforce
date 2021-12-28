package site

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/structs"
	"errors"
	"net"
	"net/url"
	"strings"
)

var (
	Host 	string	= config.YAMLConfig.S.Host
	Method 	string  = config.YAMLConfig.S.Method
	Fields  []structs.YAMLFields = config.YAMLConfig.F
)

// Error message if the request method in the config is incorrect
var ErrInvalidMethod = errors.New("please specify a valid request method")
var ErrDeadHost		 = errors.New("looks that the host is not alive, check your config again")

// All request methods that are allowed to use
var Methods_Allowed []string = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}

// Verifying if the request method is correct
func Verify_Method() error {
	for _, value := range Methods_Allowed {
		if ok := strings.EqualFold(Method, value); ok {
			return nil
		}
	}
	return ErrInvalidMethod
}

// Verifying if the host of target is specified correctly
func Verify_Host() error {
	if _, err := url.ParseRequestURI(Host); err != nil {
		return err
	}
	if _, err := net.Dial("tcp", Host); err != nil {
		return ErrDeadHost
	}
	return nil
}