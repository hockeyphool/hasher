# hasher
Hasher starts a webserver listening on the specified port. The server handles the following connection endpoints:

**"/hash"** - Accept a POST form with one field named "password". Compute the SHA512 hash of the password value, and return it in Base64-encoded form. The password value must be less than or equal to 32 characters in length. **NOTE:** The password length limit is arbitrary. No lower limit is enforced.

**"/shutdown"** - Shut down the server gracefully, ensuring that all in-flight requests are allowed to complete, but no new requests are accepted.

**"/stats"** - Get the server's current statistics. Returns a JSON object with two fields:
* total   - the total number of requests handled since the server started
* average - the average request duration in microseconds

## USAGE
```
Usage of ./hasher:
 -port string
   	Listen port (allowable range: "1024" - "9000" (default "8080")
```

Valid 'port' values: **1024** - **9000**. The lower limit ensures that the web server can open its listening port without colliding with well-known ports (and thus requiring superuser privileges).

The upper limit was arbitrarily chosen for this exercise.

## BUILD AND EXECUTE
To build and run the application:

```
$ go run hasher.go
```

## UNIT TESTING
```
$ go test -v
=== RUN   TestEncoding
--- PASS: TestEncoding (0.00s)
=== RUN   TestPortValidation
--- PASS: TestPortValidation (0.00s)
=== RUN   TestRandomInt
--- PASS: TestRandomInt (0.00s)
=== RUN   TestSleepInterval
--- PASS: TestSleepInterval (0.00s)
=== RUN   TestInitStats
--- PASS: TestInitStats (0.00s)
=== RUN   TestUpdateStats
--- PASS: TestUpdateStats (0.00s)
=== RUN   TestMarshalStats
--- PASS: TestMarshalStats (0.00s)
=== RUN   TestPasswordHandlerGoodStatus
--- PASS: TestPasswordHandlerGoodStatus (5.26s)
=== RUN   TestPasswordHandlerBadStatus
--- PASS: TestPasswordHandlerBadStatus (0.00s)
=== RUN   TestShutdownHandler
--- PASS: TestShutdownHandler (0.00s)
=== RUN   TestStatsHandlerGoodStatus
--- PASS: TestStatsHandlerGoodStatus (0.00s)
=== RUN   TestStatsHandlerBadStatus
--- PASS: TestStatsHandlerBadStatus (0.00s)
PASS
ok  	_/home/scottt/projects/go/hasher	5.268s
```
The app and unit tests pass _**golint**_ and _**go vet**_.

## APPLICATION TESTING
Once you've started the server, you can send requests in a separate window using **'curl'**:
```
$ curl -i --data "password=testpassword" http://localhost:8080/hash
HTTP/1.1 200 OK
Content-Transfer-Encoding: BASE64
Date: Wed, 10 Oct 2018 02:18:32 GMT
Content-Length: 90
Content-Type: text/plain; charset=utf-8

"6eYzCXq5zrPkjsP3DuK+ukHQXVQg7+5dqF+X2XAFcnWH/aM+9P8jIgiPTHnoEzzJzZ81EvTTowPL21vFhUFaAA=="
```

To retrieve server statistics after processing some number of requests:
```
$ curl -i http://localhost:8080/stats && printf "\n"
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 10 Oct 2018 02:16:37 GMT
Content-Length: 29

{"total":2,"average":5342445}
```

To test sending multiple requests:
```
$ for PASSWD in testpassword1 testpassword2 testpassword3 testpassword4
do
  time curl -i --data "password=${PASSWD}" http://localhost:8080/hash
# sleep between requests to ensure they overlap and timing for each request is clear
  sleep 2
done
```

To shutdown the server and exit the application, send this
request:
```
$ curl -i http://localhost:8080/shutdown
HTTP/1.1 200 OK
Date: Wed, 10 Oct 2018 02:20:06 GMT
Content-Length: 20
Content-Type: text/plain; charset=utf-8

200 - Shutting down
```

Combine the last two steps in separate terminal windows:
1. Submit multiple requests
1. Send the shutdown request before all requests are submitted

Confirm the following:
* Requests submitted **before** the shutdown receive a _200 OK_ status and a **Base64-encoded** password string
* Requests submitted **after** the shutdown receive either
	* a _403 Forbidden_ status and the message **403 - Server is shutting down**, or
	* a _Connection refused_ error at the command line because the server has already shut down
