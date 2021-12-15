
# Overview

### 安装

在现有集群上执行
`curl -fsSL http://island-resource.develenv.com/crab.sh | bash`

### 流程

1. 浏览器访问系统后台, 调用 [用户信息] 接口
2. 默认用户名密码登录
3. 登陆后的默认页面是 [应用管理]
4. 只要访问 [应用管理] 页面, 则请求 [获取集群根域] 判断接口返回是否是默认值(默认值为空), 默认值时弹窗提示跳转(提示信息, 调转路径需要前端定义)
5. 分割线
6. 删除之前的查看日志, 隐藏根域设置, 上传zip包改为上传yaml
7. 调整实例列表字段, 删除状态字段, 查看 [实例列表] 接口
8. [应用管理]页面点击[导出配置], node端请求 [实例详情] 接口(需修改)取 deployment字段, node构造前端进行下载
9. [应用管理]页面点击[导出部署], 请求 [实例详情] 接口取 configuration 字段, node构造前端进行下载
10. [应用管理]页面点击[查看部署], node端请求 [实例详情] 接口, 取 pods 字段
11. 在 [查看部署] 页面, 点击具体的POD, 请求 [实例日志] 接口
12. 在[工作负载]选项页, 分别请求的接口是 [Trait列表] 接口, [Type列表] 接口, [Vendor列表] 接口
13. 在[工作负载]选项页, 要编辑时, 分别请求  [Trait修改] 接口, [Type修改] 接口, [Vendor修改] 接口
14. 在[创建应用]选项页, 点击[下载], 前端构造为 manifest.yaml 文件下载
15. 在[创建应用]选项页, 点击[部署], node端构造zip后请求 [添加实例] 接口
16. 在[创建Trait]选项页, 点击[保存], 请求 [创建Trait] 接口
17. 在[创建WorkloadType]选项页, 点击[保存], 请求 [创建WorkloadType] 接口
18. 在[创建WorkloadVendor]选项页, 请求 [SystemSpec默认值] 接口
19. 在[创建WorkloadVendor]选项页, 点击[保存], 请求 [创建WorkloadVendor] 接口
20. 在[创建WorkloadVendor]选项页, 点击[转换], 请求 [转换YAML到CUE] 接口
21. 在[创建WorkloadVendor]选项页, 点击[检查], 请求 [检查CUE语法] 接口
22. 原样转发 [部署] 接口

### 默认值

在[创建应用]选项页, Metadata默认值
```yaml
apiVersion: aam.globalsphare.com/v1alpha1
kind: Application
metadata:
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
type: globalsphare.com/v1alpha1/workloadType/webservice
vendor: globalsphare.com/v1alpha1/workloadVendor/webservice
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
    description: init user description
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

dependencies默认值
```yaml
    - name: gitlab
      version: ">=0.0.1"
      location: user-defined(https://gitlab.com)
      items:
        /*:
          - create
          - read
          - update
          - delete
```

exports默认值
```yaml
exports:
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
apiVersion: aam.globalsphare.com/v1alpha1
kind: WorkloadVendor
metadata:
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

