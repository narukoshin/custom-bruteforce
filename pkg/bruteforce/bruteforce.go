package bruteforce

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/structs"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

var Types_Available []string = []string{"list", "file"}

var (
	Field 	string 		= config.YAMLConfig.B.Field
	From  	string 		= config.YAMLConfig.B.From
	File	string 		= config.YAMLConfig.B.File
	List	[]string	= config.YAMLConfig.B.List
	Fail	structs.YAMLOn_fail = config.YAMLConfig.OF
	Threads int			= config.YAMLConfig.B.Threads
)

type Attack_Result struct {
	Status		bool
	Password	string
}
var Attack Attack_Result

// adding some error messages
var ErrNoPasswords 		= errors.New("there is no passwords available for bruteforce, please specify some passwords")
var ErrOpeningFIle 		= errors.New("we have issues with opening a file, make sure that file exists and is readable")
var ErrWrongType   		= errors.New("you specified the wrong source of dictionary, allowed types are (file, list)")
var ErrEmptyField  		= errors.New("the field that you want to bruteforce is empty")
var ErrTooMuchThreads	= errors.New("too much threads for such small wordlist, please decrease amount of threads") 

// verifying if the list type is correct, currently there is only two types available - file and list
func verify_type() bool{
	for _, t := range Types_Available {
		if From == t {
			return true
		}
	}
	return false
}

// Loading dictionary from file or from the list that is defined in the config file
func Dictionary() ([][]string, error) {
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
	
	// spliting the wordlist for eaxh thread
	// To split a wordlist we need to know how much threads we want to use, if threads are not set, we will set them to default one
	if Threads == 0 {
		Threads = 5
	}
	// calculating the length how much passwords will be in one thread
	var size int = len(wordlist) / Threads
	// creating the output slice
	var result = make([][]string, 0)
	// checking if the size is not less than 1, because too much threads can cause that
	if size < 1 {
		return nil, ErrTooMuchThreads
	}
	// the brain of spliting
	for i := 0; i < len(wordlist); i += size{
		end := i + size
		// for the last element
		if end > len(wordlist) {
			end = len(wordlist)
		}
		slice := wordlist[i:end]
		// this code moves passwords from one thread to another if there are too less passwords
		if len(slice) != size {
			// we need to get the last element
			last_slice := result[len(result)-1]
			// calculating the location of the last element
			index := len(result) - 1
			// deleting the last element
			result = append(result[:index], result[index+1:]...)
			// adding passwords to previous thread and adding back to result
			result = append(result, append(last_slice, slice...))
		} else {
			result = append(result, slice)
		}
	}
	// again checking if threads are not so much
	// length of the result slice need to match the count of threads
	if len(result) != Threads {
		return nil, ErrTooMuchThreads
	}
	return result, nil
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
	var wg sync.WaitGroup
	// adding +1 job
	wg.Add(len(wordlist))
	// starting reading the wordlist
	for _, w := range wordlist {
		// adding goroutine to run each thread in sync
		go func(w []string) {
			// finishing the job
			defer wg.Done()
			// finally reading the passwords
			for _, pass := range w {
				// launching attacks for next steps
				_run_attack(pass)
				if Attack.Status {
					// deleting the wordlist to stop the threads
					wordlist = [][]string{}
					return
				}
			}
		}(w)
	}
	wg.Wait()
	_attack_finished()
	return nil
}

// launching the thread brute-force attack
func _run_attack(pass string) {
	// brute-forcing password until the status is set to true
	if !Attack.Status{
		fmt.Printf("\033[34m[~] trying password: %v\033[0m\n", pass)
		if pass == "nogales1" {
			// setting attack status to true and setting the password that we brute-forced
			Attack = Attack_Result {Status: true, Password: pass}
		}
	}
}

func _attack_finished(){
	if Attack.Status && Attack.Password != "" {
		fmt.Printf("\033[32m[~] the thing that you are looking for is found: %v\033[0m\n", Attack.Password)
	}
}