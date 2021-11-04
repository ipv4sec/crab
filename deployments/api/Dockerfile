FROM centos:7

WORKDIR /app

COPY assets/bin/kubectl /usr/local/bin/
RUN chmod 755 /usr/local/bin/kubectl

ADD deployments/api/api.ini /etc/supervisord.d/api.ini

ADD config.yaml config.yaml
COPY assets/bin/api api

RUN yum update -y && yum install epel-release -y && yum -y install supervisor && yum clean all
CMD ["bash", "-c", "/usr/bin/supervisord;while true;do sleep 100;done"]