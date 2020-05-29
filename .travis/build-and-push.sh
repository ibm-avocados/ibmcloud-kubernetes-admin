#!/bin/bash
echo "Building for Commit : $TRAVIS_COMMIT  and Tag : $TRAVIS_TAG"
docker build -t moficodes/ibmcloud-kubernetes-admin:$TRAVIS_TAG -f docker/Dockerfile.web .
docker build -t moficodes/ibmcloud-kubernetes-cron:$TRAVIS_TAG -f docker/Dockerfile.cron .

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker push moficodes/ibmcloud-kubernetes-admin:$TRAVIS_TAG
docker push moficodes/ibmcloud-kubernetes-cron:$TRAVIS_TAG