DOCKER_IMAGE_VER=0.1

DOCKER_CONTAINER=k8sutils:${DOCKER_IMAGE_VER}

.PHONY: build
build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .

.PHONY: img-build
img-build:
	docker build -t ${DOCKER_CONTAINER} .

.PHONY: img-push
img-push:
	docker tag ${DOCKER_CONTAINER} docker.io/dmitsh/${DOCKER_CONTAINER} && docker push docker.io/dmitsh/${DOCKER_CONTAINER}
