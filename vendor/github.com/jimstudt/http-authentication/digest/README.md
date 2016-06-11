# auth/digest [![GoDoc](https://godoc.org/github.com/jimstudt/http-authentication/digest?status.png)](http://godoc.org/github.com/jimstudt/http-authentication/digest)

Authenticate using HTTP Digest Authentication for [Martini](https://github.com/codegangsta/martini)
or similar HTTP request routers.

## Overview

HTTP _Digest Authentication_ is one of the mechanisms described in 
[RFC-2617](http://pretty-rfc.herokuapp.com/RFC2617), along with
_Basic Authentication_. Digest Authentication is more secure because the password is 
never transmitted across the network and need not even be stored on the servers. Despite this
advantage, Digest Authentication is less popular than Basic and less well supported.

This package implements the server side authentication of Digest Authentication using
client provided credentials, Apache-style htdigest files, or a custom client provided
mechanism.

As might be dreaded, Digest authentication comes with a suite of options to be combined
into implementation taxing cross products. Not all of Digest Authentication
is supported by this package. It stops implementing where libcurl and most web browsers stop
on the theory that you'd have no one to talk to anyway.

| Feature | Support |
|---------|---------|
| algorithm | asks for MD5-sess, accepts MD5 or MD5-sess or none |
| qop | asks for auth, accepts auth or none, **no auth-int** |

## Usage

This package is designed to be used by Martini or a similar URL router. It provides a .ServeHTTP()
method which matches the http.ServeHTTP. On successful authentication it does nothing to the response.
With missing, invalid, or uncheckable credentials to responds with the appropriate HTTP status and
a brief diagnostic body. If this happens the request has been fully satisfied and you should 
stop handling the request.

To specify your account information in a map, one might do this...

~~~ go
import (
  "github.com/codegangsta/martini"
  "github.com/jimstudt/http-authentication/digest"
  "log"
)

func main() {
  m := martini.Classic()
  myUserStore := digest.NewSimpleUserStore(  map[string]string{
			"foo:mortal": "3791e8e14a10b3666ba15d9e78e4b359",    // pw is 'bar'
			"Mufasa:testrealm@host.com": "939e7578ed9e3c518a452acee763bce9",   // pw is 'Circle Of Life'
                 })

  digester := digest.NewDigestHandler( "mortal", nil, nil, myUserStore )
  m.Use( digester.ServeHTTP )   // this will force authentication of all requests, you can be more specific.

  //...
}
~~~

To use an Apache style htdigest file, one would instead...

~~~ go
     m := martini.Classic()
     ...
     // Read file: hint, the nil is standing in for a malformed line reporter function.
     myUserFile,err := digest.NewHtdigestUserStore("path/to/my/htdigest/file", nil)
     if err != nil {
		log.Fatalf("Unable to load password file: %s", err.Error())
     }
     ...
     myUserFile.ReloadOn(syscall.SIGHUP, nil)   // optional, that nil is the bad line reporter again
     ...
     digester := digest.NewDigestHandler( "My Realm", nil, nil, myUserFile )
     ...
     m.Post("/my-sensitive-uri", digester.ServeHTTP, mySensitiveHandler)  // just protect this one, notice chained handlers.
~~~

## Documentation

The API documentation is available using godoc and at [godoc.org](http://godoc.org/github.com/jimstudt/http-authentication/digest)




