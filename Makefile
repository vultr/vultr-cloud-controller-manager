.PHONY: build-linux
build-linux:
	echo "building vultr ccm for linux"
	GOOS=linux GOARCH=amd64 GCO_ENABLED=0 go build -o dist/vultr-cloud-controller-manager .

docker-build:
	echo "building docker image"
	docker build . -t vultr/vultr-cloud-controller-manager