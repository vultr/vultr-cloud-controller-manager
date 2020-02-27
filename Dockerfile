FROM alpine:latest

RUN apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
ADD dist/vultr-cloud-controller-manager /
ENTRYPOINT ["/vultr-cloud-controller-manager"]