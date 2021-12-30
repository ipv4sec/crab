
# Overview

- [说明](#说明)
- [获取管理员信息](#获取管理员信息)
- [修改管理员密码](#修改管理员密码)
- [添加实例](#添加实例)
- [运行实例](#运行实例)
- [删除实例](#删除实例)
- [获取节点地址](#获取节点地址)
- [获取集群根域](#获取集群根域)
- [设置集群根域](#设置集群根域)
- [设置工作负载源](#设置工作负载源)
- [查询工作负载源](#查询工作负载源)

- [实例列表](#实例列表)
- [实例详情](#实例详情)
- [实例日志](#实例日志)

- [流水线接口](#流水线接口)
- [Trait列表](#Trait列表)
- [WorkloadType列表](#WorkloadType列表)
- [WorkloadVendor列表](#WorkloadVendor列表)

- [修改Trait](#修改Trait)
- [修改WorkloadType](#修改WorkloadType)
- [修改WorkloadVendor](#修改WorkloadVendor)

- [删除Trait](#删除Trait)
- [删除WorkloadType](#删除WorkloadType)
- [删除WorkloadVendor](#删除WorkloadVendor)

- [添加Trait](#添加Trait)
- [添加WorkloadType](#添加WorkloadType)
- [添加WorkloadVendor](#添加WorkloadVendor)

- [SystemSpec默认值](#SystemSpec默认值)
- [转换YAML到CUE](#转换YAML到CUE)
- [检查CUE语法](#检查CUE语法)

- [资源详情](#资源详情)

- [Trait详情](#Trait详情)
- [WorkloadType详情](#WorkloadType详情)
- [WorkloadVendor详情](#WorkloadVendor详情)



<a name="说明"></a>
## 说明
### 返回的数据结构如下
```
{
    "code": 错误代码,
    "result": 任意类型
}
```

其中 `code` 表示的为错误代码, 数字类型

正常返回时`code`为`0`, `result`为数据, 可能是数字, 布尔, 字符串, 数组, 对象等等

错误返回时`code`不为`0`, `result`为错误信息, 字符串类型, 可直接显示在浏览器页面

错误代码原则上前端用不到, 前端仅需要判断非0时显示`result`字段即可

接口请求和输出均为JSON格式, 接口的字段为小驼峰命名, 特有名词除外(注意:当前特有名词有`userconfigs`)

<a name="获取管理员信息"></a>
## 获取管理员信息
### 请求语法
```
GET /user/root HTTP/1.1
```
### 请求参数
无

### 返回值
```json
{
    "code": 0,
    "result": {
      "password": "toor",
      "username": "root"
    }
}
```

<a name="修改管理员密码"></a>
## 修改管理员密码
### 请求语法
```
PUT /user/root HTTP/1.1
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|password|密码|无|是|
|oldPassword|密码|无|是|
### 请求示例
```json
{
  "password": "admin233",
  "oldPassword": "toor"
}
```
### 返回值
```json
{
    "code": 0,
    "result": "设置成功"
}
```

<a name="添加实例"></a>
## 添加实例
### 请求语法
```
POST /app HTTP/1.1
Content-Type: multipart/form-data; 
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|file|实例描述文件（即manifest.yaml压缩后的zip文件）|无|是|
### 返回值
```json
{
  "code": 0,
  "result": {
    "id": "ins1635146904",
    "dependencies": {
      "github": {
        "instances": [
          {
            "id": "ins1634971790",
            "name": "github"
          }
        ],
        "location": "https://www.github.com",
        "type": "immutable"
      }
    },
    "userconfigs": {
      "properties": {
        "param1": {
          "type": "integer"
        },
        "param2": {
          "type": "string"
        },
        "param3": {
          "properties": {
            "param3_1": {
              "type": "number"
            },
            "param3_2": {
              "type": "number"
            }
          },
          "required": [
            "param3_1"
          ],
          "type": "object"
        },
        "param4": {
          "items": {
            "type": "string"
          },
          "minItems": 1,
          "type": "array",
          "uniqueItems": true
        }
      },
      "required": [
        "param2"
      ],
      "type": "object"
    }
  }
}
```

该接口的错误返回, 按照全局说明中的错误返回格式处理

<a name="实例列表"></a>
## 实例列表
### 请求语法
```
GET /app?limit=<limit>&offset=<offset> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|offset| |0|否|
|limit|  |10|否|

### 返回值
```json
{
  "code": 0,
  "result": {
    "rows": [
      {
        "id": "ins1634971791",
        "name": "harbor",
        "version": "2.0.0",
        "entry": "http://ins1634971791.example.com",
        "created_at": "2021-10-23T06:49:51.498Z",
        "updated_at": "2021-10-23T06:49:51.498Z"
      }
    ],
    "total": 1
  }
}
```

<a name="实例详情"></a>
## 实例详情
### 请求语法
```
GET /app/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id| 实例主键 |无|是|

### 返回值

[点击查看](app.json)


cronJob展示的字段有:
名称, 创建时间, 最后执行时间
metadata.name, metadata.creationTimestamp, status.lastScheduleTime

daemonSet展示的字段有:
名称, 创建时间, 当前可用数
metadata.name, metadata.creationTimestamp, status.numberAvailable

deployment展示的字段有:
名称, 创建时间, 当前可用数
metadata.name, metadata.creationTimestamp, status.availableReplicas

job展示的字段有:
名称, 创建时间, 完成时间
metadata.name, metadata.creationTimestamp, status.completionTime

pod展示的字段有:
名称, 创建时间, 状态, 启动时间
metadata.name, metadata.creationTimestamp, status.phase, status.startTime

replicaSet展示的字段有:
名称, 创建时间, 当前可用数
metadata.name, metadata.creationTimestamp, status.availableReplicas

replicationController展示的字段有:
名称, 创建时间, 当前副本数
metadata.name, metadata.creationTimestamp, status.replicas

statefulSet展示的字段有:
名称, 创建时间, 当前副本数
metadata.name, metadata.creationTimestamp, status.currentReplicas

service展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

configMap展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

pvc展示的字段有:
名称, 创建时间, 状态
metadata.name, metadata.creationTimestamp, status.phase

secret展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

roleBinding展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

role展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

serviceAccount展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

<a name="实例日志"></a>
## 实例日志
### 请求语法
```
GET /app/<id>/logs HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|pod| 实例主键 |无|是|

### 返回值


```json
{
  "code": 0,
  "result": [{
    "name": "POD名称",
    "value": "春江潮水连海平，海上明月共潮生"
  }]
}
```


<a name="运行实例"></a>
## 运行实例
### 请求语法
```
PUT /app/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id| 实例主键 |无|是|

以下参数为BODY参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|dependencies|依赖描述 |无|否|
|userconfigs| 实例自定义配置信息 |无|否|
|status| 实例状态,此处固定为1 |无|是|
#### dependencies 示例
```json
[
    {
		"name": "demo-app1",
		"id": "iqtxycne"
	},
	{
		"name": "demo-app1",
		"location": "https://www.huanqiu.com"
	}
]
```
#### userconfigs 示例
```json
{
    "param1": "aaa",
    "param2": 123,
    "param3": {
        "param3_1": "bbb",
        "param3_2": 456
    }
}
```
### 请求示例
```json
{
  "status": 1,
  "dependencies": [
    {
      "name": "demo-app1",
      "id": "iqtxycne"
    },
    {
      "name": "demo-app1",
      "location": "https://www.huanqiu.com"
    }
  ],
  "userconfigs": {
    "param1": "aaa",
    "param2": 123,
    "param3": {
      "param3_1": "bbb",
      "param3_2": 456
    }
  }3
}
```
### 返回值
```json
{
    "code": 0,
    "result": "正在部署中"
}
```

<a name="删除实例"></a>
## 删除实例
### 请求语法
```
DELETE /app/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id|应用实例 id|无|是|
### 返回值
```json
{
    "code": 0,
    "result": "删除成功"
}
```

<a name="获取节点地址"></a>
## 获取节点地址
显示所有节点的IP地址(仅能显示此节点网卡上绑定的IP地址)

### 请求语法
```
GET /cluster/addrs HTTP/1.1
```
### 请求参数
无

### 返回值
```json
{
    "code": 0,
    "result": [
        {
            "name": "master1",
            "addrs": [
              "192.168.0.1",
              "192.168.0.2",
              "192.168.0.3"
            ]
        },
        {
            "name": "salve1",
            "addrs": [
              "192.168.0.4"
            ]
        },
        {
            "name": "salve2",
            "addrs": [
              "192.168.0.5",
              "192.168.0.6"
            ]
        }
    ]
}
```

<a name="获取集群根域"></a>
## 获取集群根域

### 请求语法
```
GET /cluster/domain HTTP/1.1
```
### 请求参数
无

### 返回值
```json
{
    "code": 0,
    "result": "example.com"
}
```

<a name="设置集群根域"></a>
## 设置集群根域

### 请求语法
```
PUT /cluster/domain HTTP/1.1
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|domain|绑定到此集群的域名|无|是|
#### 请求参数示例
```json
{
    "domain": "abc.com"
}
```

### 返回值
```json
{
    "code": 0,
    "result": {
      "status": 3,
      "message": "成功"
    }
}
```

`result.status` 为设置根域后的状态码, 数字类型, 具体意义见下表
`result.message` 为当前的状态信息, 字符串,  可直接显示在网页

|  status   | 意义  |
|  ----  | ----  |
| 0  | 检测域名的解析失败 |
| 1  | 检测域名的解析成功, 保存失败 |
| 2  | 检测域名的解析成功, 保存成功 |


<a name="设置工作负载源"></a>
## 设置工作负载源

### 请求语法
```
PUT /cluster/mirror HTTP/1.1
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|mirror|作负载源|无|是|
#### 请求参数示例
```json
{
    "mirror": "https://github.com/xxx.git"
}
```

### 返回值
```json
{
    "code": 0,
    "result": "设置成功"
}
```

<a name="查询工作负载源"></a>
## 查询工作负载源

### 请求语法
```
GET /cluster/mirror HTTP/1.1
```
### 请求参数
无

### 返回值
```json
{
    "code": 0,
    "result": "https://github.com/xxx.git"
}
```

<a name="流水线接口"></a>
## 流水线接口

### 请求语法
```
PUT /deployment/<id> HTTP/1.1
Content-Type: multipart/form-data; 
```
### 请求参数

以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id| 实例主键 |无|是|

以下参数为BODY参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|file|应用描述文件|无|是|


### 返回值
```json
{
    "code": 0,
    "result": "部署成功"
}
```


<a name="Trait列表"></a>
## Trait列表
### 请求语法
```
GET /trait?limit=<limit>&offset=<offset> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|offset| |0|否|
|limit|  |10|否|

### 返回值

type: 0 内置, 不可删除 1 可删除
表头为: 主键, 名称, 版本, 值, 类型, 创建时间, 修改时间

```
{
  "code": 0,
  "result": {
    "rows": [
      {
        "id": 1, 主键
        "name": "ingress",  名称
        "apiVersion": "aam.globalsphare.com/v1alpha1", 版本
        "value": "具体定义", 值
        "type": 0, 类型
        "created_at": "2021-10-23T06:49:51.498Z",
        "updated_at": "2021-10-23T06:49:51.498Z"
      }
    ],
    "total": 1
  }
}
```


<a name="WorkloadType列表"></a>
## WorkloadType列表
### 请求语法
```
GET /workloadType?limit=<limit>&offset=<offset> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|offset| |0|否|
|limit|  |10|否|

### 返回值

type: 0 内置, 不可删除 1 可删除

```json
{
  "code": 0,
  "result": {
    "rows": [
      {
        "id": 1,
        "name": "worker",
        "apiVersion": "aam.globalsphare.com/v1alpha1",
        "value": "具体定义",
        "type": 0,
        "created_at": "2021-10-23T06:49:51.498Z",
        "updated_at": "2021-10-23T06:49:51.498Z"
      }
    ],
    "total": 1
  }
}
```


<a name="WorkloadVendor列表"></a>
## WorkloadVendor列表
### 请求语法
```
GET /workloadVendor?limit=<limit>&offset=<offset> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|offset| |0|否|
|limit|  |10|否|

### 返回值

type: 0 内置, 不可删除 1 可删除

```json
{
  "code": 0,
  "result": {
    "rows": [
      {
        "id": 1,
        "name": "webservice",
        "apiVersion": "aam.globalsphare.com/v1alpha1",
        "value": "具体定义",
        "yaml": "",
        "cue": "",
        "type": 0,
        "created_at": "2021-10-23T06:49:51.498Z",
        "updated_at": "2021-10-23T06:49:51.498Z"
      }
    ],
    "total": 1
  }
}
```

<a name="修改Trait"></a>
## 修改Trait
### 请求语法
```
PUT /trait/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id|主键|无|是|

以下参数为BODY参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|value| 具体定义 |无|是|

### 返回值
```json
{
    "code": 0,
    "result": "修改成功"
}
```

<a name="修改WorkloadType"></a>
## 修改WorkloadType
### 请求语法
```
PUT /workloadType/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id|主键|无|是|

以下参数为BODY参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|value| 具体定义 |无|是|

### 返回值
```json
{
    "code": 0,
    "result": "修改成功"
}
```


<a name="修改WorkloadVendor"></a>
## 修改WorkloadVendor
### 请求语法
```
PUT /workloadVendor/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id|主键|无|是|

以下参数为BODY参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|value| 最终的值 |无|是|
|yaml| YAML的值, 页面上YAML框的值 |无|是|
|cue| CUE的值, 即translate spec框的值 |无|是|

### 返回值
```json
{
    "code": 0,
    "result": "修改成功"
}
```

<a name="删除Trait"></a>
## 删除Trait
### 请求语法
```
DELETE /trait/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id|主键|无|是|
### 返回值
```json
{
    "code": 0,
    "result": "删除成功"
}
```

<a name="删除WorkloadType"></a>
## 删除WorkloadType
### 请求语法
```
DELETE /workloadType/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id| 主键|无|是|
### 返回值
```json
{
    "code": 0,
    "result": "删除成功"
}
```

<a name="删除WorkloadVendor"></a>
## 删除WorkloadVendor
### 请求语法
```
DELETE /workloadVendor/<id> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|id| 主键|无|是|
### 返回值
```json
{
    "code": 0,
    "result": "删除成功"
}
```


<a name="SystemSpec默认值"></a>
## SystemSpec默认值
### 请求语法
```
GET /tool/systemTemplate HTTP/1.1
```

### 请求参数
无

### 返回值
```json
{
    "code": 0,
    "result": "模板的内容"
}
```

<a name="添加Trait"></a>
## 添加Trait
### 请求语法
```
POST /trait HTTP/1.1 
```
### 请求参数

```json
{
  "value": "具体的yaml"
}

```


### 返回值
```json
{
  "code": 0,
  "result": "创建成功"
}
```


<a name="添加WorkloadType"></a>
## 添加WorkloadType
### 请求语法
```
POST /workloadType HTTP/1.1 
```
### 请求参数

```json
{
  "value": "具体的yaml"
}

```


### 返回值
```json
{
  "code": 0,
  "result": "创建成功"
}
```

<a name="添加WorkloadVendor"></a>
## 添加WorkloadVendor
### 请求语法
```
POST /workloadVendor HTTP/1.1
```
### 请求参数

```json
{
  "yaml": "编写的原生YAML格式的K8S文件内容",
  "cue":  "点击转换之后的CUE格式的文件内容",
  "value": "最终保存的文件内容"
}
```


### 返回值
```json
{
  "code": 0,
  "result": "创建成功"
}
```


<a name="转换YAML到CUE"></a>
## 转换YAML到CUE
### 请求语法
```
POST /tool/convertion HTTP/1.1 
```
### 请求参数

```json
{
  "value": "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-deployment\n  labels:\n    app: nginx\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n        - name: nginx\n          image: nginx:1.14.2\n          ports:\n            - containerPort: 80"
}

```


### 返回值
```json
{
  "code": 0,
  "result": "翻译后的CUE模板或者为报错的信息, 直接展示即可, 让用户判断"
}
```



<a name="检查CUE语法"></a>
## 检查CUE语法
### 请求语法
```
POST /tool/spelling HTTP/1.1 
```
### 请求参数

```json
{
  "value": "deployment: \"nginx-deployment\": {\n\tapiVersion: \"apps/v1\"\n\tkind:       \"Deployment\"\n\tmetadata: {\n\t\tname: \"nginx-deployment\"\n\t\tlabels: app: \"nginx\"\n\t}\n\tspec: {\n\t\treplicas: 3\n\t\tselector: matchLabels: app: \"nginx\"\n\t\ttemplate: {\n\t\t\tmetadata: labels: app: \"nginx\"\n\t\t\tspec: containers: [{\n\t\t\t\tname:  \"nginx\"\n\t\t\t\timage: \"nginx:1.14.2\"\n\t\t\t\tports: [{\n\t\t\t\t\tcontainerPort: 80\n\t\t\t\t}]\n\t\t\t}]\n\t\t}\n\t}\n}"
}

```


### 返回值

```json
{
  "code": 0,
  "result": "执行CUE命令的结果, 错误时为报错信息, 正确时为空"
}
```


<a name="资源详情"></a>
## 资源详情
### 请求语法
```
GET /resource/<instanceId>/<resourceType>/<resourceName> HTTP/1.1 
```
### 请求参数

instanceId为实例ID

resourceType的值可能为
```
cronJob
daemonSet
deployment
job
pod
replicaSet
replicationController
statefulSet
service
configMap
pvc
secret
roleBinding
role
serviceAccount
```

resourceName的值为[实例详情]接口返回的result.details.<resourceType>.metadata.name的值



### 返回值

在详情中, cronJob展示的字段有:
名称, 创建时间, 最后执行时间
metadata.name, metadata.creationTimestamp, status.lastScheduleTime

在详情中, daemonSet展示的字段有:
名称, 创建时间, 当前可用数
metadata.name, metadata.creationTimestamp, status.numberAvailable

在详情中, deployment展示的字段有:
名称, 创建时间, 当前可用数
metadata.name, metadata.creationTimestamp, status.availableReplicas

在详情中, job展示的字段有:
名称, 创建时间, 完成时间
metadata.name, metadata.creationTimestamp, status.completionTime

在详情中, pod展示的字段有:
名称, 创建时间, 状态, 启动时间
metadata.name, metadata.creationTimestamp, status.phase, status.startTime

在详情中, replicaSet展示的字段有:
名称, 创建时间, 当前可用数
metadata.name, metadata.creationTimestamp, status.availableReplicas

在详情中, replicationController展示的字段有:
名称, 创建时间, 当前副本数
metadata.name, metadata.creationTimestamp, status.replicas

在详情中, statefulSet展示的字段有:
名称, 创建时间, 当前副本数
metadata.name, metadata.creationTimestamp, status.currentReplicas

在详情中, service展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

在详情中, configMap展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

在详情中, pvc展示的字段有:
名称, 创建时间, 状态
metadata.name, metadata.creationTimestamp, status.phase

在详情中, secret展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

在详情中, roleBinding展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

在详情中, role展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp

在详情中, serviceAccount展示的字段有:
名称, 创建时间
metadata.name, metadata.creationTimestamp



每种resourceType的返回值请看 [实例详情] 接口, 下面样例是service的返回
```json
{
  "code": 0,
  "result": {
    "metadata": {
      "name": "crab",
      "namespace": "island-system",
      "uid": "d0e95583-15ca-4198-8dfa-68836f509c19",
      "resourceVersion": "4684705",
      "creationTimestamp": "2021-12-16T10:12:19Z",
      "annotations": {
        "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"name\":\"crab\",\"namespace\":\"island-system\"},\"spec\":{\"ports\":[{\"port\":80,\"targetPort\":3000}],\"selector\":{\"app\":\"island\",\"component\":\"ui\"},\"type\":\"NodePort\"}}\n"
      }
    },
    "spec": {
      "ports": [{
        "protocol": "TCP",
        "port": 80,
        "targetPort": 3000,
        "nodePort": 31997
      }],
      "selector": {
        "app": "island",
        "component": "ui"
      },
      "clusterIP": "10.108.64.173",
      "clusterIPs": ["10.108.64.173"],
      "type": "NodePort",
      "sessionAffinity": "None",
      "externalTrafficPolicy": "Cluster"
    },
    "status": {
      "loadBalancer": {}
    }
  }
}
```


<a name="Trait详情"></a>
## Trait详情
### 请求语法
```
GET /trait/<IdorName> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|IdorName| 主键或名称(名称唯一) |无|是|

### 返回值
```json
{
    "code": 0,
    "result": {
        "id": 1,
        "name": "ingress",
        "apiVersion": "aam.globalsphare.com/v1alpha1",
        "value": "\napiVersion: aam.globalsphare.com/v1alpha1\nkind: Trait\nmetadata:\n    name: ingssssress\nspec:\n    parameter: |\n",
        "type": 1,
        "created_at": "2021-12-16T15:48:48.129+08:00",
        "updated_at": "2021-12-16T16:06:19.228+08:00"
    }
}
```


<a name="WorkloadType详情"></a>
## WorkloadType详情
### 请求语法
```
GET /workloadType/<IdorName> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|IdorName| 主键或名称(名称唯一) |无|是|

### 返回值
```json
{
    "code": 0,
    "result": {
        "id": 1,
        "name": "ingress",
        "apiVersion": "aam.globalsphare.com/v1alpha1",
        "value": "\napiVersion: aam.globalsphare.com/v1alpha1\nkind: Trait\nmetadata:\n    name: ingssssress\nspec:\n    parameter: |\n",
        "type": 1,
        "created_at": "2021-12-16T15:48:48.129+08:00",
        "updated_at": "2021-12-16T16:06:19.228+08:00"
    }
}
```


<a name="WorkloadVendor详情"></a>
## WorkloadVendor详情
### 请求语法
```
GET /workloadVendor/<IdorName> HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|IdorName| 主键或名称(名称唯一) |无|是|

### 返回值
```json
{
    "code": 0,
    "result": {
      "id": 1,
      "name": "ingress",
      "apiVersion": "aam.globalsphare.com/v1alpha1",
      "yaml": "",
      "cue": "",
      "value": "\napiVersion: aam.globalsphare.com/v1alpha1\nkind: Trait\nmetadata:\n    name: ingssssress\nspec:\n    parameter: |\n",
      "type": 1,
      "created_at": "2021-12-16T15:48:48.129+08:00",
      "updated_at": "2021-12-16T16:06:19.228+08:00"
    }
}
```