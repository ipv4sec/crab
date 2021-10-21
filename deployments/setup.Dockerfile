FROM centos:7


WORKDIR /app
ADD assets/bin/kubectl /usr/local/bin/kubectl
ADD assets/bin/istioctl /usr/local/bin/istioctl

RUN chmod +x /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/istioctl

ADD assets/istio/ assets/istio/
ADD assets/ceph/ assets/ceph/
ADD assets/island/ assets/island/

ADD scripts/ceph.sh scripts/ceph.sh
ADD scripts/istio.sh scripts/istio.sh
RUN chmod +x scripts/ceph.sh
RUN chmod +x scripts/istio.sh

COPY assets/bin/setup setup

CMD ["./setup"]