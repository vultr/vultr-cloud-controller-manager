FROM golang:1.23-alpine AS build

RUN apk add --no-cache git

WORKDIR /workspace

COPY . .
ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o vultr-cloud-controller-manager .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/vultr-cloud-controller-manager /usr/local/bin/vultr-cloud-controller-manager
ENTRYPOINT ["vultr-cloud-controller-manager"] 
