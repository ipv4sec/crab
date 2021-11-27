# Overview

- [说明](#说明)
- [错误代码](#错误代码)
- [解析manifest文件](#解析manifest文件)

<a name="说明"></a>

## 说明

### Changelog

### 返回的数据结构如下
```
{
    "code": 错误代码,
    "result": 任意类型
}
```

`code` 表示的为错误代码, 数字类型

正常返回时`code`为`0`, `result`为数据, 可能是数字, 布尔, 字符串, 数组, 对象等等

错误返回时`code`不为`0`, `result`为错误信息, 字符串类型, 可直接显示在浏览器页面

错误代码原则上前端用不到, 前端仅需要判断非0时显示`result`字段即可


<a name="错误代码"></a>
## 错误代码

|  code     |意义  | 
|  ----   |----  |
| 10201   | 参数错误 |
| 10202   | 服务器内部错误 |

错误代码的第一位, 目前1 标识此应用
错误代码的第二三位, 标识组件
错误代码的第四五位, 标识错误代码

<a name="解析manifest文件"></a>
## 解析manifest文件


### 请求语法
```
POST / HTTP/1.1
```

### 请求参数
|名称|说明|类型|默认值|是否必填|
|---|---|---|---|---|
|Content|文件内容|string|无|是|
|InstanceId|实例id|string|无|是|
|UserConfig|运行时配置|object|{}|否|
|Dependencies|实例依赖|dependency object|{}|否|
|RootDomain|根域|string|无|是|
|WorkloadPath|绝对路径|string|无|是|

Dependencies.Internal 内部的服务, 数组类型, 非必填, 内容如下：

* Dependencies.Internal.[i].Name 应用的name string类型, 必填
* Dependencies.Internal.[i].InstanceId 内部服务实例id，string类型, 必填
* Dependencies.Internal.[i].EntryService 服务暴露的组件的名称，string类型, 必填

Dependencies.External 外部的服务, 数组类型, 非必填, 内容如下：

* Dependencies.External[i].Name 应用的name string类型, 必填
* Dependencies.External[i].Location string类型，必填

### 返回值
```
{
    "code": 0,
    "result": "中间格式yaml"
}
```

中间格式反序列化后yaml字段：

|名称|说明|类型|
|---|---|---|
|name|应用名称|string|
|init|初始化语句|string|
|workloads|工作负载|数组|

workloads的字段：
* workloads.[i].parameter 为工作负载的参数 string类型
* workloads.[i].construct 为工作负载主体部分  map[string]string类型
* workloads.[i].traits 为workload的trait map[string]string类型


```yaml
name: cs
init: |2

  apiVersion: v1
  kind: Namespace
  metadata:
   name: ins1638001697
   labels:
     istio-injection: enabled
workloads:
  ac:
    parameter: |
      after: ac-db
      image: harbor1.zlibs.com/island/ac:0.0.1
      port: 80
    construct: //workoadType的实现
      ac-deployment: |
        apiVersion: apps/v1
        kind: Deployment
        //...
      island-ac-0: |
        apiVersion: v1
        data:
          userconfig: "null"
        kind: ConfigMap
        //...
      island-allow-ins1638001697-to-ins1638001697-ac-db: |
        apiVersion: security.istio.io/v1beta1
        kind: AuthorizationPolicy
        //...
    traits: //特性
      ingressgateway-http: |
        apiVersion: networking.istio.io/v1alpha3
        kind: Gateway
        //...
      ingressgateway-https: |
        apiVersion: networking.istio.io/v1alpha3
        kind: Gateway
        //...
      virtualservice-http: |
        apiVersion: networking.istio.io/v1alpha3
        kind: VirtualService
        //...
      virtualservice-https: |
        apiVersion: networking.istio.io/v1alpha3
        kind: VirtualService
        //... 
```
