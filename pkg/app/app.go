package app

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/bruteforce"
	"custom-bruteforce/pkg/site"
	"fmt"
)

func Run(){
	// verifying the config
	if err := config.HasError(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	// verifying the target host uri
	if err := site.Verify_Host(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	// verifying the request method
	if err := site.Verify_Method(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	// starting a bruteforce attack
	err := bruteforce.Start()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}