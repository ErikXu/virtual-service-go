# Introduction

The virtual service config with longer uri will have a higher priority than the others.

Use `kubectl apply -f uri.yaml` to run the example.

Use `kubectl get vs uri -o yaml` to see the generated virtual service, eg:

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  creationTimestamp: "2022-03-09T01:52:24Z"
  generation: 1
  name: uri
  namespace: default
  resourceVersion: "117765723"
  uid: 57695480-f2e2-49f5-8154-4a53aa2b977f
spec:
  hosts:
  - uri
  http:
  - match:
    - uri:
        prefix: /long
    name: long
    route:
    - destination:
        host: long
        subset: latest
  - match:
    - uri:
        prefix: /
    name: short
    route:
    - destination:
        host: short
        subset: latest
```

Use `kubectl delete -f uri.yaml` to cleanup the example.
