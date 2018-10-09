package main

import (
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
