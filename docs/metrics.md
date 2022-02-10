
# Overview

- [获取CPU/内存的使用数据](#获取CPU/内存的使用数据)


<a name="获取CPU/内存的使用数据"></a>
## 获取CPU/内存的使用数据
### 请求语法
```
GET /metrics/:namespace/:resourceName HTTP/1.1
```
### 请求参数
以下参数为URL PATH参数

|名称|说明|默认值|是否必填|
|---|---|---|---|
|namespace| 实例主键 |无|是|
|resourceName| 资源名称 |无|是|


### 返回值

```json
{
  "code": 0,
  "result": {
    "cpu": [
      {
        "x": 111,
        "y": 222
      }
    ],
    "mem": [
      {
        "x": 111,
        "y": 222
      }
    ]
  }
}
```