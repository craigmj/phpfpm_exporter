#!/bin/bash
set -e
export GOPATH=`pwd`
if [[ ! -d bin ]]; then 
	mkdir bin
fi
go build -o bin/phpfpm_exporter src/cmd/phpfpm_exporter.go
