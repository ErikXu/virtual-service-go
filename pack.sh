#!/bin/bash

IMAGE_NAME=${IMAGE_NAME:-vs-conf-operator}
echo "IMAGE_NAME: "$IMAGE_NAME

IMAGE_TAG=${IMAGE_TAG:-1.0.0}
echo "IMAGE_TAG: "$IMAGE_TAG

docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
#docker push ${IMAGE_NAME}:${IMAGE_TAG}
#docker rmi ${IMAGE_NAME}:${IMAGE_TAG}