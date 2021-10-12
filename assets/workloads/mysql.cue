import "mod/context"

import "mod/auth"

parameter: {
	namespace: string
	rootpwd: string | "123456"
	storage: {
		capacity: string | "1Gi"
	}
	init: string
	after?: string
}

outputs: "\(context.componentName)-master-configmap": {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: {
		name:      "\(context.componentName)-master"
		namespace: parameter.namespace
		labels: {
			"app": "\(context.appName)"
			"component": "\(context.componentName)"
			"item": "\(context.componentName)-master"
		}
	}
	data: {
		"my.cnf": """
			[mysqld]
			log-bin = mysql-bin
			server-id = 100
			binlog_format=row
			gtid_mode=on
			enforce_gtid_consistency=on
			"""
		"init.sql": """
			\(parameter.init)
			"""
	}
}

outputs: "\(context.componentName)-master-service-headless": {
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
			port: 3306
		}]
		selector: {
			"app": "\(context.appName)"
			"component": "\(context.componentName)"
			"item": "\(context.componentName)-master"
		}
	}
}


outputs: "\(context.componentName)-master-service": {
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
			port: 3306
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
			port: 3306
		}]
		selector: {
			"app": "\(context.appName)"
			"component": "\(context.componentName)"
			"item": "\(context.componentName)-master"
		}
	}
}


outputs: "\(context.componentName)-master-statefulset": {
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
		serviceName: "\(context.componentName)-master-headless"
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
              args: ["-c", "while true; do curl 'http://island-status.island-system/status/?name=\(parameter.namespace)&component=\(parameter.after)' | grep '\"result\":1'; if [ $? -ne 0 ]; then sleep 4s; continue; else break; fi; done"]
						}
					]
				}
				containers: [{
					name:  "\(context.componentName)-master"
					image: "harbor1.zlibs.com/dockerhub/mysql:5.7"
					env: [{
						name:  "MYSQL_ROOT_PASSWORD"
						value: parameter.rootpwd
					}]
					ports: [{
						containerPort: 3306
						name:          "mysql"
					}]
					volumeMounts: [{
						name:      "\(context.componentName)-master"
						mountPath: "/var/lib/mysql"
					}, {
						name:      "conf"
						mountPath: "/etc/mysql/conf.d/mysql.cnf"
						subPath:   "my.cnf"
					}, {
						name:      "conf"
						mountPath: "/docker-entrypoint-initdb.d/init.sql"
						subPath:   "init.sql"
					}]
					command: [
						"bash",
						"-c",
						"""
              rm -rf /var/lib/mysql/lost+found
              echo "start server!"
              /usr/local/bin/docker-entrypoint.sh mysqld
            """,
					]
				},{
					name:  "\(context.componentName)-sidecar"
					image: "harbor1.zlibs.com/island/centos:7"
					command: ["/bin/sh", "-c", "while true; do curl -X POST http://island-status.island-system/status/ -H 'Content-Type: application/json' -d '{\"name\": \"\(parameter.namespace)\",\"component\": \"\(context.componentName)\"}'; sleep 30s; done;"]
				}]
				volumes: [{
					name: "conf"
					configMap: name: "\(context.componentName)-master"
				}]
			}
		}
		volumeClaimTemplates: [{
			metadata: name: "\(context.componentName)-master"
			spec: {
				accessModes: ["ReadWriteOnce"]
				storageClassName: "rook-ceph-block"
				resources: requests: storage: parameter.storage.capacity
			}
		}]
	}
}
outputs: "\(context.componentName)-slave-configmap": {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: {
		name:      "\(context.componentName)-slave"
		namespace: parameter.namespace
		labels: {
			"app": "\(context.appName)"
			"component": "\(context.componentName)"
			"item": "\(context.componentName)-slave"
		}
	}
	data: {
		"my.cnf": """
      [mysqld]
      log-bin = mysql-bin
      binlog_format=row
      gtid_mode=on
      enforce_gtid_consistency=on
    """
		"init.sql": """
      change master to master_host='\(context.componentName)-master-0.\(context.componentName)-master-headless', master_port=3306, master_user='root', master_password='\(parameter.rootpwd)', master_auto_position=1;
      start slave;
    """
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
			port: 3306
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
			port: 3306
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
		serviceName: "\(context.componentName)-slave"
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
					image: "harbor1.zlibs.com/dockerhub/mysql:5.7"
					env: [{
						name:  "MYSQL_ROOT_PASSWORD"
						value: parameter.rootpwd
					}]
					ports: [{
						containerPort: 3306
						name:          "mysql"
					}]
					volumeMounts: [{
						name:      "\(context.componentName)-slave"
						mountPath: "/var/lib/mysql"
					}, {
						name:      "conf"
						mountPath: "/etc/mysql/conf.d/mysql.cnf"
						subPath:   "my.cnf"
					}, {
						name:      "conf"
						mountPath: "/docker-entrypoint-initdb.d/init.sql"
						subPath:   "init.sql"
					}]
					command: [
						"bash",
						"-c",
						"""
            rm -rf /var/lib/mysql/lost+found
            until mysql -h \(context.componentName)-master-0.\(context.componentName)-master-headless -P 3306 -p\(parameter.rootpwd) -e \"SELECT 1\"; do sleep 1; done
            [[ `hostname` =~ -([0-9]+)$ ]] || exit 1
            ordinal=${BASH_REMATCH[1]}
            echo [mysqld] > /etc/mysql/conf.d/server-id.cnf
            echo server-id=$((101 + $ordinal)) >> /etc/mysql/conf.d/server-id.cnf
            echo "run mysql!!"
            /usr/local/bin/docker-entrypoint.sh mysqld
            """]
				}]
				volumes: [{
					name: "conf"
					configMap: name: "\(context.componentName)-slave"
				}]
			}
		}
		volumeClaimTemplates: [{
			metadata: name: "\(context.componentName)-slave"
			spec: {
				accessModes: ["ReadWriteOnce"]
				storageClassName: "rook-ceph-block"
				resources: requests: storage: parameter.storage.capacity
			}
		}]
	}
}