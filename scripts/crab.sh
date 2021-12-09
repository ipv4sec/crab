#!/bin/bash

domain=''

while [ $# -gt 0 ]; do
  case "$1" in
  --domain)
    domain="$2"
    shift
    ;;
  --*)
    echo "Illegal option $1"
    ;;
  esac
  shift $(($# > 0 ? 1 : 0))
done

cat <<EOF | kubectl apply -f -

apiVersion: v1
kind: Namespace
metadata:
  name: island-system
  labels:
    istio-injection: enabled
---

apiVersion: v1
kind: ServiceAccount
metadata:
  name: island-system
  namespace: island-system
---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: island-system
subjects:
  - kind: ServiceAccount
    name: island-system
    namespace: island-system
    apiGroup: ""
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
---

apiVersion: batch/v1
kind: Job
metadata:
  name: island-setup
  namespace: island-system
spec:
  template:
    metadata:
      name: island-setup
    spec:
      containers:
        - name: main
          image: harbor1.zlibs.com/island/island-setup:0.1
          imagePullPolicy: Always
          env:
            - name: CRAB_DOMAIN
              value: $domain
      restartPolicy: OnFailure
      serviceAccountName: island-system
---

apiVersion: v1
kind: Service
metadata:
  name: crab
  namespace: island-system
spec:
  selector:
    app: island
    component: ui
  ports:
    - port: 80
      targetPort: 3000
  type: NodePort

---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: island-system
  namespace: island-system
spec:
  {}
---

apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: island-system-to-island-system
  namespace: island-system
spec:
  action: ALLOW
  rules:
    - from:
        - source:
            principals:
              - cluster.local/ns/island-system/sa/island-system
  selector:
    matchLabels:
      app: island
---
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: crab
  namespace: island-system
spec:
  gateways:
    - island-system/crab
  hosts:
    - '*'
  http:
    - name: router
      route:
        - destination:
            host: crab
            port:
              number: 80
          headers:
            request:
              add:
                X-Forwarded-Host: crab.$domain
---

apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: crab
  namespace: island-system
spec:
  selector:
    istio: ingressgateway
  servers:
    - hosts:
        - crab.$domain
      port:
        name: http
        number: 80
        protocol: HTTP
EOF

echo
echo "================================================================================"
echo
echo "系统正在部署中, 之后请访问: http://carb."$domain
echo
echo "注意: 默认登录用户名为root, 默认密码为toor, 登录成功后请尽快修改此密码"
echo
echo "================================================================================"
echo