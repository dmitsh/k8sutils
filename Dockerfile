FROM ubuntu:22.04

COPY k8sutils /usr/local/bin/k8sutils

ENTRYPOINT /usr/local/bin/k8sutils
