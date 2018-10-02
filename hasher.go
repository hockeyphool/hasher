package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	hash2 "hash"
)

func hash(password string) hash2.Hash {
	mySha512 := sha512.New()
	mySha512.Write([]byte(password))
	return mySha512
}

func encode(hashData hash2.Hash) string {
	encHash := base64.StdEncoding.EncodeToString(hashData.Sum(nil))
	return encHash
}

func main() {
	fmt.Println("Hasher\n")
	var expEncPw = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	password := "angryMonkey"

	hashedPw := hash(password)
	encPw := encode(hashedPw)

	if encPw != expEncPw {
		fmt.Println("Passwords don't match")
	} else {
		fmt.Println("Everything is fine")
	}

	fmt.Printf("base64 (exp):\t%s\n", expEncPw)
	fmt.Printf("base64 (enc):\t%s\n", encPw)
}
