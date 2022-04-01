# Introduction

The virtual service config with more headers will have a higher priority than the others when the uri length are the same.

- The `header-dr.yaml` display how to generate a virtual service to handle 1 service with 3 deployments.

  Use `kubectl apply -f header-dr.yaml` to run the example.

  Use `kubectl get vs header-dr -o yaml` to see the generated virtual service, eg:
  
  ``` yaml
  apiVersion: networking.istio.io/v1beta1
  kind: VirtualService
  metadata:
    creationTimestamp: "2022-03-09T02:05:38Z"
    generation: 10
    name: header-dr
    namespace: default
    resourceVersion: "117788637"
    uid: 43780727-6edc-414a-aefb-a76a2b5432e9
  spec:
    hosts:
    - header-dr
    http:
    - match:
      - headers:
          version:
            exact: b
        uri:
          prefix: /
      name: header-b
      route:
      - destination:
          host: header-dr
          subset: b
    - match:
      - headers:
          version:
            exact: a
        uri:
          prefix: /
      name: header-a
      route:
      - destination:
          host: header-dr
          subset: a
    - match:
      - uri:
          prefix: /
      name: header-dr
      route:
      - destination:
          host: header-dr
          subset: base
  ```

  Use `kubectl delete -f header-dr.yaml` to cleanup the example.

- The `header-svc.yaml` display how to generate a virtual service to handle 3 service with 3 deployments.

  Use `kubectl apply -f header-svc.yaml` to run the example.

  Use `kubectl get vs header-svc -o yaml` to see the generated virtual service, eg:
  
  ``` yaml
  apiVersion: networking.istio.io/v1beta1
  kind: VirtualService
  metadata:
    creationTimestamp: "2022-03-08T13:42:38Z"
    generation: 15
    name: header-svc
    namespace: default
    resourceVersion: "117826597"
    uid: 207fe8c0-26d3-4d95-a54f-1ded13b51ebe
  spec:
    hosts:
    - header-svc
    http:
    - match:
      - headers:
          version:
            exact: b
        uri:
          prefix: /
      name: header-b
      route:
      - destination:
          host: header-svc-b
          subset: latest
    - match:
      - headers:
          version:
            exact: a
        uri:
          prefix: /
      name: header-a
      route:
      - destination:
          host: header-svc-a
          subset: latest
    - match:
      - uri:
          prefix: /
      name: header-svc
      route:
      - destination:
          host: header-svc
          subset: latest  
  ```
  
  Use `kubectl delete -f header-svc.yaml` to cleanup the example.  
