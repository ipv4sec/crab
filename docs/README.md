
# Overview

### 安装

在现有集群上执行
`curl -fsSL http://island-resource.develenv.com/crab.sh | bash -s -- --domain {root-domain}`

### 流程

1. 浏览器访问系统后台, 调用 [用户信息] 接口
2. 默认用户名密码登录
3. 登陆后的默认页面是 [应用管理]
4. 只要访问 [应用管理] 页面, 则请求 [获取集群根域] 判断接口返回是否是默认值(默认值为空), 默认值时弹窗提示跳转(提示信息, 调转路径需要前端定义)
5. ================分割线, 接口文档还未完成, 以下是需要修改或新增的内容================
6. 删除之前的实例列表中的[查看日志]
7. 隐藏[根域设置], 隐藏[工作负载源]的设置
8. 修改[上传]为[添加应用]
9. 调整实例列表字段, 删除[状态]字段
10. [应用管理]页面点击[部署链接], 请求 [实例详情] 接口取 id 字段 拼接到原样转发 [流水线接口] 接口之后 , 前端复制到粘贴板并给出提示(或者前端提示用户要复制)
11. [应用管理]页面点击[导出K8S描述文件], node端请求 [实例详情] 接口, 取 deployment 字段, 前端下载
12. [应用管理]页面点击[部署详情], node端请求 [实例详情] 接口, 取 details 字段, 见下文 [Details  字段的解释]
13. 在 [Pod详情] 中有 [实例日志] 接口
14. 在[工作负载]选项页, 分别请求的接口是 [Trait列表] 接口, [WorkloadType列表] 接口, [WorkloadVendor列表] 接口
15. 在[工作负载]选项页, 要编辑时, 分别请求  [Trait修改] 接口, [WorkloadType修改] 接口, [WorkloadVendor修改] 接口
16. 在[创建应用]选项页, 点击[下载], 前端构造为 manifest.yaml 文件下载
17. 在[创建应用]选项页, 点击[部署], node端构造zip后请求 [添加实例] 接口
18. 在[创建Trait]选项页, 点击[保存], 请求 [创建Trait] 接口
19. 在[创建WorkloadType]选项页, 点击[保存], 请求 [创建WorkloadType] 接口
20. 在[创建WorkloadVendor]选项页, 请求 [SystemSpec默认值] 接口
21. 在[创建WorkloadVendor]选项页, 点击[保存], 请求 [创建WorkloadVendor] 接口
22. 在[创建WorkloadVendor]选项页, 点击[转换], 请求 [转换YAML到CUE] 接口
23. 在[创建WorkloadVendor]选项页, 点击[检查], 请求 [检查CUE语法] 接口
24. 原样转发 [流水线接口] 接口

### Details  字段的解释

在 details 字段中有:
1.CronJob列表, 字段名为 cronJob
2.DaemonSet列表, 字段名为 daemonSet
3.Deployment列表, 字段名为 deployment
4.Job列表, 字段名为 job
5.Pod列表, 字段名为 pod
6.ReplicaSet列表, 字段名为 replicaSet 
7.ReplicationController列表, 字段名为 replicationController
8.StatefulSet列表, 字段名为 statefulSet
9.Service列表, 字段名为 service
10.ConfigMap列表, 字段名为 configMap
11.PVC列表, 字段名为 pvc
12.Secret列表, 字段名为 secret
13.RoleBinding列表, 字段名为 roleBinding
14.Role列表, 字段名为 role
15.ServiceAccount列表, 字段名为 serviceAccount

可再次请求 /{字段名}/{名称} 查看详情, 具体看 [资源详情] 接口




### 默认值

保持缩进, 暂定2个空格

在[创建应用]选项页, Metadata默认值
```yaml
  name: example
  version: 0.0.1
  description: 样例应用
  keywords:
    - 样例应用
  author: example@example.com
  maintainers:
    - email: example@example.com
      name: example
      web: https://example.com
  repositories: ["https://github.com/example/example.git"]
  bugs: https://github.com/example/example/issues
  licenses:
    - type: LGPL
      url: https://license.spec.com

```

Workloads默认值
```yaml
name: example
type: webservice
vendor: webservice
properties:
    image: nginx:1.21
traits:
    - type: expose
      properties:
        k1: "v1"
```

Userconfigs默认值
```yaml
    "$schema": http://json-schema.org/draft-07/schema#
    "$id": http://example.com/product.schema.json
    title: User
    description: User Description
    type: object
    properties:
      username:
        type: string
      password:
        type: string
    required:
      - username
      - password
```

dependencies默认值(数组)
```yaml
  name: gitlab
  version: ">=0.0.1"
  location: user-defined(https://gitlab.com)
  items:
    /*:
      - create
      - read
      - update
      - delete
```

exports默认值(map类型)
```yaml
/user:
  - create
  - read
  - update
  - delete
/admin:
  - create
  - read
  - update
  - delete
```

在[创建Trait]选项页, 默认值
```yaml
apiVersion: aam.globalsphare.com/v1alpha1
kind: Trait
metadata:
  name: expose
spec:
  parameter: |
    k1: *"v1" | string

```

在[创建WorkloadType]选项页, 默认值
```yaml
apiVersion: aam.globalsphare.com/v1alpha1
kind: WorkloadType
metadata:
  name: example
spec:
  parameter: |
    image: *"example" | string
```

在[创建WorkloadVendor]选项页, Metadata默认值
```yaml
name: example
```

SystemSpec默认值, 请求 [Spec默认值] 接口

YAML默认值
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

Template默认值: 无

