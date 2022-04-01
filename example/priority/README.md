# Introduction

The virtual service configs use `order` to indicate the priority, the larger vaule will have a higher priority than the others. 

eg: `order = 2` will have a higher priority than `order = 1`.

Use `kubectl apply -f priority.yaml` to run the example.

Use `kubectl get vs priority -o yaml` to see the generated virtual service, eg:

``` yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  creationTimestamp: "2022-03-09T07:11:53Z"
  generation: 1
  name: priority
  namespace: default
  resourceVersion: "118178042"
  uid: 7483976b-8f33-4a7d-b7c7-c78ea32e95d4
spec:
  hosts:
  - priority
  http:
  - match:
    - uri:
        prefix: /
    name: short-no-header
    route:
    - destination:
        host: priority
        subset: short-no-header
  - match:
    - uri:
        prefix: /long
    name: long-no-header
    route:
    - destination:
        host: priority
        subset: long-no-header
  - match:
    - headers:
        version:
          exact: b
      uri:
        prefix: /long
    name: long-with-header
    route:
    - destination:
        host: priority
        subset: long-with-header
```

Use `kubectl delete -f priority.yaml` to cleanup the example.
