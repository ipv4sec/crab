FROM golang:1.16 as build

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app/


ADD go.mod .
ADD go.sum .
RUN go mod download
ADD assets/bin/cue /usr/local/bin/cue

RUN chmod +x /usr/local/bin/cue

COPY . .
RUN go build -o parser cmd/parser/main.go
CMD ["./parser"]