#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

docker build -t harbor1.zlibs.com/crab-demo/hezhongjiang/redis-test .
docker push harbor1.zlibs.com/crab-demo/hezhongjiang/redis-test