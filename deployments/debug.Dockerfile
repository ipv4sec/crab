FROM centos:7

WORKDIR /app
COPY assets/bin/debug debug

CMD ["./debug"]