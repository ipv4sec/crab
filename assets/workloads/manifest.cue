parameter: {
  manifest: {
    appName: string
    namespace: string
  }
}

outputs: "island-deny-\(parameter.manifest.namespace)":{
  apiVersion: "security.istio.io/v1beta1"
  kind: "AuthorizationPolicy"
  metadata: {
    name: parameter.manifest.namespace
    namespace: parameter.manifest.namespace
  }
  spec:
    {}
}