package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Projfile struct {
	Queue string `json:"queue"`
}

//GetProj returns projnr from inparam or from file from inparam
func GetProj(projnr, file string) (string, error) {
	var proj string
	var err error
	if projnr != "" {
		proj = projnr
	} else {
		res := Projfile{}
		raw, err := ioutil.ReadFile(file)
		if err != nil {
			return proj, err
		}
		err = json.Unmarshal(raw, &res)
		if err != nil {
			return proj, err
		} else {
			proj = res.Queue
		}
	}
	return proj, err
}

//Get one item from project
func Get(proj string) (string, int, error) {
	url := "http://restq.io/rest_api/" + proj
	resp, err := http.Get(url)
	if resp.StatusCode == 200 {
		responseData, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		return string(responseData), resp.StatusCode, err
	}
	if resp.StatusCode == 204 {
		return "Empty queue", resp.StatusCode, err
	}
	err = errors.New("Error getting item")
	return "", resp.StatusCode, err
}

//Put item on proj queue
func Put(message, proj string) error {
	url := "http://restq.io/rest_api/" + proj
	b := bytes.NewBufferString(message)
	req, _ := http.NewRequest("PUT", url, b)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err = errors.New("Error posting new item")
	}
	defer resp.Body.Close()
	return err
}

//Create a new project return string from server
func Create() (string, int, error) {
	var queue string
	status := 0
	resp, err := http.Post("http://restq.io/rest_api/", "", nil)
	if resp.StatusCode != 200 {
		err = errors.New("Error status on creating new project")
		return queue, resp.StatusCode, err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	queue = string(responseData)
	defer resp.Body.Close()
	return queue, status, err
}

//FindStdin read stdin to find message
func FindStdin() string {
	var message string
	messin, _ := os.Stdin.Stat()
	if messin.Mode()&os.ModeNamedPipe != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message = message + scanner.Text()
		}
	}
	return message
}

var filePtr string
var projPtr string
var msgPtr string
var createPtr bool
var putPtr bool
var getPtr bool
var quietPtr bool

func init() {
	flag.StringVar(&filePtr, "f", "", "Project file")
	flag.StringVar(&projPtr, "i", "", "Project id")
	flag.StringVar(&msgPtr, "m", "", "Message if not on stdin")
	flag.BoolVar(&createPtr, "c", false, "Create new project")
	flag.BoolVar(&putPtr, "p", false, "Post")
	flag.BoolVar(&getPtr, "g", false, "Get message")
	flag.BoolVar(&quietPtr, "q", false, "Silent empty queue message")
}

//Dispatch to different actions
func Dispatch(args []string) (string, int, error) {
	var message, output string
	var status int
	var err error

	flag.Parse()
	if (createPtr == false &&
		putPtr == false &&
		getPtr == false) ||
		(filePtr == "" &&
			projPtr == "" &&
			(putPtr == true ||
				getPtr == true)) {
		err = errors.New("Invalid input!")
		return output, status, err
	}

	switch {
	case createPtr == true:
		output, status, err := Create()
		return output, status, err
	case putPtr == true:
		proj, _ := GetProj(projPtr, filePtr)
		if msgPtr == "" {
			message = FindStdin()
		} else {
			message = msgPtr
		}
		err = Put(message, proj)
		if err != nil {
			return output, 1, err
		}
	case getPtr == true:
		proj, _ := GetProj(projPtr, filePtr)
		output, status, err := Get(proj)
		if err != nil {
			return output, 1, err
		}
		if status == 204 && quietPtr == false {
			return output, status, err
		} else {
			return output, status, err
		}
	}
	return output, status, err
}

func main() {
	output, status, _ := Dispatch(flag.Args())
	fmt.Print(output)
	os.Exit(status)
}
