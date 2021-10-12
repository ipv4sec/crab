FROM golang:1.16 as build

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app/

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -o crab cmd/web/main.go

FROM centos:7 as prod

WORKDIR /app
COPY --from=build /app/crab .

CMD ["./crab"]