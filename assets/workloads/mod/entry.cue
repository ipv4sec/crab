import "mod/context"

parameter: {
  namespace: string
  entry?: {
    host: string
    path?: [...string]
  }
}

if parameter["entry"]["host"] != "" {
  outputs: "ingressgateway-http": {
    apiVersion: "networking.istio.io/v1alpha3"
    kind: "Gateway"
    metadata: {
      name: "\(parameter.namespace)-http"
      namespace: "island-system"
    }
    spec: {
      selector: istio: "ingressgateway"
      servers: [
        {
          port: {
            number: 80
            name: "http"
            protocol: "HTTP"
          }
          hosts: [
            parameter["entry"]["host"],
          ]
        },
      ]
    }
  }
  outputs: "ingressgateway-https": {
    apiVersion: "networking.istio.io/v1alpha3"
    kind: "Gateway"
    metadata: {
      name: "\(parameter.namespace)-https"
      namespace: "island-system"
    }
    spec: {
      selector: istio: "ingressgateway"
      servers: [
        {
          port: {
            number: 443
            name: "https"
            protocol: "HTTPS"
          }
          tls: {
            mode: "SIMPLE"
            serverCertificate: "/etc/istio/ingressgateway-certs/tls.crt"
            privateKey: "/etc/istio/ingressgateway-certs/tls.key"
          }
          hosts: [
            parameter["entry"]["host"],
          ]
        },
      ]
    }
  }

  outputs: "virtualservice-http": {
    apiVersion: "networking.istio.io/v1alpha3"
    kind: "VirtualService"
    metadata: {
      name: "\(context.appName)-http"
      namespace: parameter.namespace
    }
    spec: {
      hosts: ["*"]
      gateways: ["island-system/\(parameter.namespace)-http"]
      http: [
        {
          name: context.componentName
          if parameter.entry.path != _|_ {
            match: [
              for k, v in parameter.entry.path {
                {uri: regex: v}
              },
            ]
          }
          route: [{
            destination: {
              port: number: 80
              host: context.componentName
            }
          }]
        },
      ]
    }
  }

  outputs: "virtualservice-https": {
    apiVersion: "networking.istio.io/v1alpha3"
    kind: "VirtualService"
    metadata: {
      name: "\(context.appName)-https"
      namespace: parameter.namespace
    }
    spec: {
      hosts: ["*"]
      gateways: ["island-system/\(parameter.namespace)-https"]
      http: [
        {
          match: [
            {
              uri: {
                regex: "/.*"
              }
            }
          ]
          route: [
            {
              destination: {
                host: context.componentName
                port: {
                  number: 80
                }
              }
            }
          ]
        }
      ]
    }
  }
}

