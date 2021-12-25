# Overview

- [说明](#说明)
- [错误代码](#错误代码)
- [解析manifest文件](#解析manifest文件)
- [获取系统cue模板](#获取系统cue模板)

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
|Host|实例访问域名|string|无|是|

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
    "result": "k8s yaml"
}
```

字段result的内容

```yaml
apiVersion: apps/v1
kind: Deployment
//...
---
apiVersion: v1
data:
  userconfigs: "{}"
kind: ConfigMap
//...
---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
//...
---
apiVersion: apps/v1
kind: Deployment
//...
---
apiVersion: v1
data:
userconfig: "null"
kind: ConfigMap
//...
---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
//...
```

<a name="获取系统cue模板"></a>
## 获取系统cue模板

### 请求语法
```
GET /systemTemplate HTTP/1.1
```

### 请求参数
无

### 返回值
```
{
    "code": 0,
    "result": "cue template"
}
```

字段result的内容

```
context: {
	appName:       string
	workloadName:  string
	namespace:     string
}
parameter:{
    authorization?: [...{
    	service:   string
    	namespace: string
    	resources?: [...{
    		uri: string
    		action: [...string]
    	}]
    }]
    serviceEntry?: [...{
    	name:     string
    	host:     string
    	address:  string
    	port:     int
    	protocol: string
    }]
    dependencies: [...{[string]: host: string}]
    userconfigs: string
    ingress?: {
    	host: string
    	path?: [...string]
    }
}

construct: namespace: {
	apiVersion: "v1"
	kind:       "Namespace"
	metadata: {
		name: context.namespace
		labels: {
			"istio-injection": "enabled"
		}
	}
}
construct: serviceAccount: {
	apiVersion: "v1"
	kind:       "ServiceAccount"
	metadata: {
		name:      context.appName
		namespace: context.namespace
	}
}
construct: "default-authorizationPolicy": {
	apiVersion: "security.istio.io/v1beta1"
	kind:       "AuthorizationPolicy"
	metadata: {
		name:      context.namespace
		namespace: context.namespace
	}
	spec: {}
}
if parameter.serviceEntry != _|_ {
	for k, v in parameter.serviceEntry {
		"construct": "serviceEntry-\(context.workloadName)-to-\(v.name)": {
			apiVersion: "networking.istio.io/v1alpha3"
			kind:       "ServiceEntry"
			metadata: {
				name:      "\(context.workloadName)-to-\(v.name)"
				namespace: context.namespace
			}
			spec: {
				exportTo: ["."]
				hosts: [
					v.host,
				]
				if v.address != _|_ {
					addresses: [
						v.address,
					]
				}
				location: "MESH_EXTERNAL"
				ports: [
					{
						number:   v.port
						name:     "port-name"
						protocol: v.protocol
					},
				]
			}
		}
	}
}
if authorization != _|_ {
	for k, v in authorization {
		"construct": "island-allow-\(context.namespace)-to-\(v.namespace)-\(v.service)": {
			apiVersion: "security.istio.io/v1beta1"
			kind:       "AuthorizationPolicy"
			metadata: {
				name:      "\(context.namespace)-to-\(v.namespace)-\(v.service)"
				namespace: v.namespace
			}
			spec: {
				action: "ALLOW"
				selector: {
					matchLabels: {
						"workload": v.service
					}
				}
				rules: [
					{
						from: [
							{source: principals: ["cluster.local/ns/\(context.namespace)/sa/\(context.appName)"]},
						]
						if v.resources != _|_ {
							to: [
								for resource in v.resources {
									operation: {
										methods: resource.actions
										paths: [resource.uri]
									}
								},
							]
						}
					},
				]
			}
		}
	}
}

if parameter.ingress != _|_ {
	"ingress": "ingressgateway-http": {
		apiVersion: "networking.istio.io/v1alpha3"
		kind:       "Gateway"
		metadata: {
			name:      "\(context.namespace)-http"
			namespace: "island-system"
		}
		spec: {
			selector: istio: "ingressgateway"
			servers: [
				{
					port: {
						number:   80
						name:     "http"
						protocol: "HTTP"
					}
					hosts: [
						parameter.ingress.host,
					]
				},
			]
		}
	}
	"ingress": "ingressgateway-https": {
		apiVersion: "networking.istio.io/v1alpha3"
		kind:       "Gateway"
		metadata: {
			name:      "\(context.namespace)-https"
			namespace: "island-system"
		}
		spec: {
			selector: istio: "ingressgateway"
			servers: [
				{
					port: {
						number:   443
						name:     "https"
						protocol: "HTTPS"
					}
					tls: {
						mode:              "SIMPLE"
						serverCertificate: "/etc/istio/ingressgateway-certs/tls.crt"
						privateKey:        "/etc/istio/ingressgateway-certs/tls.key"
					}
					hosts: [
						parameter.ingress.host,
					]
				},
			]
		}
	}
	"ingress": "virtualservice-http": {
		apiVersion: "networking.istio.io/v1alpha3"
		kind:       "VirtualService"
		metadata: {
			name:      "\(context.appName)-http"
			namespace: context.namespace
		}
		spec: {
			hosts: ["*"]
			gateways: ["island-system/\(context.namespace)-http"]
			http: [
				{
					name: context.workloadName
					if ingress.http != _|_ {
						match: []
					}
					route: [{
						destination: {
							port: number: 80
							host: context.workloadName
						}
						headers: {
							request: {
								add: {
									"X-Forwarded-Host": parameter.ingress.host
								}
							}
						}
					}]
				},
			]
		}
	}
	"ingress": "virtualservice-https": {
		apiVersion: "networking.istio.io/v1alpha3"
		kind:       "VirtualService"
		metadata: {
			name:      "\(context.appName)-https"
			namespace: context.namespace
		}
		spec: {
			hosts: ["*"]
			gateways: ["island-system/\(context.namespace)-https"]
			http: [
				{
					match: []
					route: [
						{
							destination: {
								host: context.workloadName
								port: {
									number: 80
								}
							}
							headers: {
								request: {
									add: {
										"X-Forwarded-Host": parameter.ingress.host
									}
								}
							}
						},
					]
				},
			]
		}
	}
}
```