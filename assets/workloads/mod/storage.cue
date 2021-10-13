import "mod/context"

parameter: {
	storage?: {
		capacity: string | *"1Gi"
		path:     string
	}
}

if parameter.storage != _|_ {
  if parameter.storage.capacity != _|_ {
    outputs: {
      "storage": {
        apiVersion: "v1"
        kind:       "PersistentVolumeClaim"
        metadata: {
          name:      "storage-\(context.componentName)"
          namespace: parameter.namespace
        }
        spec: {
          storageClassName: "rook-ceph-block"
          accessModes: [
            "ReadWriteOnce",
          ]
          resources: requests: storage: parameter.storage.capacity
        }
      }
    }
	}
}
