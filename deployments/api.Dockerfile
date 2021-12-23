FROM centos:7

WORKDIR /app
COPY assets/bin/kubectl /usr/local/bin/
RUN chmod 755 /usr/local/bin/kubectl

COPY assets/bin/cue /usr/local/bin/
RUN chmod 755 /usr/local/bin/cue

COPY deployments/run.sh /app
RUN chmod +x /app/run.sh

COPY assets/cue/ assets/cue/

ADD config.yaml config.yaml
COPY assets/bin/api api
COPY assets/bin/parser parser

CMD [ "./run.sh" ]