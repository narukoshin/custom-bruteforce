package app

import (
	"custom-bruteforce/pkg/bruteforce"
	"custom-bruteforce/pkg/updater"
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/email"
	"custom-bruteforce/pkg/site"
	"fmt"
	"os"
)

const Version string = "v2.4.5"

func Run(){
	// checking if there's any command used
	{
		if len(os.Args) != 1 {
			command := os.Args[1]
			switch command {
				case "version":
					fmt.Printf("Build version: %s\r\n", Version)
					err := updater.CheckForUpdate(Version)
					if err != nil {
					  fmt.Printf("error: updater: %v\n", err)
					  return
					}
				case "update":
					err := updater.InstallUpdate()
					if err != nil {
						fmt.Printf("error: updater: %v\n", err)
						return
					  }
			}
			return
		}
	}
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