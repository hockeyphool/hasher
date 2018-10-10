/*
Hasher implements a webserver that listens on the specified port (default: 8080).
It handles the following endpoints:
   /hash     - Accept a form field value named "password", compute its SHA512 hash, and return the Base64-encoded string representing the hash
   /stats    - Return a JSON object containing two fields:
               total:   The total number of requests handled by the server since startup
               average: The average duration of requests in microseconds
   /shutdown - Stop accepting new requests, allow in-flight requests to complete, the shut down the server gracefully
 */
package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type serverStats struct {
	NumEncodings    int   `json:"total"`
	AverageDuration int64 `json:"average"`
	mux             sync.RWMutex
}

var (
	processTransactions bool
	stats               serverStats
)

func init() {
	processTransactions = true
	rand.Seed(time.Now().Unix())
	initStats()
}

func initStats() {
	stats.NumEncodings = 0
	stats.AverageDuration = 0
}

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

	return rand.Intn(maximum-minimum) + minimum
}

func getSleepInterval() time.Duration {
	const nanosecondsPerMillisecond = 1000000

	randSleep := getRandomInt()
	randSleep = randSleep * nanosecondsPerMillisecond
	return time.Duration(randSleep)
}

func updateStats(duration time.Duration) {
	stats.mux.Lock()
	defer stats.mux.Unlock()

	stats.NumEncodings++

	durationInUSecs := int64(duration / time.Microsecond)

	if stats.NumEncodings == 1 {
		stats.AverageDuration = durationInUSecs
	} else {
		stats.AverageDuration = (stats.AverageDuration + durationInUSecs) / 2
	}
}

func marshalStats() []byte {
	stats.mux.RLock()
	defer stats.mux.RUnlock()

	jsonStats, err := json.Marshal(&stats)
	if err != nil {
		log.Fatal(err)
	}
	return jsonStats
}

func buildShutdownHandler(hdlrWaitGroup *sync.WaitGroup, quit chan<- bool) http.Handler {
	shutdownFunc := func(respWriter http.ResponseWriter, _ *http.Request) {
		processTransactions = false
		respWriter.WriteHeader(http.StatusOK)
		respWriter.Write([]byte("200 - Shutting down\n"))
		hdlrWaitGroup.Wait()
		marshalStats()
		quit <- true
	}
	return http.HandlerFunc(shutdownFunc)
}

func buildPasswordHandler(hdlrWaitGroup *sync.WaitGroup) http.Handler {
	passwordFunc := func(respWriter http.ResponseWriter, req *http.Request) {
		proceed := make(chan bool)
		hdlrWaitGroup.Add(1)
		go func() {
			var (
				status  = http.StatusOK
				message = ""
			)
			const passwordKey = "password"
			const maxPasswordLength = 32

			defer hdlrWaitGroup.Done()
			startTime := time.Now()

			if processTransactions {
				time.Sleep(getSleepInterval())
				if req.Method != http.MethodPost {
					message = "400 - Request method must be 'POST'\n"
					status = http.StatusBadRequest
				} else {
					if len(req.FormValue(passwordKey)) <= maxPasswordLength {
						clearPassword := req.FormValue(passwordKey)
						message = encode(clearPassword)
						message = "\"" + message + "\""
						respWriter.Header().Set("Content-Transfer-Encoding", "BASE64")
					} else {
						status = http.StatusBadRequest
						message = "400 - Password must be < " + strconv.Itoa(maxPasswordLength) + " characters\n"
					}
				}
			} else {
				status = http.StatusForbidden
				message = "403 - Server is shutting down\n"
			}

			respWriter.WriteHeader(status)
			respWriter.Write([]byte(message))

			endTime := time.Now()
			updateStats(endTime.Sub(startTime))

			proceed <- true
		}()
		<-proceed
	}
	return http.HandlerFunc(passwordFunc)
}

func buildStatsHandler(hdlrWaitGroup *sync.WaitGroup) http.Handler {
	statsFunc := func(respWriter http.ResponseWriter, _ *http.Request) {
		proceed := make(chan bool)
		hdlrWaitGroup.Add(1)

		go func() {
			var (
				status  = http.StatusOK
				message = []byte("")
			)
			defer hdlrWaitGroup.Done()

			if processTransactions {
				time.Sleep(getSleepInterval())
				message = marshalStats()
				respWriter.Header().Set("Content-Type", "application/json")

			} else {
				status = http.StatusForbidden
				message = []byte("403 - Cannot process stats - server shutting down")
			}
			respWriter.WriteHeader(status)
			respWriter.Write(message)

			proceed <- true
		}()
		<-proceed
	}
	return http.HandlerFunc(statsFunc)
}

func main() {
	fmt.Println("Hasher")
	const defaultPort = "8080"

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	portPtr := flag.String("port", defaultPort, "Listen port (allowable range: \"1024\" - \"9000\"")
	flag.Parse()

	if !isPortValid(*portPtr) {
		fmt.Printf("Port %s is invalid; setting to default (%s)\n", *portPtr, defaultPort)
		*portPtr = defaultPort
	}

	*portPtr = ":" + *portPtr

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    *portPtr,
		Handler: mux,
	}

	var wg sync.WaitGroup
	done := make(chan bool)
	quit := make(chan bool, 1)

	go func() {
		<-quit
		log.Println("Server is shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not shut down gracefully: %v\n", err)
		}
		close(done)
	}()

	srvPassHandler := buildPasswordHandler(&wg)
	mux.Handle("/hash", srvPassHandler)

	srvShutdownHandler := buildShutdownHandler(&wg, quit)
	mux.Handle("/shutdown", srvShutdownHandler)

	srvStatHandler := buildStatsHandler(&wg)
	mux.Handle("/stats", srvStatHandler)

	log.Printf("Starting server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-done
	log.Printf("serverStats: %+v", &stats)
	log.Printf("Finished")
}
