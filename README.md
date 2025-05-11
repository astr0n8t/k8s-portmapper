# k8s-portmapper

Auto-add new listening ports in a container to a k8s service.

See config.yaml for options.

## k8s Pod Privileges

Make sure to add a service account to be able to get and edit the service:
```
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: service-patcher
rules:
- apiGroups: [""]
  resources: ["services"]
  verbs: ["get", "patch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: default
  name: service-patcher-binding
subjects:
- kind: ServiceAccount
  name: default # Replace with your service account
  namespace: default
roleRef:
  kind: Role
  name: service-patcher
  apiGroup: rbac.authorization.k8s.io
```

And then apply it to your pod, as well as give the pod `shareProcessNamespace: true` if wanting to filter on a specific program.
