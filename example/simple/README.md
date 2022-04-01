# Introduction

This example will show how to use a virtual service config to generate a istio virtual service. 

Use `kubectl apply -f simple.yaml` to run the example.

Use `kubectl get vs simple -o yaml` to see the generated virtual service, eg:

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  creationTimestamp: "2022-03-09T01:20:52Z"
  generation: 1
  name: simple
  namespace: default
  resourceVersion: "117725056"
  uid: f41c4715-5cbd-41aa-9877-ab3d25cc9227
spec:
  hosts:
  - simple
  http:
  - match:
    - uri:
        prefix: /
    name: simple
    route:
    - destination:
        host: simple
        subset: latest
```

Use `kubectl delete -f simple.yaml` to cleanup the example.
