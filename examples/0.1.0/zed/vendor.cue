
if parameter.dependencies != _|_ {
	construct: {
        apiVersion: "v1"
        kind:  "ConfigMap"
        metadata: {
          name:      context.appName
          namespace: context.namespace
        }
        data: {
            for k, v in parameter.dependencies {
                "\(k)": v.host
            }
            "runtime": parameter.userconfigs
        }
    }
}


"deployment-zed": {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name:      context.appName
		namespace: context.namespace
		labels: workload: context.workloadName
	}
	spec: {
		replicas: 1
		template: {
			metadata: {
				name: context.appName
				labels: workload: context.workloadName
			}
			spec: {
				containers: [{
					name:            "main"
					image:           "harbor1.zlibs.com/tars/zed:alpha"
					imagePullPolicy: "Always"
					volumeMounts: [{
						mountPath: "/opt"
						name:      "conf"
					}]
				}]
				restartPolicy: "Always"
				serviceAccountName: context.appName
				volumes: [{
					name: "conf"
					configMap: name: context.appName
				}]
			}
		}
		selector: matchLabels: workload: context.workloadName
	}
}
"service-zed": {
	apiVersion: "v1"
	kind:       "Service"
	metadata: {
		name:      context.workloadName
		namespace: context.namespace
	}
	spec: {
		selector: workload: context.workloadName
		ports: [{
			port: 80
			name: "http"
		}]
		type: "ClusterIP"
	}
}