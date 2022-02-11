
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

resourceName需要具体的POD名称， 当查询多个POD的使用数据时， resourceName的值以逗号分隔即可

比如当要查询cronjob、daemonsets、deployment、jobs、pods、replicasets、replicationcontrollers、statefulsets的资源
使用情况时，resourceName的值为该资源下的所有POD，请求路径为 `GET /metrics/ins233/island-api-78577f88c8-wpkbg,island-webssh-67d9b969fb-9t6fg`
比如当要查询名为island-api-78577f88c8-wpkbg的具体POD的资源使用情况时，请求路径为
`GET /metrics/ins233/island-api-78577f88c8-wpkbg`

### 返回值

返回JSON中x为坐标轴x轴的值，表示的是时间戳
返回JSON中y为坐标轴y轴的值，result.cpu[].y 表示的是占用的CPU核心数量， result.mem[].y 表示的是占用的内存字节

```json
{
  "code": 0,
  "result": {
    "cpu": [{
      "x": 1644487740,
      "y": 5
    }, {
      "x": 1644487800,
      "y": 6
    }, {
      "x": 1644487860,
      "y": 4
    }, {
      "x": 1644487920,
      "y": 5
    }],
    "mem": [{
      "x": 1644487740,
      "y": 18272256
    }, {
      "x": 1644487800,
      "y": 18477056
    }, {
      "x": 1644487860,
      "y": 18534400
    }, {
      "x": 1644487920,
      "y": 19021824
    }]
  }
}
```

错误返回时， 前端不显示表格
```json
{"code":10102,"result":"错误提示信息"}
```