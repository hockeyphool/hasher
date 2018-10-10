package main

import (
	"encoding/json"
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

func TestInitStats(t *testing.T) {
	stats.NumEncodings = 2
	stats.AverageDuration = 5000000

	initStats()

	if stats.NumEncodings != 0 {
		t.Errorf("Expected number of requests %d, actual %d", 0, stats.NumEncodings)
	}

	if stats.AverageDuration != 0 {
		t.Errorf("Expected average duration %d, actual %d", 0, stats.AverageDuration)
	}
}

func TestUpdateStats(t *testing.T) {
	initStats()
	test := 0

	reqDurationInt1 := 5001
	t1 := time.Duration(reqDurationInt1 * nanosecondsPerMillisecond)
	expDuration1 := int64(t1 / time.Microsecond)
	updateStats(t1)

	test++
	if stats.NumEncodings != test {
		t.Errorf("Test %d: Expected %d requests, got %d instead", test, test, stats.NumEncodings)
	}

	if stats.AverageDuration != expDuration1 {
		t.Errorf("Test %d: Expected avg duration %d, got %d instead", test, expDuration1, stats.AverageDuration)
	}

	reqDurationInt2 := 5999
	t2 := time.Duration(reqDurationInt2 * nanosecondsPerMillisecond)
	expDuration2 := int64(t2 / time.Microsecond)
	updateStats(t2)

	expDuration2 = (expDuration1 + expDuration2) / 2
	test++
	if stats.NumEncodings != test {
		t.Errorf("Test %d: Expected %d requests, got %d instead", test, test, stats.NumEncodings)
	}

	if stats.AverageDuration != expDuration2 {
		t.Errorf("Test %d: Expected avg duration %d, got %d instead", test, expDuration2, stats.AverageDuration)
	}
}

func TestMarshalStats(t *testing.T) {
	const expNumEncodings = 5
	const expAvgDuration = 5987654

	initStats()
	stats.NumEncodings = expNumEncodings
	stats.AverageDuration = expAvgDuration

	testStats := marshalStats()
	var testStatsBuf serverStats
	json.Unmarshal(testStats, &testStatsBuf)

	if testStatsBuf.NumEncodings != expNumEncodings {
		t.Errorf("Expected %d requests, got %d", expNumEncodings, testStatsBuf.NumEncodings)
	}

	if testStatsBuf.AverageDuration != expAvgDuration {
		t.Errorf("Expected %d average duration, got %d", expAvgDuration, testStatsBuf.AverageDuration)
	}
}

func TestPasswordHandlerGoodStatus(t *testing.T) {
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

func TestPasswordHandlerBadStatus(t *testing.T) {
	var passwordKey = "password"
	var testPassword = "itShouldNotMatter"

	req, err := http.NewRequest("POST", "/hsh", nil)

	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add(passwordKey, testPassword)
	req.PostForm = form

	processTransactions = false
	var wg sync.WaitGroup
	testRecorder := httptest.NewRecorder()
	passwordHandler := buildPasswordHandler(&wg)

	passwordHandler.ServeHTTP(testRecorder, req)

	if status := testRecorder.Code; status != http.StatusForbidden {
		t.Errorf("Handler returned incorrect status - rcvd: %v, exp: %v\n", status, http.StatusOK)
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

func TestStatsHandlerGoodStatus(t *testing.T) {
	const expNumEncodings = 2
	const expAvgDuration = 5572500

	initStats()

	req, err := http.NewRequest("GET", "/stats", nil)

	if err != nil {
		t.Fatal(err)
	}

	reqDurationInt1 := 5234
	reqDurationInt2 := 5911

	t1 := time.Duration(reqDurationInt1 * nanosecondsPerMillisecond)
	t2 := time.Duration(reqDurationInt2 * nanosecondsPerMillisecond)

	updateStats(t1)
	updateStats(t2)

	processTransactions = true

	testRecorder := httptest.NewRecorder()
	statsHandler := buildStatsHandler()

	statsHandler.ServeHTTP(testRecorder, req)

	if status := testRecorder.Code; status != http.StatusOK {
		t.Errorf("Test Good Status - Handler returned incorrect status - rcvd: %v, exp: %v\n", status, http.StatusOK)
	}

	testBody := testRecorder.Body
	var testStats serverStats
	json.Unmarshal(testBody.Bytes(), &testStats)

	if testStats.NumEncodings != expNumEncodings {
		t.Errorf("serverStats: expected %d requests, got %d instead", expNumEncodings, testStats.NumEncodings)
	}

	if testStats.AverageDuration != expAvgDuration {
		t.Errorf("serverStats: expected %d average duration, got %d instead", expAvgDuration, testStats.AverageDuration)
	}
}

func TestStatsHandlerBadStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats", nil)

	if err != nil {
		t.Fatal(err)
	}

	processTransactions = false

	testRecorder := httptest.NewRecorder()
	statsHandler := buildStatsHandler()

	statsHandler.ServeHTTP(testRecorder, req)

	if status := testRecorder.Code; status != http.StatusForbidden {
		t.Errorf("Test Bad Status - Handler returned incorrect status - rcvd: %v, exp: %v\n", status, http.StatusOK)
	}
}
