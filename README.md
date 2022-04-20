# virtual-service-go

An operator to deal with the priority of istio virtual service by watching the virtual service config `crd`.

## Prerequest

- Kubernetes 1.20+ is `required`
 
- Istio 1.12+ is `required`, and only tested in 1.12.1, should be ok in higher version

- Docker is `required`

- Golang 1.17 is `optional` if you want to develop your own features

- Make is `optional` if you want to generate the crd or installation yamls

## Installation

- Gen or re-gen crd

  Using `make manifests` to gen or re-gen crd.

- Pack docker image
  
  Using `bash pack.sh` to pack the docker image, If you want to push the image to your docker registry, please modify the [pack.sh](pack.sh) and using your registry address before running `bash pack.sh`.

- Install crd and operator

  - Enter the [config](config) directory

  - With your cluster `kubeconfig`, running `kubectl apply -k crd/` to install crd

  - With your cluster `kubeconfig`, running `kubectl apply -k default/` to install operator
  
  If you are using your own docker registry, please modify the image info of [manager.yaml](config/manager/manager.yaml) before running `kubectl apply -k default/`.

## Examples

- Enter the [example](example) directory

- Using `kubectl apply -f xxx.yaml` to easily start an example

- Using `kubectl get vs xxx -o yaml` to see the generated virtual service

- Using `kubectl get vsc` to see the virtual service config(s)

- Using `kubectl delete -f xxx.yaml` to cleanup the example
