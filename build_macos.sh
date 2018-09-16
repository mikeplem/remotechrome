#!/bin/bash
env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o remotechrome_mac main.go
