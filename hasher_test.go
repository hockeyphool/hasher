package main

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestEnc(t *testing.T) {
	testPw := "angryMonkey"
	var expEncPw = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	testEncPw := encode(testPw)

	if testEncPw != expEncPw {
		t.Errorf("Calculated Base64 encoding did not match expected;\nExpected: '%s'\nReceived: '%s'\n", expEncPw, testEncPw)
	}
}

func TestInvalidPort(t *testing.T) {
	portTooLow := "1023"
	if isPortValid(portTooLow) {
		t.Errorf("Too-low port %sw as not validated correctly", portTooLow)
	}

	portTooHigh := "9001"
	if isPortValid(portTooHigh) {
		t.Errorf("Too-high port %s was not validated correctly", portTooHigh)

	}
}

func TestHttpEncode(t *testing.T) {
	//var expEncPw = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	reqUrl := "http://localhost:8080"
	reqResource := "/hash"
	reqData := url.Values{}
	reqData.Set("password", "angryMonkey")

	u, _ := url.ParseRequestURI(reqUrl)
	u.Path = reqResource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(reqData.Encode()))

	resp, _ := client.Do(r)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status: %v", resp.Status)
	}

	//if resp != expEncPw {
	//	t.Errorf("Received encoding did not match expected;\nExpected: '%s'\nReceived: '%s'\n", expEncPw, resp)
	//}
}
