FROM centos:7

WORKDIR /app
COPY assets/bin/kubectl /usr/local/bin/
RUN chmod 755 /usr/local/bin/kubectl

ADD assets/runtime/ assets/runtime/

ADD config.yaml config.yaml
COPY assets/bin/scheduler scheduler

CMD ["./scheduler"]