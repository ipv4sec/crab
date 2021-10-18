import "mod/context"

parameter: {
  storage?: {
    capacity: string
    path:     string
  }
}

if parameter.storage != _|_ {
  if parameter.storage.capacity != ""  {
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

