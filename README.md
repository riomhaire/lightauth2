# LightAuth2

This project _SHOULD NOT_ be used in production for various reasons. This project is one of our many experiments in learning GO and in this case [clean architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) , REST, Crypto and JWT for a demo Authentication and Session application. We prefer to learn using something more than the usual 'hello world'. This is a re-write of https://github.com/riomhaire/lightauth .

In essence the API allows a UI or other API to authenticate a user with username/password (passwords are hashed so no leak) and returns a 'token' which can be used and submitted when it can be queried for validity or checked for claims (if 'admin' token is also supplied).

All communication is done via json, REST and HTTP.

## Installation

The simplest way if you 'make' installed is to run 'make' which will install all the dependencies and install the apps. Otherwise After cloning the repository you need to install the few dependencies. Execute the following within the main directory.

```bash
$ go get github.com/riomhaire/lightauth2
$ cd <gopath-root>/src/github.com/riomhaire/lightauth2
$ go get ./...
```

There are three applications - the server is in the 'lightauth2' directory, a session generation application (useful for generating API tokens with roles such as admin), and a 'user' app for creating users suitable for including in the users.csv file (similar to passwd).

The best way if you just want run is to build and install the app:

```bash
 go install github.com/riomhaire/lightauth2/frameworks/application/lightauth2
```

## Getting Started

The lightauth2 application now consists on one application (was three: server, user and session) which understands four commands:

```bash
$ lightauth2 -h
lightauth2 is an authentication server and session key generator

Usage:
  lightauth2 [flags]
  lightauth2 [command]

Available Commands:
  help        Help about any command
  serve       Starts an authentication server
  session     Create a session key suitable for accessing a resource.
  user        Create a user record and signed/encrypted password info for including in your user store.
  version     Returns application version

Flags:
  -h, --help   help for lightauth2

Use "lightauth2 [command] --help" for more information about a command.
```

### User Creation Application

The 'users.csv' file is read from the folder where the 'lightauth2' server is executed. This is a simple CSV file where the 1st line is a header consisting of:

```csv
username,password,enabled,roles,claim1,claim2
```

An example is:

```csv
username,password,enabled,roles,claim1,claim2
test,939c1f673b7f5f5c991b4d4160642e72880e783ba4d7b04da260392f855214a6,true,none,8f0daa02-4f89-4ad8-9db9-d1d8bcc60f07,d70e6e06-a8b7-4daf-a2d0-476acd396e82
admin,50b911deac5df04e0a79ef18b04b29b245b8f576dcb7e5cca5937eb2083438ba,true,admin:api,83df72c5-f3e7-47a4-bd1b-a8d7e8f6db0f,90a08fec-783f-48b0-b2ac-325f71f6ebfb
```

The password is hashed based on a secret and a salt. To add a user you need to use the lightauth2 'user' command which takes parameters and creates a line suitable to append to the user csv file:

```bash
$lightauth2 user -h

Usage:
  lightauth2 user [flags]

Flags:
      --claim1 string     Security Claim - Eg QRCode Hash
      --claim2 string     Security Claim - Eg API Key
  -h, --help              help for user
  -p, --password string   Password (in the raw - will be encoded).
  -r, --roles string      List of roles separated by ':'. (default "guest:public")
  -u, --user string       Username associated with the token. (default "anonymous")
```

### Session Creation Application

If a session token is created via the authenticate method they have a limited life span (usually 3600 seconds) before they become invalid. Tokens for api's typically have a requirement for a longer lived period - sometimes months or longer. Long lived Tokens need not be stored within the sessions file since they only need to be encoded using the same parameters as used by the lightauth server itself.

The session token creation is done via the session command of lightauth:

```bash
$ lightauth2 session -h
Create a session key suitable for accessing a resource.

Usage:
  lightauth2 session [flags]

Flags:
  -h, --help                help for session
  -r, --roles string        List of roles separated by ':'. (default "guest:public")
  -s, --secret string       Secret used to sign things. (default "secret")
  -k, --sessionPeriod int   How long session  will be active for. (default 3600)
  -t, --token string        If populated means decode token
  -u, --user string         Username associated with the token. (default "anonymous")
```

An example usage would be:

```bash
 $ lightauth2 session --user someapp --roles "api:admin" --sessionPeriod 9999999 --secret hush

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTUxNTYzMTIsImppZCI6IjJiM2ZhYjcwLTA3M2MtNGRiNi05ZTEwLThlOWJlMTQwZWM5NCIsInJvbGVzIjpbImFwaSIsImFkbWluIl0sInN1YiI6InNvbWVhcHAifQ.dmsHOMzspru-HBL7QsuLILhFuEOlNSXMksVUismFi8U
```

