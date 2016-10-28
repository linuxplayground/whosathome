# Arp scanner written in golang.

* The mac addresses you are interessted in are hardcoded into the script.
* The constant, DECAY is for how long you want to wait before marking someone as OFFLINE.
* The constant, INTERVAL is a text representation (see time.ParseDuration) between polling.

## Dependancies
* arp-scan

## Usage
`sudo ./arp &>> arp.log`

## Build
``` bash
$ cd $GOPATH
$ mkdir -pv src/github.com/linuxplayground/
$ cd src/github.com/linuxplayground/
$ git clone https://github.com/linuxplayground/arp.git
$ cd $GOPATH
$ go build github.com/linuxplayground/arp
```

