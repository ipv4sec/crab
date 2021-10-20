import "mod/context"

import "mod/auth"

import "mod/entry"

import "mod/storage"

import "mod/configmap"

parameter: {
  namespace: string
  image:     string
  port:      *80 | int
  cmd?: [...string]
  args?: [...string]
  cpu?: string
  env?: [...{
    name:   string
    value?: string
    valueFrom?: {
      secretKeyRef: {
        name: string
        key:  string
      }
    }
  }]
  after?: string
}

outputs: "\(context.componentName)-deployment":{
  apiVersion: "apps/v1"
  kind:       "Deployment"
  metadata: {
    name:      context.componentName
    namespace: parameter.namespace
  }
  spec: {
    selector: matchLabels: {
      "app":       context.appName
      "component": context.componentName
    }

    template: {
      metadata: labels: {
        "app":       context.appName
        "component": context.componentName
      }

      spec: {
        serviceAccountName: context.appName
        if parameter["after"] != _|_ {
          initContainers: [
            {
              name:  "init"
              image: "harbor1.zlibs.com/island/centos:7"
              command: ["/bin/sh"]
              args: ["-c", "while true; do curl 'http://island-status.island-system/status/?id=\(parameter.namespace)&component=\(parameter.after)' | grep '\"result\":1'; if [ $? -ne 0 ]; then sleep 4s; continue; else break; fi; done"]
            },
          ]
        }
        containers: [{
          name:  context.componentName
          image: parameter.image

          if parameter["cmd"] != _|_ {
            command: parameter.cmd
          }

          if parameter["args"] != _|_ {
            args: parameter.args
          }

          if parameter["env"] != _|_ {
            env: parameter.env
          }

          if parameter["cpu"] != _|_ {
            resources: {
              limits:
                cpu: parameter.cpu
              requests:
                cpu: parameter.cpu
            }
          }

          ports: [{
            containerPort: parameter.port
          }]

          volumeMounts: [
            for k, v in parameter.configs if v["subPath"] != _|_ {
              name:      "\(context.componentName)-\(k)"
              mountPath: "\(v.path)/\(v.subPath)"
              subPath:   v.subPath
            }
            for k, v in parameter.configs if v["subPath"] == _|_ {
              name:      "\(context.componentName)-\(k)"
              mountPath: v.path
            }
            if context.appName == "island-console" {
              name:      "island-info"
              mountPath: "/etc/island-info"
            },
            if parameter.storage.capacity != _|_ {
              name:      "storage-\(context.componentName)"
              mountPath: parameter.storage.path
            },
          ]
        }, {
          name:  "\(context.componentName)-sidecar"
          image: "harbor1.zlibs.com/island/centos:7"
          command: ["/bin/sh", "-c", "while true; do curl -X POST http://island-status.island-system/status/ -H 'Content-Type: application/json' -d '{\"id\": \"\(parameter.namespace)\",\"component\": \"\(context.componentName)\"}'; sleep 30s; done;"]
        }]

      volumes: [
        for k, v in parameter.configs if v["subPath"] != _|_ {
          name:      "\(context.componentName)-\(k)"
          configMap: name: "\(context.componentName)-\(k)"
        }
        for k, v in parameter.configs if v["subPath"] == _|_ {
          name:      "\(context.componentName)-\(k)"
          configMap: name: "\(context.componentName)-\(k)"
        }
        if context.appName == "island-console" {
          name: "island-info"
          configMap: name: "island-info"
        },
        if parameter.storage.capacity != _|_ {
          name: "storage-\(context.componentName)"
          persistentVolumeClaim: claimName: "storage-\(context.componentName)"
        },
      ]
    }
    }
  }
}

outputs: "serviceinternal": {
  apiVersion: "v1"
  kind:       "Service"
  metadata: {
    name:      context.componentName
    namespace: parameter.namespace
  }
  spec: {
    selector: {
      "app":       context.appName
      "component": context.componentName
    }
    ports: [
      {
        name: "http"
        port: 80
        if parameter.port != _|_ {
          targetPort: parameter.port
        }
        if parameter.port == _|_ {
          targetPort: 80
        }
      },
    ]
  }
}

outputs: "\(context.componentName)-viewer": {
  "apiVersion": "security.istio.io/v1beta1"
  "kind":       "AuthorizationPolicy"
  "metadata": {
    "name":      "\(context.componentName)-viewer"
    "namespace": parameter.namespace
  }
  "spec": {
    "selector": {
      "matchLabels": {
        "app":       context.appName
        "component": context.componentName
      }
    }
    "rules": [
      {
        to: [
          {
            operation: {
              methods: ["GET", "POST", "DELETE", "PUT","HEAD", "OPTIONS","PATCH"]
            }
          },
        ]
      },
    ]
  }
}
