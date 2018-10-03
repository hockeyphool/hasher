package main

import (
	"bytes"
	"crypto/sha512"
	"testing"
)

func TestEnc(t *testing.T) {
	testPw := "angryMonkey"
	var expEncPw = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	testHash := hash(testPw)
	testEncPw := encode(testHash)

	if testEncPw != expEncPw {
		t.Errorf("Calculated Base64 encoding did not match expected;\nExpected: '%s'\nReceived: '%s'\n", expEncPw, testEncPw)
	}
}

func TestHash(t *testing.T) {
	testPw := "angryMonkey"
	expSha512 := sha512.New()
	expSha512.Write([]byte(testPw))

	testSha512 := hash(testPw)

	if bytes.Compare(expSha512.Sum(nil), testSha512.Sum(nil)) != 0 {
		t.Errorf("Calculated SHA512 hash did not match expected")
	}
}
