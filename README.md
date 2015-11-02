# riak-ping

## Overview
Check riak connection status.

## Requirement
- [mattn/gom](https://github.com/mattn/gom)
- [tpjp/goriakpbc](https://github.com/tpjg/goriakpbc)

## Usage
```
Usage of riak-ping:
  -i="127.0.0.1:8087": Riak Server IP Address and Port
  -p="TCP": Protocol (TCP/PB)
```

Check TCP connection.
```
riak-ping
```

Check Riak Protocol Buffer connection.
```
riak-ping -p PB
```

Check remote server.
```
riak-ping -i "127.0.0.1:8087" -p "TCP"
```

## Install
### Install gom
```
go get github.com/mattn/gom
```

### create Gomfile
```
gom gen gomfile
```
### install requirment package
```
gom install
```

### riak-ping package test (support local only)
```
gom test -v
```

### riak-ping package build (for local)
```
gom build
```

### riak-ping package build (for linux)
```
GOOS=linux GOARCH=amd64 gom build
```

## Contribution
1. Pull or Clone
1. Create a feature branch
1. Commit your changes
1. Run test suite with the `gom test -v` command and confirm that it passes
1. Create new Pull Request
