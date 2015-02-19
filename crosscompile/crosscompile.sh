#!/bin/bash

function trap_handler()
{
	set +x
	LINE="$1"
	ERR="$2"
	echo "BUILD FAILED: line ${LINE}: exit status of last command: ${ERR}"
	exit
}
trap 'trap_handler ${LINENO} $?' ERR

export GOPATH=/usr/
go get -v ./...

export GOOS=windows
export GOARCH=386
go build -v -o geoip-$GOOS-$GOARCH.exe
export GOARCH=amd64
go build -v -o geoip-$GOOS-$GOARCH.exe
export GOOS=linux
export GOARCH=386
go build -v -o geoip-$GOOS-$GOARCH
export GOARCH=amd64
go build -v -o geoip-$GOOS-$GOARCH
export GOARCH=arm
go build -v -o geoip-$GOOS-$GOARCH
export GOOS=darwin
export GOARCH=386
go build -v -o geoip-$GOOS-$GOARCH
export GOARCH=amd64
go build -v -o geoip-$GOOS-$GOARCH
export GOOS=freebsd
go build -v -o geoip-$GOOS-$GOARCH
export GOOS=netbsd
go build -v -o geoip-$GOOS-$GOARCH
export GOOS=dragonfly
go build -v -o geoip-$GOOS-$GOARCH
