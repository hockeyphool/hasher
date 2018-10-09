package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

const testMinInt = 5000
const testMaxInt = 6000
const nanosecondsPerMillisecond = 1000000

func TestEncoding(t *testing.T) {
	testPw := "angryMonkey"
	var expEncPw = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	testEncPw := encode(testPw)

	if testEncPw != expEncPw {
		t.Errorf("Base64-encoded hash did not match expected;\nExpected: '%s'\nReceived: '%s'\n", expEncPw, testEncPw)
	}
}

func TestPortValidation(t *testing.T) {
	const invalidPortLow = "1023"
	const invalidPortHigh = "9001"
	const validPort = "1024"

	if isPortValid(invalidPortLow) {
		t.Errorf("Too-low port %s was not validated correctly", invalidPortLow)
	}

	if isPortValid(invalidPortHigh) {
		t.Errorf("Too-high port %s was not validated correctly", invalidPortHigh)

	}

	if !isPortValid(validPort) {
		t.Errorf("Valid port %s was not validated correctly", validPort)
	}
}

func TestRandomInt(t *testing.T) {
	testInt := getRandomInt()

	if testMinInt > testInt {
		t.Errorf("getRandomInt() returned %d which is less than minimum %d", testInt, testMinInt)
	}

	if testMaxInt < testInt {
		t.Errorf("getRandomInt() returned %d which is greater than maximum %d", testInt, testMaxInt)
	}
}

func TestSleepInterval(t *testing.T) {
	const testMinSleep = time.Duration(testMinInt * nanosecondsPerMillisecond)
	const testMaxSleep = time.Duration(testMaxInt * nanosecondsPerMillisecond)

	testDuration := getSleepInterval()

	if testMinSleep > testDuration {
		t.Errorf("getSleepInterval() returned duration %d which is less than minimum %d", testDuration, testMinSleep)
	}

	if testMaxSleep < testDuration {
		t.Errorf("getSleepInterval() returned duration %d which is greater than maximum %d", testDuration, testMaxSleep)
	}
}

func TestPasswordHandler(t *testing.T) {
	var expEncPw = "\"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==\""
	var passwordKey = "password"
	var testPassword = "angryMonkey"

	req, err := http.NewRequest("POST", "/hash", nil)

	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add(passwordKey, testPassword)
	req.PostForm = form

	var wg sync.WaitGroup
	testRecorder := httptest.NewRecorder()
	passwordHandler := buildPasswordHandler(&wg)

	passwordHandler.ServeHTTP(testRecorder, req)

	if status := testRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned incorrect status - rcvd: %v, exp: %v\n", status, http.StatusOK)
	}

	rcvdVal := testRecorder.Body.String()

	if rcvdVal != expEncPw {
		t.Errorf("Handler did not return expected value - \n  rcvd: '%v'\n  exp:  '%v'", rcvdVal, expEncPw)
	}
}

func TestShutdownHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/shutdown", nil)

	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	testQuitChan := make(chan bool)

	go func() {
		<-testQuitChan
		t.Log("Proceeding")
	}()

	testRecorder := httptest.NewRecorder()
	shutdownHandler := buildShutdownHandler(&wg, testQuitChan)

	shutdownHandler.ServeHTTP(testRecorder, req)

	if status := testRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned incorrect status - rcvd: %v, exp: %v\n", status, http.StatusOK)
	}

	if processTransactions {
		t.Errorf("Handler did not modify processTransactions correctly")
	}

}
