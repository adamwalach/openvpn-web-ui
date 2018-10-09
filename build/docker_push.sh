#!/bin/bash -e
echo "Travis tag: $TRAVIS_TAG"

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push awalach/openvpn-web-ui:$TRAVIS_TAG
