FROM centos:7

WORKDIR /app
ADD assets/bin/cue /usr/local/bin/cue
RUN chmod +x /usr/local/bin/cue

ADD assets/workloads/ assets/workloads/

COPY assets/bin/parser parser

CMD ["./parser"]