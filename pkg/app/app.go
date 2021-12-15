package app

import (
	"custom-bruteforce/pkg/headers"
	// "custom-bruteforce/pkg/bruteforce"
	"custom-bruteforce/pkg/site"
	"fmt"
	"log"
	"net/http"
)

func Run(){
	req, err := http.NewRequest("POST", "http://localhost", nil)
	if err != nil {
		log.Fatal(err)
	}
	if headers.Is(){
		for _, header := range headers.Get() {
			req.Header.Set(header.Name, header.Value)
		}
	}
	// fmt.Println(req.Header.Get("User-Agent"))
	// fmt.Println(bruteforce.Dictionary())

	fmt.Println(site.Host)
	fmt.Println(site.Method)
	fmt.Println(site.Fields[0].Value)
}