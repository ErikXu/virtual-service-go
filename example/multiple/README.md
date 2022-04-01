# Introduction

This example will show how to use multiple virtual service configs to generate a istio virtual service. 

Use `kubectl apply -f multiple.yaml` to run the example.

Use `kubectl get vs multiple -o yaml` to see the generated virtual service, eg:

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  creationTimestamp: "2022-03-09T01:27:59Z"
  generation: 2
  name: multiple
  namespace: default
  resourceVersion: "117734227"
  uid: 674bdbe3-fdf6-4d4c-9a36-9727dfe419aa
spec:
  hosts:
  - multiple
  http:
  - match:
    - uri:
        prefix: /v1
    name: v1
    route:
    - destination:
        host: multiple
        subset: v1
  - match:
    - uri:
        prefix: /v2
    name: v2
    route:
    - destination:
        host: multiple
        subset: v2
```

Use `kubectl delete -f multiple.yaml` to cleanup the example.
