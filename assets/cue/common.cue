context: {
	appName:       string
	componentName: string
	namespace:     string
}
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
if serviceEntry != _|_ {
	for k, v in serviceEntry {
		"construct": "serviceEntry-\(context.componentName)-to-\(v.name)": {
			apiVersion: "networking.istio.io/v1alpha3"
			kind:       "ServiceEntry"
			metadata: {
				name:      "\(context.componentName)-to-\(v.name)"
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

if ingress != _|_ {
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
						ingress.host,
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
						ingress.host,
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
					name: context.componentName
					if ingress.http != _|_ {
						match: []
					}
					route: [{
						destination: {
							port: number: 80
							host: context.componentName
						}
						headers: {
							request: {
								add: {
									"X-Forwarded-Host": ingress.host
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
								host: context.componentName
								port: {
									number: 80
								}
							}
							headers: {
								request: {
									add: {
										"X-Forwarded-Host": ingress.host
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
