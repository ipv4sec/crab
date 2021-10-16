import "mod/auth"
import "mod/configmap"
import "mod/storage"

parameter: {
  image: string
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

outputs: {
  statefulsetheadless: {
    apiVersion: "v1"
    kind:       "Service"
    metadata: {
      name:      "\(context.componentName)-headless"
      namespace: parameter.namespace
      labels: {
        "component": context.componentName
        "app": context.appName
      }
    }
    spec: {
      clusterIP: "None"
      selector: {
        "component": context.componentName
        "app": context.appName
      }
    }
  }
}

outputs: "\(context.componentName)-statefulset":{
  apiVersion: "apps/v1"
  kind:       "StatefulSet"
  metadata: {
    name:      context.componentName
    namespace: parameter.namespace
  }
  spec: {
    selector: matchLabels: {
      "app": context.appName
      "component": context.componentName
    }
    replicas:    1
    serviceName: "\(context.componentName)-headless"
    template: {
      metadata: labels: {
        "app": context.appName
        "component": context.componentName
      }

      spec: {
        serviceAccountName: context.appName
        if parameter["after"] != _|_ {
          initContainers: [
            {
              name: "init"
              image: "harbor1.zlibs.com/island/centos:7"
              command: ["/bin/sh"]
              args: ["-c", "while true; do curl 'http://island-status.island-system/status/?name=\(parameter.namespace)&component=\(parameter.after)' | grep '\"result\":1'; if [ $? -ne 0 ]; then sleep 4s; continue; else break; fi; done"]
            }
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

          volumeMounts: [
            for k, v in parameter.configs if v["subPath"] != _|_ {
              name:      "\(context.componentName)-\(k)"
              mountPath: "\(v.path)/\(v.subPath)"
              subPath:   v.subPath
            },
            for k, v in parameter.configs if v["subPath"] == _|_ {
              name:      "\(context.componentName)-\(k)"
              mountPath: v.path
            }
            if context.appName == "island-console" {
              name:      "island-info"
              mountPath: "/etc/island-info"
            },
            if parameter.storage!= _|_ {
              if len(parameter.storage.capacity) > 0 {
                name:      "storage-\(context.componentName)"
                mountPath: parameter.storage.path
              },
            }
          ]
        }, {
          name:  "\(context.componentName)-sidecar"
          image: "harbor1.zlibs.com/island/centos:7"
          command: ["/bin/sh", "-c", "while true; do curl -X POST http://island-status.island-system/status/ -H 'Content-Type: application/json' -d '{\"name\": \"\(parameter.namespace)\",\"component\": \"\(context.componentName)\"}'; sleep 30s; done;"]
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
        if parameter.storage!= _|_ {
          if len(parameter.storage.capacity) > 0 {
            name: "storage-\(context.componentName)"
            persistentVolumeClaim: claimName: "storage-\(context.componentName)"
          }
        }
      ]
    }
    }
  }
}