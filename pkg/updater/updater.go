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

type HasUpdatesToInstall struct {
	LatestVersion string
	ExecutableName string
}


func CheckForUpdate(currentVersion string) (HasUpdatesToInstall, error) {
	resp, err := http.Get(api_endpoint)
	if err != nil {
		return HasUpdatesToInstall{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return HasUpdatesToInstall{}, err
		}
		releases := []Release{}
		err = json.Unmarshal(body, &releases)
		if err != nil {
			return HasUpdatesToInstall{}, err
		}
		// checking if there is any release to avoid errors in future
		if len(releases) > 0 {
			// latest release
			Latest = releases[0]
			if currentVersion >= Latest.Version {
				return HasUpdatesToInstall{}, nil
			} else {
				path, err := os.Executable()
				if err != nil {
					return HasUpdatesToInstall{}, err
				}
				updates := HasUpdatesToInstall {
					LatestVersion: Latest.Version,
					ExecutableName: filepath.Base(path),
				}
				return updates, nil
			}
		}
	}
	return HasUpdatesToInstall{}, nil
}

func InstallUpdate(currentVersion string) error {
	// checking if there is an update to install
	if updates, err := CheckForUpdate(currentVersion); err == nil {
		if (HasUpdatesToInstall{}) != updates {
			fmt.Printf("\033[36m[-] Starting update...\n\033[0m")
			path, err := os.Executable()
			if err != nil {
				return err
			}
			// Reading the binary from the github repository
			var gitFileName string
			// adding .exe to the end if it's a Windows
			if system == "windows" {
				gitFileName = system + ".exe"
			} else {
				gitFileName = system
			}
			resp, err := http.Get(fmt.Sprintf(binariesf, gitFileName))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			// deleting an old binary file
			err = os.Remove(path)
			if err != nil {
				// if we can't delete an existing file, then we will rename it and download a newer version.
				err = os.Rename(path, fmt.Sprintf("%s.%s.bak", path, currentVersion))
				if err != nil {
				return err
				}
			}
			// creating a new binary file
			fp, err := os.Create(path)
			if err != nil {
				return err
			}
			defer fp.Close()
			// copying content from the github repository to the file
			size, err := io.Copy(fp, resp.Body)
			if err != nil {
				return err
			}
			if size > 0 {
				fmt.Println("\033[32m[-] Latest version has been downloaded and replaced,\n please verify version with a version check command.\033[0m")
			}
			// making it executable
			err = fp.Chmod(0544)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("You already has the latest version")
		}
	}
	return nil
}