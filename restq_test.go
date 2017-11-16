package main

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

//TestGetProj verifies that json file can be read and projnr returned
func TestGetProj(t *testing.T) {
	var projnr string
	//Test a correct file
	f, err := ioutil.TempFile("", "testqueueid")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte("{\"Queue\":\"a-long-queue-name-123\"}"), 0644)
	queue, err := GetProj(projnr, f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if queue != "a-long-queue-name-123" {
		t.Fatal(queue + " is not equal to a-long-queue-name-123")
	}
	//Verify none existing file breaks
	_, err = GetProj(projnr, "none")
	if err == nil {
		t.Fatal(err)
	}

	//verify broken json breaks
	f, err = ioutil.TempFile("", "testqueueidbroken")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte("\"Queue\":\"a-long-queue-name-123\"}"), 0644)
	queue, err = GetProj(projnr, f.Name())
	if err == nil {
		t.Fatal("Broken json in config should return error")
	}

	//Verify projnr works
	proj, err := GetProj("123", "none")
	if err != nil {
		t.Fatal(err)
	}
	if proj != "123" {
		t.Fatal(err)
	}

}

//TestCreate verifies that a project can be created
func TestCreate(t *testing.T) {
	url := "http://restq.io/rest_api/"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", url,
		httpmock.NewStringResponder(200, `{"Queue": "1"}`))
	queuejson, _, err := Create()
	res := Projfile{}
	json.Unmarshal([]byte(queuejson), &res)
	if err != nil || res.Queue != "1" {
		t.Fatal("Could not create Project!")
	}
	httpmock.RegisterResponder("POST", url,
		httpmock.NewStringResponder(404, ""))
	_, _, err = Create()
	if err == nil {
		t.Fatal("No error on status 404 on creating project!")
	}
}

//TestGet verifies that Get returns a json object and empty queue
func TestGet(t *testing.T) {
	proj := "a-long-queue-name-123"
	url := "http://restq.io/rest_api/" + proj
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//read response from item from queue
	complete_response := `{"Queue": "1"}`
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, complete_response))
	itemjson, status, err := Get("a-long-queue-name-123")
	if err != nil {
		t.Fatal("Could not get item from project!")
	}
	if status != 200 {
		t.Fatal("Wrong status!")
	}
	res := Projfile{}
	json.Unmarshal([]byte(itemjson), &res)
	if res.Queue != "1" {
		t.Log(res)
		t.Fatal("Could not read item!")
	}

	//receive status that queue is empty
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(204, ""))
	itemjson, status, err = Get("a-long-queue-name-123")
	if status != 204 {
		t.Fatal("Wrong status should be 204!")
	}

	//receive status that queue is broken
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(500, ""))
	itemjson, status, err = Get("a-long-queue-name-123")
	if err == nil {
		t.Fatal("Could get item from broken server")
	}

}

//TestGet verifies that Put returns err if response is not 200
func TestPut(t *testing.T) {
	proj := "a-long-queue-name-123"
	url := "http://restq.io/rest_api/" + proj
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//put one message on queue
	httpmock.RegisterResponder("PUT", url,
		httpmock.NewStringResponder(200, ""))
	err := Put(`{"Queue": "1"}`, proj)
	if err != nil {
		t.Fatal("Could not Put")
	}

	//put one which should fail
	httpmock.RegisterResponder("PUT", url,
		httpmock.NewStringResponder(404, ""))
	err = Put(`{"Queue": "1"}`, proj)
	if err == nil {
		t.Fatal("Could Put!")
	}

	//using broken http
	httpmock.RegisterResponder("PUT", url,
		httpmock.NewStringResponder(404, ""))
	err = Put("1", " ?")
	if err == nil {
		t.Fatal("Could !")
	}
}

