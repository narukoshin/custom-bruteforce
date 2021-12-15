package site

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/structs"
)

var (
	Host 	string	= config.YAMLConfig.S.Host
	Method 	string  = config.YAMLConfig.S.Method
	Fields  []structs.YAMLFields = config.YAMLConfig.F
)