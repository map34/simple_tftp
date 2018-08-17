# Simple TFTP

[![Build Status](https://travis-ci.org/map34/simple_tftp.svg?branch=master)](https://travis-ci.org/map34/simple_tftp)

A simple in-memory TFTP server written in Golang. Just open up your favorite tftp client interface,
and start saving your data on the TFTP server.

## Get
```bash
go get github.com/map34/simple_tftp
cd $GOPATH/src/github.com/map34/simple_tftp
glide install
```

or

```bash
mkdir -p $GOPATH/src/github.com/map34/simple_tftp
cd $GOPATH/src/github.com/map34/simple_tftp
cp -r <location_of_project>/* .
glide install
```

If you want to install the server as a binary, just
do:

```bash
cd  $GOPATH/src/github.com/map34/simple_tftp
go install .
```

## Go
To start the TFTP server.
``` bash
cd  $GOPATH/src/github.com/map34/simple_tftp
go run main.go
```

Or after installation, from the CLI,
just run

``` bash
simple_tftp
```
you will see an interface such as this:
``` bash
INFO[0000] Listening UDP at [::]:52767
```
which means that the server is listening for
incoming requests. To start storing to / reading from the server from the CLI:
``` bash
$ tftp
tftp> binary
tftp> connect localhost 52767
tftp> put some_file.txt
Sent 3230 bytes in 0.0 seconds
tftp> get some_file.txt
Received 3230 bytes in 0.0 seconds
```

## Test
Unit test uses [testify](https://github.com/stretchr/testify) for assertion tests.
``` bash
cd  $GOPATH/src/github.com/map34/simple_tftp
go test tftputils/*
```



