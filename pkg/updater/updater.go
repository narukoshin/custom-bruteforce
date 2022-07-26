package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	api_endpoint string = "https://api.github.com/repos/narukoshin/custom-bruteforce/releases"
	binariesf	 string = "https://github.com/narukoshin/custom-bruteforce/blob/main/bin/%s?raw=true"
	system		 string = runtime.GOOS
)

var Latest Release

type Release struct {
	Version	   string `json:"name"`
	Prerelease bool   `json:"prerelease"`
}

func CheckForUpdate(currentVersion string) (error) {
	resp, err := http.Get(api_endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		releases := []Release{}
		err = json.Unmarshal(body, &releases)
		if err != nil {
			return err
		}
		// checking if there is any release to avoid errors in future
		if len(releases) > 0 {
			// latest release
			Latest = releases[0]
			if currentVersion == Latest.Version {
				fmt.Println("You already have the latest version")
			} else {
				path, err := os.Executable()
				if err != nil {
					return err
				}
				fmt.Printf("Latest available: %v\nUse %v update - to install the update\r\n", Latest.Version, filepath.Base(path))
			}
		}
	}
	return nil
}

func InstallUpdate() error {
	path, err := os.Executable()
	if err != nil {
		return err
	}
	// Reading the binary from the github repository
	resp, err := http.Get(fmt.Sprintf(binariesf, system))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// deleting an old binary file
	err = os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
	// creating a new binary file
	fp, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()
	// copying content from the github repository to the file
	size, err := io.Copy(fp, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if size > 0 {
		fmt.Println("Successfuly updated.")
	}
	return nil
}