FROM golang:1.26.6-alpine AS build

RUN apk add --no-cache git

WORKDIR /workspace

COPY . .
ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o vultr-cloud-controller-manager .

FROM dhi.io/static:20260909-alpine3.24@sha256:296ab7284ac616e1f03b9ae929852b968315242311da974c57de342894276418

COPY --from=build /workspace/vultr-cloud-controller-manager /usr/local/bin/vultr-cloud-controller-manager
ENTRYPOINT ["/usr/local/bin/vultr-cloud-controller-manager"]
