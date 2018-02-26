#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o handleRequest handleRequest.go
zip handleRequest.zip handleRequest
sls deploy