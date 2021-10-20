
# Overview

- [说明](#说明)

<a name="说明"></a>
## 说明

在现有集群上执行
`curl -fsSL http://island-resource.develenv.com/crab.sh`

再执行
`kubectl get svc -n island-system`, 查看暴露的端口号


|  status   | 意义  |
|  ----  | ----  |
| 0  | 未部署 |
| 1  | 正在部署中  |
| 2  | 部署完成  |
| 3  | 卸载中  |
| 4  | 卸载完成  |