Created sessions are only valid for lightauth servers which have been started with the same secret.

### The LightAuth2 Server

The lightauth2 is a single server which includes authentication and session token management. In a production system these should be implemented as separate servers - but this is a learning experience and we dont plan to do this at the moment.

The server can be started with the following parameters:

```bash
$lightauth2 serve -h
Starts an authentication server on the given port and using given
               secret based on users stored in a predefined location - by default a csv file.

Usage:
  lightauth2 serve [flags]

Flags:
  -c, --consul                 Enable consul support
  -t, --consulHost string      Host where consul resides usually something like http://consul:8500
  -h, --help                   help for serve
  -l, --loggingLevel string    Logging Level: Trace,Debug,Info,Warn,Error (default "Debug")
  -p, --port int               Default Port to Listen to. (default 3030)
      --profile                Enable Profiling.
      --serverCert string      Server SSL Cert File  (default "server.crt")
      --serverKey string       Server SSL Key File  (default "server.key")
  -k, --sessionPeriod int      How long session returned will be active for. (default 3600)
  -s, --sessionSecret string   Secret used to sign things. (default "secret")
      --useSSL                 Use SSL
      --useUserAPI             Use User API (see https://githib.com/riomhaire/lightauthuserapi
      --userAPIHost string     The base endpoint where service resides
      --userAPIKey string      The API Access Token Needed for User API
  -u, --usersFile string       If User File used this is the one to use. (default "users.csv")
```

The parameters are pretty much self evident. An example startup would produce:

```bash
$ lightauth2  serve --usersFile ../../ansible/users.csv
2017/10/11 18:45:05 [INFO] Initializing
2017/10/11 18:45:05 Reading User Database ../../ansible/users.csv
2017/10/11 18:45:05 #Number of users = 2
2017/10/11 18:45:05 [INFO] Running LightAuth2 Version 1.8.1
2017/10/11 18:45:05 [INFO]
CONFIGURATION
          SigningSecret : 'secret'
           TokenTimeout : '3600'
                  Store : '../../ansible/users.csv'
                    SSL : 'false'
                   Port : '3030'

[negroni] listening on :3030
```

## The API

The API is a fairly simple one and consists of:

1.  Authenticate/Login.
2.  Verify a session token.
3.  Get session token details (get roles).
4.  Get Call Statistics.

The endpoint in the default startup can be found at "http://somehost:3030/api/v2/authentication" or "http://somehost:3030/api/v2/authentication/session" or "http://somehost:3030/api/v1/authentication/statistics" - and yes we know we should be using HTTPS or some other transport medium, but this a simple project to help us learn GO and not a prod app.

Content-Type should be "application/json"

### Authenticate/Login.

There are two flavours of this - Authenticate via password and Authenticate via claims,

Example request for authenticate via password POST ("http://somehost:3030/api/v2/authentication"):

```json
{
  "username": "vader",
  "password": "anakin"
}
```

Wont mention password is in clear - should be a hash or some other method in conjunction with https... but this is a learning experience app.

Example response on success:

```json
{
  "token":
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDc3NDc3MDEsImppZCI6IjBiNWQ3M2FmLWYwMmQtNHMxYy1iMjVkLTIyODc2MzNkM2E4ZSIsInJvbGVzIjpbIm5vbmUiXSwic3ViIjoidGVzdCJ9.6cCDtJw50nu4bfUG4KSrgm0P8pZlxYn8dxsZn1_NACY"
}
```

Example response on error is the useual HTTP Status Error 'No Authorized':

```
Response headers:
HTTP/1.1 401 Unauthorized
Content-Type: application/json; charset=UTF-8
Date: Wed, 11 Oct 2017 17:49:59 GMT
Content-Length: 22
```

Example request for authenticate via claims POST ("http://somehost:3030/api/v2/authentication"):

```json
{
  "username": "vader",
  "claims": ["wef33rfefefef", "3535wefw4ff3"]
}
```

Claims are information about a user which can be used to identify them. Strictly speaking a password is a claim - but the common usage is the user/password hence the via password route above. In this case a claim could be a registered 'QRCode' or 'Iris Scan Info' etc. The response and token returned is consistent across all calls.

### Verify a session token.

This is simply send the token as a bearer token within the HEADER to the validate endpoint:

Example GET request ("http://somehost:3030/api/v2/authentication/session") ..
Where you send:

```
Connection: keep-alive
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDc3NDgwMTMsImppZCI6ImFhMAJmNTYxLTI3ZTEtNDE4My1hYTEwLWVjYThkNmU4MGE4ZSIsInJvbGVzIjpbIm5vbmUiXSwic3ViIjoidGVzdCJ9.chkjOBBnmyjx6pTbmcqSTn9enMm6Cy9B-1n2g8MgLxA
Content-Type: application/json
Content-Length: 0
Host: somehost:3030
User-Agent: Apache-HttpClient/4.5.2 (Java/1.8.0_131)
```

