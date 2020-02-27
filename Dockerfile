FROM alpine:latest

RUN apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
ADD vultr-cloud-controller-manager /
CMD ["/vultr-cloud-controller-manager"]