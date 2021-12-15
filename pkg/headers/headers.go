package headers

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/structs"
)

func Is() bool {
	return len(config.YAMLConfig.H) != 0
}

func Get() []structs.YAMLHeaders {
	return config.YAMLConfig.H
}