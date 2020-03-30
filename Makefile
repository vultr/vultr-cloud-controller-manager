VERSION ?= $VERSION
REGISTRY ?= $REGISTRY

.PHONY: deploy
deploy: clean build-linux docker-build docker-push

.PHONY: build
build:
	@echo "building vultr ccm"
	go build -o dist/vultr-cloud-controller-manager .

.PHONY: build-linux
build-linux:
	@echo "building vultr ccm for linux"
	GOOS=linux GOARCH=amd64 GCO_ENABLED=0 go build -o dist/vultr-cloud-controller-manager .

.PHONY: docker-build
docker-build:
	@echo "building docker image to dockerhub $(REGISTRY) with version $(VERSION)"
	docker build . -t $(REGISTRY)/vultr-cloud-controller-manager:$(VERSION)

.PHONY: docker-push
docker-push:
	docker push $(REGISTRY)/vultr-cloud-controller-manager:$(VERSION)

.PHONY: clean
clean:
	go clean -i -x ./...

.PHONY: test
test:
	go test -race github.com/vultr/vultr-cloud-controller-manager/vultr -v