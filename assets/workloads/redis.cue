import "mod/context"

import "mod/auth"

parameter: {
  namespace: string
  after?: string
}

outputs: "\(context.appName)-configmap": {
  apiVersion: "v1"
  kind:       "ConfigMap"
  metadata: {
    name:      "\(context.appName)-redis-conf"
    namespace: parameter.namespace
  }
  data: {
    master: """
    pidfile /var/run/redis.pid
    port 6379
    bind 0.0.0.0
    timeout 3600
    tcp-keepalive 1
    loglevel verbose
    logfile /data/redis.log
    slowlog-log-slower-than 10000
    slowlog-max-len 128
    databases 16
    protected-mode no
    save \"\"
    appendonly no

    """

    slave: """
    pidfile /var/run/redis.pid
    port 6379
    bind 0.0.0.0
    timeout 3600
    tcp-keepalive 1
    loglevel verbose
    logfile /data/redis.log
    slowlog-log-slower-than 10000
    slowlog-max-len 128
    databases 16
    protected-mode no
    save \"\"
    appendonly no
    slaveof \(context.componentName)-master 6379

    """
  }
}

outputs: "\(context.componentName)-service-master": {
  apiVersion: "v1"
  kind:       "Service"
  metadata: {
    name:      "\(context.componentName)-master"
    namespace: parameter.namespace
    labels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
    }
  }
  spec: {
    ports: [{
      name: "\(context.componentName)"
      port: 6379
    }]
    selector: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
    }
  }
}

outputs: "\(context.componentName)-service-master-headless": {
  apiVersion: "v1"
  kind:       "Service"
  metadata: {
    name:      "\(context.componentName)-master-headless"
    namespace: parameter.namespace
    labels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
    }
  }
  spec: {
    clusterIP: "None"
    ports: [{
      name: "\(context.componentName)"
      port: 6379
    }]
    selector: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
    }
  }
}

outputs: "\(context.componentName)-service": {
  apiVersion: "v1"
  kind:       "Service"
  metadata: {
    name:      context.componentName
    namespace: parameter.namespace
    labels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
    }
  }
  spec: {
    ports: [{
      name: "\(context.componentName)"
      port: 6379
    }]
    selector: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
    }
  }
}

outputs: "\(context.componentName)-statefulset-master": {
  apiVersion: "apps/v1"
  kind:       "StatefulSet"
  metadata: {
    name:      "\(context.componentName)-master"
    namespace: parameter.namespace
  }
  spec: {
    selector: matchLabels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
    }
    serviceName: "\(context.componentName)"
    replicas:    1
    template: {
      metadata: labels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-master"
      }
      spec: {
        serviceAccountName: context.appName
        if parameter["after"] != _|_ {
          initContainers: [
            {
              name: "init"
              image: "harbor1.zlibs.com/island/centos:7"
              command: ["/bin/sh"]
              args: ["-c", "while true; do curl 'http://island-status.island-system/status/?id=\(parameter.namespace)&component=\(parameter.after)' | grep '\"result\":1'; if [ $? -ne 0 ]; then sleep 4s; continue; else break; fi; done"]
            }
          ]
        }
        containers: [{
          name:  "\(context.componentName)-master"
          image: "harbor1.zlibs.com/dockerhub/redis:6.2.4"
          ports: [{
            containerPort: 6379
            name:          "redis"
          }]
          command: [
            "redis-server",
            "/etc/redis/redis.conf",
          ]
          volumeMounts: [{
            name:      "redis-conf"
            mountPath: "/etc/redis/redis.conf"
            subPath:   "master"
          }]
        }, {
          name:  "\(context.componentName)-sidecar"
          image: "harbor1.zlibs.com/island/centos:7"
          command: ["/bin/sh", "-c", "while true; do curl -X POST http://island-status.island-system/status/ -H 'Content-Type: application/json' -d '{\"id\": \"\(parameter.namespace)\",\"component\": \"\(context.componentName)\"}'; sleep 30s; done;"]
        }]
        volumes: [{
          name: "redis-conf"
          configMap: name: "\(context.appName)-redis-conf"
        }]
      }
    }
  }
}

outputs: "\(context.componentName)-slave-service": {
  apiVersion: "v1"
  kind:       "Service"
  metadata: {
    name:      "\(context.componentName)-slave"
    namespace: parameter.namespace
    labels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-slave"
    }
  }
  spec: {
    ports: [{
      name: "\(context.componentName)"
      port: 6379
    }]
    selector: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-slave"
    }
  }
}

outputs: "\(context.componentName)-slave-service-headless": {
  apiVersion: "v1"
  kind:       "Service"
  metadata: {
    name:      "\(context.componentName)-slave-headless"
    namespace: parameter.namespace
    labels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-slave"
    }
  }
  spec: {
    clusterIP: "None"
    ports: [{
      name: "\(context.componentName)"
      port: 6379
    }]
    selector: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-slave"
    }
  }
}

outputs: "\(context.componentName)-slave-statefulset": {
  apiVersion: "apps/v1"
  kind:       "StatefulSet"
  metadata: {
    name:      "\(context.componentName)-slave"
    namespace: parameter.namespace
  }
  spec: {
    selector: matchLabels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-slave"
    }
    serviceName: "\(context.componentName)"
    replicas:    2
    template: {
      metadata: labels: {
      "app": "\(context.appName)"
      "component": "\(context.componentName)"
      "item": "\(context.componentName)-slave"
      }
      spec: {
        serviceAccountName: context.appName
        containers: [{
          name:  "\(context.componentName)-slave"
          image: "harbor1.zlibs.com/dockerhub/redis:6.2.4"
          ports: [{
            containerPort: 6379
            name:          "redis"
          }]
          command: [
            "bash",
            "-c",
            """
            until [ \"$(echo 'set check_status 1'|timeout 3 redis-cli -h \(context.componentName)-master)\" = \"OK\" ];do sleep 4s;echo \"waiting for the master ready\";done
            redis-server /etc/redis/redis.conf
            """]
          volumeMounts: [{
            name:      "redis-conf"
            mountPath: "/etc/redis/redis.conf"
            subPath:   "slave"
          }]
        }]
        volumes: [{
          name: "redis-conf"
          configMap: name: "\(context.appName)-redis-conf"
        }]
      }
    }
  }
}