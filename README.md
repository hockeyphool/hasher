# hasher
Hasher starts a webserver listening on the specified port. The
server handles the following connection endpoints:

"/hash" - Accept a POST form with one field named "password". 
Compute the SHA512 hash of the password value, and return it in
Base64-encoded form.

"/shutdown" - Shut down the server gracefully, ensuring that all
in-flight requests are allowed to complete, but no new requests
are accepted.


Usage of ./hasher:
  -port string
    	Listen port (allowable range: "1024" - "9000" (default "8080")

Valid 'port' values: 1024 - 9000. The lower limit ensures that the
web server can open its listening port without colliding with
well-known ports (and thus requiring superuser privileges).

The upper limit was arbitrarily chosen for this exercise.
