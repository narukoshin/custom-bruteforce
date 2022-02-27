package bruteforce

import (
	"custom-bruteforce/pkg/config"
	"custom-bruteforce/pkg/headers"
	"custom-bruteforce/pkg/site"
	"custom-bruteforce/pkg/structs"
	"custom-bruteforce/pkg/proxy"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var Types_Available []string = []string{"list", "file", "stdin"}

var (
	Field 	string 		= config.YAMLConfig.B.Field
	From  	string 		= config.YAMLConfig.B.From
	File	string 		= config.YAMLConfig.B.File
	List	[]string	= config.YAMLConfig.B.List
	Fail	structs.YAMLOn_fail = config.YAMLConfig.OF
	Pass	structs.YAMLOn_pass = config.YAMLConfig.OP
	Threads int			= config.YAMLConfig.B.Threads
	NoVerbose bool		= config.YAMLConfig.B.NoVerbose
	Output	string		= config.YAMLConfig.B.Output

	// Crawl
	Crawl_Search string = config.YAMLConfig.C.Search
	Crawl_Url	 string = config.YAMLConfig.C.Url
	Crawl_Name	 string = config.YAMLConfig.C.Name
)

// some status messages
var (
	StatusFinished string = "finished"
	StatusFound	string = "found"
)

type Attack_Result struct {
	Status		string
	Password	string
	Stop 		bool
	ErrorMessage string
}
var Attack Attack_Result

// adding some error messages
var ErrNoPasswords 		= errors.New("there is no passwords available for bruteforce, please specify some passwords")
var ErrOpeningFile 		= errors.New("we have issues with opening a file, make sure that file exists and is readable")
var ErrWrongType   		= errors.New("you specified the wrong source of dictionary, allowed types are (file, list)")
var ErrEmptyField  		= errors.New("the field that you want to bruteforce is empty")
var ErrTooMuchThreads	= errors.New("too much threads for such small wordlist, please decrease amount of threads") 
var ErrUnixRequired     = errors.New("you can not use this feature on Windows, you can use WSL instead")
var ErrMissingGroup		= errors.New("you forget to add group to the crawl/search option")
var ErrNoCrawlName		= errors.New("you forget to add the name of the field for token, without that option we can't set the token")
var ErrThreadsLessZero	= errors.New("threads can't be less than zero, zero threads default value is 5")

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
	switch(From){
		case "list":
			wordlist = List
		case "file":
			contents, err := ioutil.ReadFile(File)
			if err != nil {
				return nil, ErrOpeningFile
			}
			wordlist = strings.Split(string(contents), "\n")
		case "stdin":
			if runtime.GOOS == "windows" {
				return nil, ErrUnixRequired
			}
			contents, err := ioutil.ReadFile("/dev/stdin")
			if err != nil {
			  return nil, ErrOpeningFile
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
	// checking if the threads is not less than 0
	if Threads < 0 {
		return nil, ErrThreadsLessZero
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
				err := _run_attack(pass)
				if err != nil {
					fmt.Printf("error: %v\n", err)
				}
				// TODO: Handle error message from the attack function
				if Attack.Stop {
					// deleting the wordlist to stop the threads
					wordlist = [][]string{}
					return
				}
			}
		}(w)
	}
	wg.Wait()
	// When the script stopped working, this will be printed out
	_attack_finished()
	return nil
}

// launching the thread brute-force attack
func _run_attack(pass string) error {
	if !Attack.Stop {
		client := http.Client{}
		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}
		client.Jar = jar

		// Adding proxy, if there is any
		if proxy.IsProxy() {
			client.Transport = proxy.Drive()
		}
		values := url.Values{}

		// checking if the token pattern is added
		if Crawl_Search != "" {
			// finding the token
			token, err := Bypassing_Security_Token(&client)
			if err != nil {
				Attack = Attack_Result {Status: StatusFinished, Stop: true, ErrorMessage: err.Error()}
				return nil
			}
			// Checking if the Crawl_Name is added because it's a very important option
			if len(Crawl_Name) == 0 {
				Attack = Attack_Result {Status: StatusFinished, Stop: true, ErrorMessage: err.Error()}
				return nil
			}
			values.Set(Crawl_Name, token)
		} 
		for _, field := range site.Fields {
			values.Set(field.Name, field.Value)
		}
		values.Set(Field, pass)

		// adding a slash to the host if there is no slash
		if site.Host[:len(site.Host)-1] != "/" {
			site.Host = site.Host + "/"
		}
		req, err := http.NewRequest(site.Method, site.Host, strings.NewReader(values.Encode()))
		if err != nil {
			return err
		}
		defer req.Body.Close()
		if headers.Is() {
			for _, header := range headers.Get(){
				req.Header.Set(header.Name, header.Value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if (len(Fail.Message) != 0 && !strings.Contains(string(body), Fail.Message)) && (len(Pass.Message) == 0) || (len(Pass.Message) != 0 && strings.Contains(string(body), Pass.Message)) {
				Attack = Attack_Result {Status: StatusFound, Password: pass, Stop: true}
				return nil
			}
			_while_running(pass)
		} else {
			Attack = Attack_Result {Status: StatusFinished, Stop: true, ErrorMessage: "The server says 404"}
			return nil
		}
	}
	return nil
}

// message that will be printed while the script is running
func _while_running(pass string){
	if !NoVerbose {
		fmt.Printf("\033[34m[~] trying password: %v\033[0m\n", pass)
	}
}

// message that will be printed out when the script is finished
func _attack_finished(){
	// checking if the attack is stopped and the password is found
	if Attack.Stop && Attack.Status == StatusFound  && Attack.Password != "" {
		fmt.Printf("\033[32m[~] the thing that you were looking for is found: %v\033[0m\n", Attack.Password)
		// there we will save the password
		WritePasswordToFile()
		return
	}
	fmt.Printf("\033[33m[~] Well, looks that we can't find a thing that you need, sorry. :/\033[0m\n")
	if len(Attack.ErrorMessage) != 0 {
		fmt.Printf("\033[33m[~] Error: %s\033[0m\n", Attack.ErrorMessage)
		return
	}
}

// Saving the password in the file
func WritePasswordToFile(){
	// checking if the "output" option is added
	if len(Output) != 0 {
		// writting password to the file
		ioutil.WriteFile(Output, []byte(Attack.Password), 0644)
	}
}

// Crawling out the token
func Bypassing_Security_Token(client *http.Client) (string, error) {
	var Token_Uri string
	if len(Crawl_Url) != 0 {
		Token_Uri = Crawl_Url
	} else {
		Token_Uri = site.Host
	}

	req, err := http.NewRequest(http.MethodGet, Token_Uri, nil)
	if err != nil {
		return "", err
	}
	if headers.Is() {
		for _, header := range headers.Get(){
			req.Header.Set(header.Name, header.Value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(Crawl_Search)

		if len(re.FindSubmatch(body)) != 2 {
			return "", ErrMissingGroup
		}
		return string(re.FindSubmatch(body)[1]), nil
	}
	return "", nil
}