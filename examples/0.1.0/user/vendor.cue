"deployment-user": {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name: context.appName
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
					image:           "nginx:1.21"
					imagePullPolicy: "Always"
				}]
				restartPolicy: "Always"
				serviceAccountName: context.appName
			}
		}
		selector: matchLabels: workload: context.workloadName
	}
}
"service-user": {
	apiVersion: "v1"
	kind:       "Service"
	metadata: {
		name: context.workloadName
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