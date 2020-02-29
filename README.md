# Kubernetes Cloud Controller Manager for Vultr

The Vultr Cloud Controller Manager (ccm) provides a fully supported experience of Vultr features in your Kubernetes cluster.

This project is currently in active development and is not feature complete. There is a current milestone that outlines what is left remaining to be implemented [Milestone v0.1.0](https://github.com/vultr/vultr-cloud-controller-manager/milestone/1).


## Development 

Go minimum version `1.12.0`

The `vultr-cloud-controller-manager` uses go modules for its dependencies.

### Building the Binary

Since the `vultr-cloud-controller-manager` is meant to run inside a kubernetes cluster you will need to build the binary to be Linux specific.

`GOOS=linux GOARCH=amd64 go build -o dist/vultr-cloud-controller-manager .`

or by using our `Makefile`

`make build-linux`

This will build the binary and output it to a `dist` folder.

### Building the Docker Image

To build a docker image of the `vultr-cloud-controller-manager` you can run either

`make docker-build`

Running the image

`docker run -ti vultr/vultr-cloud-controller-manager`

### Deploying to a kubernetes cluster

You will need to make sure that your kubernetes cluster is configured to interact with a `external cloud provider`

More can be read about this in the [Running Cloud Controller](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/)