#!/bin/bash


while [ $# -gt 0 ]; do
  case "$1" in
  --domain)
    domain="$2"
    shift
    ;;
  esac
  case "$1" in
  --webssh)
    webssh="$2"
    shift
    ;;
  esac
  shift $(($# > 0 ? 1 : 0))
done

if [ x$domain == x ]
then
  echo "Missing domain."
  exit 0
fi

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: island-system
---

apiVersion: v1
kind: ServiceAccount
metadata:
  name: crab
  namespace: island-system
---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: crab
subjects:
  - kind: ServiceAccount
    name: crab
    namespace: island-system
    apiGroup: ""
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
---

apiVersion: v1
kind: ConfigMap
metadata:
  name: island-plugin
  namespace: island-system
data:
  description: plugin
  webssh: '$webssh'
---

apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: island-system
value: 1000000000
globalDefault: false
description: "island-system"
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
        - name: island-setup
          image: harbor1.zlibs.com/island/island-setup:alpha
          imagePullPolicy: Always
          env:
            - name: ISLAND_DOMAIN
              value: $domain
            - name: ISLAND_WEBSSH
              value: '$webssh'
      restartPolicy: OnFailure
      serviceAccountName: crab
EOF

echo
echo "================================================================================"
echo
echo "系统正在部署中, 之后请访问: http://crab."$domain
echo
echo "注意: 默认登录用户名为root, 默认密码为toor, 登录成功后请尽快修改此密码"
echo
echo "================================================================================"
echo