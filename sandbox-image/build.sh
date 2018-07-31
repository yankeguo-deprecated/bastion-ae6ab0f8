#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o dummy-init dummy-init.go

docker build -t bastion-sandbox .