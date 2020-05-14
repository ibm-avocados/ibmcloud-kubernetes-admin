#!/bin/bash
echo "Building for Commit : $TRAVIS_COMMIT  and Tag : $TRAVIS_TAG"
docker build -t moficodes/ibmcloud-kubernetes-admin:$TRAVIS_TAG .

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker push moficodes/ibmcloud-kubernetes-admin:$TRAVIS_TAG