FROM golang:1.16 as build

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app/

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -o crab cmd/setup/main.go

FROM centos:7 as prod


WORKDIR /app
ADD assets/bin/istioctl /usr/local/bin/istioctl

ADD assets/istio/operator.yaml assets/istio/operator.yaml
ADD assets/island/island-nginx.yaml assets/island/island-nginx.yaml

COPY --from=build /app/crab .

CMD ["./crab"]