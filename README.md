# Diagnostic Data Generation Server

- **Author**: Steve Carpenter
- **Version**: 1.0.0

## Overview
This is a simple HTTP server written in GoLang that supports generating random data for various
device types following a GET request to it's `/usage` endpoint. Currently, the devices supported
include the following:
- CPU
- Disk
- Memory
- Network

The server supports generating data that spans an hour or a week in increments of 10 seconds or
10 minutes respectively.

## Getting Started
- Clone the repository to your local directory from [here](https://github.com/TimeSeriesDataApp/server)
- Move the cloned directory to the Go workspace on the current machine
	- Many times, this is some form of `$GOPATH/src/github.com/<username-dir>/`
	- Note, this assumes you have already configured your `$GOPATH` and have the `go tool` installed
- Install the dependencies for this project.
	- `go get github.com/gorilla/handlers`
	- `go get github.com/gorilla/mux`
	- `go get github.com/joho/godotenv`
- Add a `.env` file with the specified port the server should listen on (see example below)
	- For this example port 3000 is being used. Replace this with the port number specified in the `.env` file if a different port is used
- Once finished build the code with `go build` and run the resulting binary file
- The server should be listening on `http://localhost:<port-number-from-env-file>` for requests

`.env`
```
PORT=3000
```

## Making Requests
This server supports `GET` requests to the `/usage` endpoint. Information passed to the server is
done so for this `GET` endpoint in query strings. The two query string variables are as follows:

`duration`:
- `hr` value denotes an hour of data
- `wk` value denotes a week of data

`device`:
- The four device arguments are as follows: `cpu`, `disk`, `memory`, `network`
- These can be combined in any way as one query string separated by commas
- e.g. `http://localhost:3000/usage?duration=wk&device=cpu,memory,network`

Here is an example of a request to get CPU
usage data for an hour. HTTP request examples here are demonstrated using the [HTTPie](https://httpie.org/) tool.

```
ENDPOINT: http://localhost:3000/usage?duration=hr&device=cpu

Example:

http get :3000/usage duration==hr device==cpu
HTTP/1.1 200 OK
Content-Length: 964
Content-Type: application/json
Date: Wed, 02 May 2018 06:37:43 GMT

{
    "cpu": [
        {
            "toffset": 0,
            "usage": 22
        },
        {
            "toffset": 10,
            "usage": 13
        },
        {
            "toffset": 20,
            "usage": 4
        },
        {
            "toffset": 30,
            "usage": 7
        },
        {
            "toffset": 40,
            "usage": 9
        },
...
...
...
	]
}
```

The data is returned back from the request as JSON with each device followed by
an array of time slice data. `toffset` is the time offset and `usage` are the usage
values for this data. For a valid `GET` request, at minimum a duration must be
specified and at least one device.

The following example shows requesting a week of usage data for multiple devices.

```
ENDPOINT: http://localhost:3000/usage?duration=wk&device=cpu,disk,network,memory

Example:

http get :3000/usage duration=wk device==cpu,disk,network,memory
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 02 May 2018 06:45:24 GMT
Transfer-Encoding: chunked

{
    "cpu": [
        {
            "toffset": 0,
            "usage": 28
        },
        {
            "toffset": 10,
            "usage": 31
        },
...
...
...
	],
    "disk": [
        {
            "toffset": 0,
            "usage": 37
        },
        {
            "toffset": 10,
            "usage": 37
        },
...
...
...
	],
    "memory": [
        {
            "toffset": 0,
            "usage": 14
        },
        {
            "toffset": 10,
            "usage": 22
        },
...
...
...
	],
    "network": [
        {
            "toffset": 0,
            "usage": 32
        },
        {
            "toffset": 10,
            "usage": 26
        },
...
...
...
	]
}
```
