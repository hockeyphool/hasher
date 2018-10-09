# hasher
Hasher starts a webserver listening on the specified port. The server handles the following connection endpoints:

**"/hash"** - Accept a POST form with one field named "password". Compute the SHA512 hash of the password value, and return it in Base64-encoded form. The password value must be less than or equal to 32 characters in length.

**"/shutdown"** - Shut down the server gracefully, ensuring that all in-flight requests are allowed to complete, but no new requests are accepted.

## USAGE
`Usage of ./hasher:
  -port string
    	Listen port (allowable range: "1024" - "9000" (default "8080")`

Valid 'port' values: **1024** - **9000**. The lower limit ensures that the web server can open its listening port without colliding with well-known ports (and thus requiring superuser privileges).

The upper limit was arbitrarily chosen for this exercise.

## BUILD AND EXECUTE
To build and run the application:
`$ go build
$ ./hasher`

In a separate window, send requests using **'curl'**:
`$ curl -i --data "password=testpassword" http://localhost:8080`

To test sending multiple requests:
`$ for PASSWD in testpassword1 testpassword2 testpassword3 testpassword4
do
  curl -i --data "password=${PASSWD}" http://localhost:8080
  sleep 2
done`

To shutdown the server and exit the application, send this
request:
`$ curl- -i http://localhost:8080/shutdown`

Combine the last two steps in separate terminal windows:
1. Submit multiple requests
1. Send the shutdown request before all requests are submitted
Confirm the following:
* Requests submitted **before** the shutdown receive a _200 OK_ status and a **Base64-encoded** password string
* Requests submitted **after** the shutdown either
** receive a _403 Forbidden_ status and the message **403 - Server is shutting down**, or
** receive a _Connection refused_ error at the command line because the server has already shut down