func TestDispatch(t *testing.T) {
	filePtr = ""
	projPtr = ""
	createPtr = false
	putPtr = false
	getPtr = false
	testargs := []string{}

	//all empty flags should fail
	_, _, err := Dispatch(testargs)
	if err == nil {
		t.Fatal(err)
	}

	filePtr = ""
	projPtr = ""
	createPtr = true
	putPtr = false
	getPtr = false
	//create flag only
	_, _, err = Dispatch(testargs)
	if err != nil {
		t.Fatal(err)
	}

	filePtr = ""
	projPtr = ""
	createPtr = false
	putPtr = true
	getPtr = false
	//put flag only
	_, _, err = Dispatch(testargs)
	if err == nil {
		t.Fatal(err)
	}

	//start test which requires http response
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	//put
	proj := "a-long-queue-name-123"
	url := "http://restq.io/rest_api/" + proj
	filePtr = ""
	projPtr = "a-long-queue-name-123"
	msgPtr = "message"
	createPtr = false
	putPtr = true
	getPtr = false
	quietPtr = false
	httpmock.RegisterResponder("PUT", url,
		httpmock.NewStringResponder(200, ""))
	output, status, err := Dispatch(testargs)
	if err != nil {
		t.Log(output)
		t.Log(status)
		t.Fatal(err)
	}
	//put with message should not fail
	proj = "a-long-queue-name-123"
	url = "http://restq.io/rest_api/" + proj
	filePtr = ""
	projPtr = "a-long-queue-name-123"
	msgPtr = "testmessage"
	createPtr = false
	putPtr = true
	getPtr = false
	quietPtr = false
	httpmock.RegisterResponder("PUT", url,
		httpmock.NewStringResponder(200, ""))
	output, status, err = Dispatch(testargs)
	if err != nil {
		t.Fatal(err)
	}
	//put with no message should fail but will actually accept empty input now
	proj = "a-long-queue-name-123"
	url = "http://restq.io/rest_api/" + proj
	filePtr = ""
	projPtr = "a-long-queue-name-123"
	msgPtr = ""
	createPtr = false
	putPtr = true
	getPtr = false
	quietPtr = false
	httpmock.RegisterResponder("PUT", url,
		httpmock.NewStringResponder(200, ""))
	output, status, err = Dispatch(testargs)
	if err != nil {
		t.Fatal(err)
	}
	//get fail from correct post should fail
	proj = "a-long-queue-name-123"
	url = "http://restq.io/rest_api/" + proj
	filePtr = ""
	projPtr = "a-long-queue-name-123"
	msgPtr = ""
	createPtr = false
	putPtr = true
	getPtr = false
	quietPtr = false
	httpmock.RegisterResponder("PUT", url,
		httpmock.NewStringResponder(404, ""))
	_, _, err = Dispatch(testargs)
	if err == nil {
		t.Fatal(err)
	}
	//getting message ok
	filePtr = ""
	projPtr = "a-long-queue-name-123"
	msgPtr = ""
	createPtr = false
	putPtr = false
	getPtr = true
	quietPtr = false
	//read response from item from queue
	complete_response := `{"Queue": "1"}`
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, complete_response))
	_, status, err = Dispatch(testargs)
	if err != nil {
		t.Fatal("Could not get item from project!")
	}
	if status != 200 {
		t.Fatal("Wrong status!")
	}
	//getting queue empty
	filePtr = ""
	projPtr = "a-long-queue-name-123"
	msgPtr = ""
	createPtr = false
	putPtr = false
	getPtr = true
	quietPtr = false
	complete_response = `{"Queue": "1"}`
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(204, complete_response))
	_, status, err = Dispatch(testargs)
	if err != nil {
		t.Fatal("Could not get item from project!")
	}
	if status != 204 {
		t.Fatal("Wrong status!")
	}
	//getting queue empty
	filePtr = ""
	projPtr = "a-long-queue-name-123"
	msgPtr = ""
	createPtr = false
	putPtr = false
	getPtr = true
	quietPtr = false
	complete_response = `{"Queue": "1"}`
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(500, complete_response))
	_, status, err = Dispatch(testargs)
	if err == nil {
		t.Fatal("Should be error on server error")
	}
	if status == 0 || status == 200 || status == 204 {
		t.Fatal("Server error should not be ok")
	}
}

//Test stdin
func TestFindstdin(t *testing.T) {
	//test no stdin
	message := FindStdin()
	if message != "" {
		t.Fatal("Should not get anything on stdin")
	}
}

//TestMain run presteps before m.Run and teardown after
func TestMain(m *testing.M) {
	output := m.Run()
	os.Exit(output)
}
