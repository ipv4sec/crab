#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zed main.go

docker build -t harbor1.zlibs.com/tars/zed:alpha .
docker push harbor1.zlibs.com/tars/zed:alpha