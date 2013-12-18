#!/bin/bash
set -e

# run build first as normal user, then invoke tests as root
# because we're running integration tests against docker.sock
# which requires root access

if [ "$1" == "doTests" ]; then
	export GOPATH=$2

	go test -i ./
	go test -v ./
else
	go build
	sudo $0 doTests $GOPATH
fi
