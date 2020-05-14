#!/bin/bash
docker build -t moficodes/ibmcloud-kubernetes-admin:$TRAVIS_COMMIT .

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker push moficodes/ibmcloud-kubernetes-admin:$TRAVIS_COMMIT