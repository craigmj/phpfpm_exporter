#!/bin/bash
set -e
export GOPATH=`pwd`
if [[ ! -d bin ]]; then 
	mkdir bin
fi
go build -o bin/phpfpm_exporter cmd/phpfpm_exporter.go
