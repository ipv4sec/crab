#!/bin/bash

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
      restartPolicy: OnFailure
      serviceAccountName: crab
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
EOF