
# Overview

- [说明](#说明)
- [添加实例](#添加实例)
- [实例列表](#实例列表)
- [修改管理员密码](#修改管理员密码)
- [获取集群根域](#获取集群根域)
- [设置集群根域](#设置集群根域)
- [获取节点地址](#获取节点地址)
- [获取存储状态](#获取存储状态)
- [初始化存储](#初始化存储)

<a name="说明"></a>
## 说明
目前的接口输出有两种风格

第一种为`{"error":""}`类型
第二种为`{"code":1, "result":""}`类型


### Changelog
#### 20201018
1.因为要兼容原来接口,目前的接口输出实际为两种类型, 待修改

#### 20201023
1.在[添加实例]接口中,出错时添加字段
因为会有未设置的流程, 比如设置根域, 设置密码等操作, 但是这些值都是有默认值的
所以, 要在进行业务时, 前端要提示用户先去完成根域设置, 密码设置等操作

目前是在[添加实例]接口中, 出错时添加 `todo` 字段, 值的含义如下

|  todo   | 意义  |
|  ----  | ----  |
| 1  | 需要设置根域, 提示信息可用 error 字段的内容  |
| 2  | 需要设置存储插件, 提示信息可用 error 字段的内容  |

比如在出错时, 可能的返回结果有
```json
{"error":  "参数错误"}
```

```json
{"error":  "未设置根域, 跳转到设置页面", "todo":  1}
```

目前仅在[添加实例]接口业务出错中会有此字段(此字段指的是`todo`), 后续不排除会在其他接口出现
前端可在接口返回json时, 统一处理, 此字段是可选的数字类型的字段, 可用 `hasOwnProperty` 判断

2. 在[实例列表]接口的返回值中增加 `dependencies` 和 `userconfig` 字段
因为在[添加实例]接口中, 有可能需要去设置存储,设置根域等操作, 此时刚上传的实例并未初始化
增加的两个字段是为了前端能增加初始化的步骤

3.在 [设置集群根域] 接口, 返回值字段中注意 `message` 字段, 和之前`result`字段不同

4.在 [获取存储状态] 接口, 返回值字段中注意 `message` 字段, 和之前`result`字段不同

<a name="添加实例"></a>
## 添加实例

### 请求语法
```
POST /api/app HTTP/1.1
Content-Type: multipart/form-data; 
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|file|实例描述文件（zip 包）|无|是|

### 返回值
缺少插件时的返回值:
```
{"error":  "未设置存储插件, 跳转到插件页面", "todo":  2}
```

未设置根域时的返回值:
```
{"error":  "未设置根域, 跳转到设置页面", "todo":  1}
```

正确返回时:
按照之前的接口返回格式


<a name="实例列表"></a>
## 实例列表

### 请求语法
```
GET /api/app HTTP/1.1
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|offset| |0|否|
|limit|  |10|否|
### 返回值
按照之前的接口返回格式的同时, 增加 dependencies 和 userconfig 字段

```
{
    "rows": [
        {
            "id": "ailwgkhn",
            "name": "demo-app1",
            "entry": "island.com",
            "status": "",
            "version": "0.0.2",
            "dependencies":  {
                    "demo-app1": { #应用名
                        "instances": [
                                {
                                    "name": "demo-app1",
                                    "instanceid": "iqtxycne"
                                },
                                {
                                    "name": "demo-app1",
                                    "location": "https://www.huanqiu.com"
                                }
                        ],
                        "type": "mutable"
                    }
            },
            userconfig: {
                    "properties": {
                    "param1": {
                        "type": "string"
                    },
                    "param2": {
                        "type": "number"
                    },
                    "param3": {
                        "type": "object",
                        "properties": {
                            "param3_1":{
                                "type": "string"
                            },
                            "param3_2":{
                                "type": "number"
                            }
                        }
                    }
                },
                "required": [
                    "param1"
                ],
                "type": "object"
                },
        },
        {
            "id": "aminkcgo",
            "name": "demo-app1",
            "entry": "island2.com",
            "status": "",
            "version": "0.0.2"
        }
    ],
    "total": 67
}
```

`rows[*].entry` 为实例链接


<a name="修改root用户密码"></a>
## 修改管理员密码

### 请求语法
```
PUT /api/user/root HTTP/1.1
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
```
{
    "code": 0,
    "result": "设置成功"
}
```




<a name="设置集群根域"></a>

<a name="获取集群根域"></a>
## 获取集群根域

### 请求语法
```
GET /cluster/domain HTTP/1.1
```
### 请求参数
无

### 返回值
```
{
    "code": 0,
    "result": "abc.com"
}
```


## 设置集群根域

### 请求语法
```
PUT /cluster/domain HTTP/1.1
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|value|绑定到此集群的域名|无|是|
#### 请求参数示例
```
{
    "value": "abc.com",
}
```

### 返回值
```
{
    "code": 0,
    "result": {
      "status": 3,
      "message": "成功",
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


<a name="获取节点IP地址"></a>
## 获取节点地址
显示所有节点的IP地址(仅能显示此节点网卡上绑定的IP地址)

### 请求语法
```
GET /cluster/addrs HTTP/1.1
```
### 请求参数
无

### 返回值
```
{
    "code": 0,
    "result": [
        {
            "name": "master1",
            "addrs": [
              "192.168.0.1",
              "192.168.0.2",
              "192.168.0.3",
            ]
        },
        {
            "name": "salve1",
            "addrs": [
              "192.168.0.4",
            ]
        },
        {
            "name": "salve2",
            "addrs": [
              "192.168.0.5",
              "192.168.0.6",
            ]
        }
    ]
}
```


<a name="获取存储状态"></a>
## 获取存储状态
能获取到存储的状态

### 请求语法
```
GET /cluster/storage HTTP/1.1
```
### 请求参数
无
### 返回值
```
{
    "code": 0,
    "result": {
        "store": {
            "status": 0,
            "message": "", 
        },
        "volumes": [
        {
            "name": "vdc",
            "size": "1024G",
            "hostname": "192.168.0.2",
            "hasChildren": false,
        },
        ]
    }
}
```

`result.store` 为存储集群的状态信息
`result.volumes` 为当前集群内所有的磁盘信息

`name`: 磁盘名称
`size`: 磁盘大小
`hostname`: 所在主机
`hasChildren`: 是否存在子分区

`result.store.status` 为状态码, 数字类型, 具体意义见下表
`result.store.message` 为当前的状态信息, 字符串,  可直接显示在网页

|  status   | 意义  |
|  ----  | ----  |
| 0  | 存储集群还未初始化 |
| 1  | 存储集群正在构建中 |
| 2  | 存储集群平稳运行中 |

<a name="初始化存储"></a>
## 初始化存储

初始化之后, 要再次 `GET /cluster/storage` 判断 `result.store.status` 为 2 时, 存储初始化成功

### 请求语法
```
POST /cluster/storage HTTP/1.1
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|default|是否使用默认配置, true时, 忽略其他参数| false |是|
|volumes|要使用的磁盘数组|无|否|
#### 请求参数示例
```
{
    "default": false,
    "volumes": [
     {
            "name": "vdc",
            "hostname": "192.168.0.2"
     },
     ]
}
```
### 返回值
```
{
    "code": 0,
    "result": "初始化存储成功, 等待内部操作完成"
}
```