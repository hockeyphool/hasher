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
	fmt.Println("Hasher")

	password := "angryMonkey"

	hashedPw := hash(password)
	encPw := encode(hashedPw)

	fmt.Printf("base64 (enc):\t%s\n", encPw)
}
