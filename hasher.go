package main

import (
	"crypto/sha512"
	"encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var globalWaitGroup sync.WaitGroup

func encode(password string) string {
	pwSha512 := sha512.New()
	pwSha512.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(pwSha512.Sum(nil))
}

func isPortValid(port string) bool {
	const minPort = 1024
	const maxPort = 9000

	intPort, err := strconv.Atoi(port)
	var portIsValid = true

	if err != nil || intPort < minPort || intPort > maxPort {
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

type handler func(writer http.ResponseWriter, request *http.Request, wgRef *sync.WaitGroup)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r, &globalWaitGroup)
}

func shutdownHandler(writer http.ResponseWriter, _ *http.Request, wgRef *sync.WaitGroup) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("200 - Shutting Down"))
	wgRef.Wait()
	os.Exit(0)
}

func passwordHandler(writer http.ResponseWriter, request *http.Request, wgRef *sync.WaitGroup) {
	const passwordKey = "password"
	const maxPasswdLen = 32
	const validationFailureStatus = http.StatusBadRequest

	fmt.Println("Entering passwordHandler")

	wgRef.Add(1)

	go func() {
		defer wgRef.Done()

		fmt.Println("Entering passwordHandler goroutine")

		time.Sleep(getSleepInterval())

		if request.Method != http.MethodPost {
			writer.WriteHeader(validationFailureStatus)
			writer.Write([]byte("400 - Request method must be 'POST'\n"))
		} else {
			if len(request.FormValue(passwordKey)) <= maxPasswdLen {
				clearPassword := request.FormValue(passwordKey)
				encodedPassword := encode(clearPassword)
				encodedPassword = "\"" + encodedPassword + "\"\n"
				writer.Write([]byte(encodedPassword))
			} else {
				writer.WriteHeader(validationFailureStatus)
				passWordTooLongMsg := "400 - Password must be < " + strconv.Itoa(maxPasswdLen) + " characters\n"
				writer.Write([]byte(passWordTooLongMsg))
			}
		}
	}()
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

	mux := http.NewServeMux()
	mux.Handle("/hash", handler(passwordHandler))
	mux.Handle("/shutdown", handler(shutdownHandler))

	//http.HandleFunc("/hash", passwordHandler)
	//http.HandleFunc("/shutdown", shutdownHandler)
	http.ListenAndServe(*portPtr, mux)
}
