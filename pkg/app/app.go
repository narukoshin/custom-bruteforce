package app

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/bruteforce"
	"custom-bruteforce/pkg/site"
	"custom-bruteforce/pkg/email"
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
	// testing email server
	{
		conn, client, err := email.Test_Connection()
		if err != nil {
			fmt.Printf("error: email: %v\n", err)
			return
		}
		// closing the connection after the test is done
		if conn != nil && client != nil {
			client.Close()
			conn.Close()
		}
	}
	// starting a bruteforce attack
	err := bruteforce.Start()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}