FROM golang:1.17-alpine AS build

RUN apk add --no-cache git

WORKDIR /workspace

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -o vultr-cloud-controller-manager .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/vultr-cloud-controller-manager /
ENTRYPOINT ["/vultr-cloud-controller-manager"] 