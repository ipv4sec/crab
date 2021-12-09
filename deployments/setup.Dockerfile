FROM centos:7


WORKDIR /app
ADD assets/bin/kubectl /usr/local/bin/kubectl
ADD assets/bin/istioctl /usr/local/bin/istioctl

RUN chmod +x /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/istioctl

ADD assets/istio/ assets/istio/
ADD assets/island/ assets/island/
ADD assets/setup/ assets/setup/

ADD scripts/istio.sh scripts/istio.sh
ADD scripts/label.sh scripts/label.sh

RUN chmod +x scripts/istio.sh
RUN chmod +x scripts/label.sh

COPY assets/bin/setup setup

CMD ["./setup"]