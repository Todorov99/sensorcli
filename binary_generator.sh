#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o sensorcli_linux .
GOOS=windows GOARCH=amd64 go build -o sensorcli_windows.exe .
GOOS=darwin GOARCH=amd64 go build -o sensorcli_darwin .