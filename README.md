# virtual-service-go

An operator to deal with the priority of istio virtual service by watching the virtual service config `crd`.

## Prerequest

- Kubernetes 1.20+ is `required`
 
- Istio 1.12+ is `required`, and only tested in 1.12.1, should be ok in higher version

- Docker is `required`

- Golang 1.17 is `optional` if you want to develop your own features

- Make is `optional` if you want to generate the crd or installation yamls