You will get with status code 200 (OK) if the token is valid an error otherwise.

### Get session token details (get roles).

Do a get call with the Authorization header set to the appropriate token from Authenticate

Example request ("http://somehost:3030/api/v2/authentication/session/decoder"):

```
Connection: keep-alive
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDc3NTI5MTksImppZCI6JmFiYjUwYjdjLTk0YWEtNDM3YS1iMzcyLTJiOTkzOGRjYzc3MCIsInJvbGVzIjpbImFkbWluIiwiYXBpIl0sInN1YiI6ImFkbWluIn0.qX-tdFbGZAyGW7nBx6Br6xwn7oEXdpXSqLuCka5pDbM
Content-Type: application/json
Content-Length: 0
Host: localhost:3030
User-Agent: Apache-HttpClient/4.5.2 (Java/1.8.0_131)
```

Example response is successful:

```json
{
  "id":
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5NzkyMDYsImppZCI6IjNmYWM5MjU2LTNjNTEtNGM5OC05YzZlLWU1MjA1NGMzYzIyZSIsInJvbGVzIjpbInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.4bQh09BPOjaavzhvErIb008Ot6STyd2B-ZoXta-7g2Y",
  "user": "stormtrooper",
  "expires": 1504979206,
  "roles": ["security"]
}
```

On Error the appropriate Status code will be returned.

### Get Call Metrics.

Example GET request ("http://somehost:3030/api/v2/authentication/metrics"):

No Bearer Authorization header is needed. If you make a request with 'application/json' accept type you will get:

```json
{
  "pid": 8943,
  "uptime": "3m16.264197174s",
  "uptime_sec": 196.264197174,
  "time": "2017-09-18 19:51:19.820078277 +0100 IST m=+196.265651036",
  "unixtime": 1505760679,
  "status_code_count": {},
  "total_status_code_count": {
    "200": 1,
    "400": 3
  },
  "count": 0,
  "total_count": 4,
  "total_response_time": "687.008µs",
  "total_response_time_sec": 0.000687008,
  "average_response_time": "171.752µs",
  "average_response_time_sec": 0.000171752
}
```

If you want a prometheus acceptable response, then use the "text/plain" accept type and you will get something like:

```text
# HELP lightauth2_uptime_sec How many seconds app has been up.
# TYPE lightauth2_uptime_sec counter
lightauth2_uptime_sec 17.717113684

# HELP lightauth2_total_response_time_sec Total time spent in handling requests.
# TYPE lightauth2_total_response_time_sec counter
lightauth2_total_response_time_sec 0.000838896

# HELP lightauth2_average_response_time_sec Average time spent in handling requests.
# TYPE lightauth2_average_response_time_sec guage
lightauth2_average_response_time_sec 7.6263e-05

# HELP lightauth2_response_status_200 Total Number of Requests returning http status 200
# TYPE lightauth2_response_status_200 counter
lightauth2_response_status_200 4

# HELP lightauth2_response_status_404 Total Number of Requests returning http status 404
# TYPE lightauth2_response_status_404 counter
lightauth2_response_status_404 7

# HELP lightauth2_response_total_count Total Number of Requests.
# TYPE lightauth2_response_total_count counter
lightauth2_response_total_count 11
```

## SSL

As of 0.4 Support for SSL server has been added based on information at that excellent resource https://gist.github.com/6174/9ff5063a43f0edd82c8186e417aae1dc and is enabled via three command line variables:

- useSSL - set that to 'true' EG '-useSSL true'
- serverCert - contain the name of the file containing the SSL cert to use
- serverKey - contain the name of the file containing the server key.

For self signed certificates you can use the following steps to generate them:

```bash
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key

openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

## Performance

The following graphs are based on running the enclosed JMeter script against LightAuth2 deployed within six nodes of a [ Raspberry PI cluster ](http://www.riomhaire.com/post/theign-post/) fronted by a Linkerd router.

### Transaction Summary

![Transactions Summary](docs/stats/20171011/summary.png)

### Transactions Per Second

![Transactions Per Second](docs/stats/20171011/tps.png)

### Response Distributions

![Response Distributions](docs/stats/20171011/response-distributions.png)

### Response Percentiles

![Response Percentiles](docs/stats/20171011/response-percentiles.png)

### Latencies

![Latencies](docs/stats/20171011/latencies.png)

### Profiling

![Profiling](docs/stats/20171029/profiling.png)
