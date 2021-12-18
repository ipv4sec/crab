#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o akali main.go

docker build -t harbor1.zlibs.com/island/akali:alpha .
docker push harbor1.zlibs.com/island/akali:alpha