FROM centos:7

WORKDIR /app
COPY assets/bin/kubectl /usr/local/bin/
RUN chmod 755 /usr/local/bin/kubectl

COPY assets/bin/cue /usr/local/bin/
RUN chmod 755 /usr/local/bin/cue

COPY deployments/graceful.sh /app
COPY deployments/api.sh /app
COPY deployments/parser.sh /app

RUN chmod +x /app/graceful.sh
RUN chmod +x /app/api.sh
RUN chmod +x /app/parser.sh

ADD config.yaml config.yaml
COPY assets/bin/api api
COPY assets/bin/parser parser

CMD [ "./graceful.sh" ]