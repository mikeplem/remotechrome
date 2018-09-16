#!/bin/bash
env GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o remotechrome_arm main.go
