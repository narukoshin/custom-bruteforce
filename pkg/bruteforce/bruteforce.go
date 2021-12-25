package bruteforce

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/structs"
	"errors"
	"io/ioutil"
	"strings"
)

var Types_Available []string = []string{"list", "file"}

var (
	Field 	string 		= config.YAMLConfig.B.Field
	From  	string 		= config.YAMLConfig.B.From
	File	string 		= config.YAMLConfig.B.File
	List	[]string	= config.YAMLConfig.B.List
	Fail	structs.YAMLOn_fail = config.YAMLConfig.OF
)

// adding some error messages
var ErrNoPasswords = errors.New("there is no passwords available for bruteforce, please specify some passwords")
var ErrOpeningFIle = errors.New("we have issues with opening a file, make sure that file exists and is readable")
var ErrWrongType   = errors.New("you specified the wrong source of dictionary, allowed types are (file, list)")
var ErrEmptyField  = errors.New("the field that you want to bruteforce is empty")

func verify_type() bool{
	for _, t := range Types_Available {
		if From == t {
			return true
		}
	}
	return false
}

// Loading dictionary from file or from the list that is defined in the config file
func Dictionary() ([]string, error) {
	var wordlist []string
	// Verifying that the source of the dictionary is allowed
	if ok := verify_type(); !ok {
		return nil, ErrWrongType
	}
	if From == "list" {
		wordlist = List
	} else if From == "file" {
		contents, err := ioutil.ReadFile(File)
		if err != nil {
			return nil, ErrOpeningFIle
		}
		wordlist = strings.Split(string(contents), "\n")
	}
	// Checking if the last element of the list is not empty
	// If the last element is empty, then we are deleting it 
	{
		if len(wordlist) != 0 {
			if len(wordlist[len(wordlist)-1]) == 0 {
				wordlist = wordlist[:len(wordlist)-1]
			}
		}
	}
	// When we cleared empty line at the end, let's check if the list is not empty now
	if len(wordlist) == 0 {
		return nil, ErrNoPasswords
	}	
	return wordlist, nil
}

// The function where all the magic happens
func Start() error {
	// loading wordlist
	wordlist, err := Dictionary()
	if err != nil {
		return err
	}
	// making sure that user specified the field that he wants to bruteforce
	if len(Field) == 0 {
		return ErrEmptyField
	}
	_ = wordlist
	// Threading coming soon....
	return nil
}