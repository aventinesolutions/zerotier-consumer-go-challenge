#!/usr/bin/zsh -x

export IMAGE_ID="${1}"
export RELEASE_TAG="0.1.0-beta${2}"

docker tag "${IMAGE_ID}" zerotier-consumer-service:"${RELEASE_TAG}"
docker tag "${IMAGE_ID}" europe-west1-docker.pkg.dev/aventine-k8s/aventine/zerotier-consumer-service:"${RELEASE_TAG}"
docker push europe-west1-docker.pkg.dev/aventine-k8s/aventine/zerotier-consumer-service:"${RELEASE_TAG}"