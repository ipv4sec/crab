#!/bin/bash

rm assets/bin/parser
rm assets/bin/setup
rm assets/bin/status
rm assets/bin/web

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/parser cmd/parser/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/setup cmd/setup/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/web cmd/web/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o assets/bin/status cmd/status/main.go