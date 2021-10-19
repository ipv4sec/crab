
# Overview

- [说明](#说明)

<a name="说明"></a>
## 说明

在现有集群上执行
`curl -fsSL http://island-resource.develenv.com/crab.sh | bash -s -- --domain example.com --password toor`

再执行
`kubectl get svc -n island-system`, 查看暴露的端口号