#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

docker build -t harbor1.zlibs.com/tars/echo .
docker push harbor1.zlibs.com/tars/echo