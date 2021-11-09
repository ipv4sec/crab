# Overview

### 请求语法
```
POST / HTTP/1.1
```

### 请求参数
|名称|说明|类型|默认值|是否必填|
|---|---|---|---|---|
|deployment|翻译器生成的文件内容|string|无|是|
|id|实例ID|string|无|是|

### 返回值
```
{
    "code": 0,
    "result": "ok"
}
```