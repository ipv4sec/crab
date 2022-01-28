FROM centos:7

WORKDIR /app
COPY main app
RUN yum -y update && yum -y install git

CMD ["./app"]