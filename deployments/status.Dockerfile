FROM centos:7

WORKDIR /app
COPY config.yaml config.yaml
COPY assets/bin/status status

ENTRYPOINT ["./status"]