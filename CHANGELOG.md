# Change Log

## v0.1.3 (2021-04-26)
* Fix error message on metadata retrieval 

[CCM Container v0.1.3](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)


## v0.1.2 (2021-03-25)
* Bumped vultr/metadata from v1.0.2 - v1.0.3

[CCM Container v0.1.2](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)

## v0.1.1 (2020-12-14)
* Bumped GoVultr from v2.0.0 - v2.2.0 
* Adding proxy protocol support for Load Balancer service [24](https://github.com/vultr/vultr-cloud-controller-manager/pull/34)
* Adding nightly yaml

[CCM Container v0.1.1](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)


## v0.1.0 (2020-11-24)
* Bumping Kubernetes dependencies from 1.18.5 to 1.19.4 
* Bumped GoVultr to v2.0.0 - This will use UUIDs for node IDs and not work on clusters provisioned by a CCM prior to v0.1.0

[CCM Container v0.1.0](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)


## v0.0.5 (2020-06-30)
* Bumping Kubernetes dependencies from 1.17.5 to 1.18.5 

[CCM Container v0.0.5](https://hub.docker.com/layers/vultr/vultr-cloud-controller-manager/v0.0.5/images/sha256-db70482087faa632e4852ddd69ad1586f2efdf0876daae2ace158d7f0721cf2f?context=repo)

## v0.0.4 (2020-06-01)
* Bumping Vultr Metadata client to v1.0.1 to support new Region South Korea 

[CCM Container v0.0.4](https://hub.docker.com/layers/vultr/vultr-cloud-controller-manager/v0.0.4/images/sha256-050a3bf2cf1726caa1295831a6f50b24efc10da2d76ea98a24f79d20bf8c294b?context=repo)

## v0.0.3 (2020-05-21)
* Updated Kubernetes dependency to 1.17.5
* Updated GoVultr dependency to v0.4.1
* Added Metadata client dependency which removes need to define region ID in yaml 
* Added a more defined RBAC

[CCM Container v0.0.3](https://hub.docker.com/layers/vultr/vultr-cloud-controller-manager/v0.0.3/images/sha256-bde33d08802dd9211d3faa66007639e605eded89d13d77ba9cd4cfae9161f6e9?context=repo)


## v0.0.2 (2020-03-13)
* Support LoadBalancer

[CCM Container v0.0.2](https://hub.docker.com/layers/vultr/vultr-cloud-controller-manager/v0.0.2/images/sha256-96c6ed0293fb6c444dfcf927d775798a1eec3f2de39e2155600677441531e4a8?context=repo)

## v0.0.1 (2020-02-28)
* Initial Release supports 
    * NodeController
    * ZoneController

[CCM Container v0.0.1](https://hub.docker.com/layers/vultr/vultr-cloud-controller-manager/v0.0.1/images/sha256-fc4e02792fa9794b41bedf2a9472ba755f6c68c7eca59d1951f53d2b61cd48a8?context=repo)
