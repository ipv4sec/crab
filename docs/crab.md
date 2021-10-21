
# Overview

- [说明](#说明)
- [插件列表](#插件列表)

<a name="说明"></a>
## 说明
目前的接口输出有两种风格

第一种为`{"error":""}`类型
第二种为`{"code":1, "result":""}`类型


### Changelog
#### 20201018
1.因为要兼容原来接口,目前的接口输出实际为两种类型, 待修改

#### 20201020
1.在[添加实例]接口中,出错时添加字段
因为会有未设置的流程, 比如设置根域, 设置密码等操作, 但是这些值都是有默认值的
所以, 要在进行业务时, 前端要提示用户先去完成根域设置, 密码设置等操作

目前是在[添加实例]接口中, 出错时添加 `todo` 字段, 值的含义如下

|  todo   | 意义  |
|  ----  | ----  |
| 0  | 没有需要设置的事项, 可直接显示 error 字段的内容 |
| 1  | 需要设置根域, 提示信息可用 error 字段的内容  |
| 2  | 需要设置存储插件, 提示信息可用 error 字段的内容  |

比如在出错时, 可能的返回结果有
```json
{"error":  "参数错误"}
```

```json
{"error":  "未设置根域, 跳转到设置页面", "todo":  1}
```

目前仅在[添加实例]接口中会有此字段(此字段指的是`todo`), 后续不排除会在其他接口出现
前端可在接口返回json时, 统一处理, 此字段是可选数字类型的字段, 可用 `hasOwnProperty` 判断


<a name="插件列表"></a>
## 插件列表
### 请求语法
```
GET /api/plugins HTTP/1.1
Content-Type: application/json
```
### 请求参数

无

### 返回值

```
{
    "code": 0,
    "result": [
        {
            "id": 1,
            "name": "存储插件",
            "status": 0,
        },
        ]
}
```

|  status   | 意义  |
|  ----  | ----  |
| 0  | 未安装 |
| 1  | 正在安装中  |
| 2  | 安装成功 |
| -1  | 安装失败 |


<a name="修改插件"></a>
## 修改插件
目前只有插件的状态能修改, 通过修改插件的状态来安装插件

### 请求语法
```
PUT /api/plugin/<id> HTTP/1.1
Content-Type: application/json
```
### 请求参数
|名称|说明|默认值|是否必填|
|---|---|---|---|
|status|要修改为的状态, 数字类型|无|是|
#### 请求参数示例
```
{
    "status": 1,
    "param": {
    
    }
}
```
### 返回值
```
{
    "code": 0,
    "result": "正在安装中"
}
```

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
```
```


<a name="实例列表"></a>
## 实例列表

|  status   | 意义  |
|  ----  | ----  |
| 0  | 未部署 |
| 1  | 正在部署中  |
| 2  | 部署完成  |
| 3  | 卸载中  |
| 4  | 卸载完成  |

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
                        ], #应用实例 id 列表
                        "userconfig": {
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
                        }
                        "location": "https://www.huanqiu.com",
                        "resources": [
                            {
                                "actions": [
                                    "GET"
                                ],
                                "uri": "/app"
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

`rows.entry` 为实例链接


<a name="修改root用户密码"></a>
## 修改root用户密码

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
## 设置集群根域

### 请求语法
```
PUT /api/domain HTTP/1.1
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
