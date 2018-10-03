package main

import (
	"crypto/sha512"
	"encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func encode(password string) string {
	pwSha512 := sha512.New()
	pwSha512.Write([]byte(password))
	encodedPassword := base64.StdEncoding.EncodeToString(pwSha512.Sum(nil))
	return encodedPassword
}

func isPortValid(port string) bool {
	const minPort = 1024
	const maxPort = 9000

	intPort, _ := strconv.Atoi(port)
	var portIsValid = true

	if intPort < minPort || intPort > maxPort {
		portIsValid = false
	}

	return portIsValid
}

func getRandomInt() int {
	const minimum = 5000
	const maximum = 6000

	rand.Seed(time.Now().Unix())
	return rand.Intn(maximum-minimum) + minimum
}

func getSleepInterval() time.Duration {
	const nanosecondsPerMillisecond = 1000000

	randSleep := getRandomInt()
	randSleep = randSleep * nanosecondsPerMillisecond
	return time.Duration(randSleep)
}

func passwordHandler(writer http.ResponseWriter, request *http.Request) {
	time.Sleep(getSleepInterval())

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("400 - Request method must be 'POST'"))
	} else {
		clearPassword := request.FormValue("password")
		encodedPassword := encode(clearPassword)
		encodedPassword = "\"" + encodedPassword + "\"\n"
		writer.Write([]byte(encodedPassword))
	}
}

func main() {
	fmt.Println("Hasher")
	const defaultPort string = "8080"

	portPtr := flag.String("port", defaultPort, "Listen port")
	flag.Parse()

	if !isPortValid(*portPtr) {
		fmt.Printf("Port %s is invalid; setting to default (%s)\n", *portPtr, defaultPort)
		*portPtr = defaultPort
	}

	*portPtr = ":" + *portPtr

	http.HandleFunc("/hash", passwordHandler)
	http.ListenAndServe(*portPtr, nil)
}
