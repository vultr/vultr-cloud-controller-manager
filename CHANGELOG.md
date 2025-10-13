# Change Log
## [v0.15.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.14.0...v0.15.0) (2025-10-10)
### Dependencies
* Bump github.com/spf13/pflag from 1.0.6 to 1.0.10 [PR 328](https://github.com/vultr/vultr-cloud-controller-manager/pull/328)
* Bump github.com/vultr/govultr/v3 from 3.15.0 to 3.24.0 [PR 327](https://github.com/vultr/vultr-cloud-controller-manager/pull/327) 
* Bump golang.org/x/oauth2 from 0.27.0 to 0.32.0 [PR 327](https://github.com/vultr/vultr-cloud-controller-manager/pull/325) 
* Bump k8s.io/cloud-provider from 0.31.1 to 0.34.1 [PR 322](https://github.com/vultr/vultr-cloud-controller-manager/pull/322) 

### Enhancements
* Add Auto SSL annotation support [PR 296](https://github.com/vultr/vultr-cloud-controller-manager/pull/296)
* Make logs quote IDs in output [PR 320](https://github.com/vultr/vultr-cloud-controller-manager/pull/320)
* Change kubernetes updates calls to patch [PR 312](https://github.com/vultr/vultr-cloud-controller-manager/pull/312)
* Add more robust ID checking for load balancers [PR 310](https://github.com/vultr/vultr-cloud-controller-manager/pull/310)
* Update CCM daemonset tolerations to use affinity [PR 297](https://github.com/vultr/vultr-cloud-controller-manager/pull/297)
* Support non VPC VM instances [PR 294](https://github.com/vultr/vultr-cloud-controller-manager/pull/294)

### Automation
* Migrate golangci-lint config to v2 [PR 326](https://github.com/vultr/vultr-cloud-controller-manager/pull/326)

### Documentation
* Add notice for VPC2 deprecation [PR 329](https://github.com/vultr/vultr-cloud-controller-manager/pull/329)

### New Contributors
* @vrabbi made their first contribution in [PR 294](https://github.com/vultr/vultr-cloud-controller-manager/pull/294)

## [v0.14.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.13.3...v0.14.0) (2025-03-06)

### Dependencies
* Update go to v1.24 and github workflows [PR 288](https://github.com/vultr/vultr-cloud-controller-manager/pull/280)

* Update govultr to v3.15.0 [PR 290](https://github.com/vultr/vultr-cloud-controller-manager/pull/290) 

### Enhancements
* Add HTTP2/3 and Timeout annonation support  [PR 230](https://github.com/vultr/vultr-cloud-controller-manager/pull/291)

* Handle "Invalid server" error for bare metal check [PR 287](https://github.com/vultr/vultr-cloud-controller-manager/pull/287)

### Automation

* Remove deprecated linters [PR 289](https://github.com/vultr/vultr-cloud-controller-manager/pull/289)

## [v0.13.3](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.13.2...v0.13.3) (2024-11-04)

* Update govultr from v3.9.1 to v3.11.2 [PR 283](https://github.com/vultr/vultr-cloud-controller-manager/pull/283)

## [v0.13.2](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.13.1...v0.13.2) (2024-10-28)

### Dependencies
* Update go to v1.23 and github workflows [PR 280](https://github.com/vultr/vultr-cloud-controller-manager/pull/280)

## [v0.13.1](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.13.0...v0.13.1) (2024-08-19)
### Bug Fixes
* Load Balancers: Resolve issues with hairpinning by introducing hostname workaround [PR 268](https://github.com/vultr/vultr-cloud-controller-manager/pull/268)

## [v0.13.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.12.0...v0.13.0) (2024-08-12)
### Bug Fixes
* Load Balancers: Get service before checking annotations to resolve nil map errors [PR 262](https://github.com/vultr/vultr-cloud-controller-manager/pull/262)

### Dependencies
* Update go to v1.21 [PR 248](https://github.com/vultr/vultr-cloud-controller-manager/pull/248)
* Bump golang.org/x/net from 0.22.0 to 0.23.0 [PR 245](https://github.com/vultr/vultr-cloud-controller-manager/pull/245)
* Update otel from v1.18.0 to v1.20.0 [PR 258](https://github.com/vultr/vultr-cloud-controller-manager/pull/258)
* Bump github.com/vultr/govultr/v3 from 3.6.4 to 3.8.1 [PR 259](https://github.com/vultr/vultr-cloud-controller-manager/pull/259)
* Bump k8s.io/klog/v2 from 2.120.1 to 2.130.1 [PR 261](https://github.com/vultr/vultr-cloud-controller-manager/pull/261)

### Automation
* Update mattermost notification workflows [PR 246](https://github.com/vultr/vultr-cloud-controller-manager/pull/246)
* Fix mattermost notifications [PR 247](https://github.com/vultr/vultr-cloud-controller-manager/pull/247)
* Goreleaser and golangci-lint workflows updates [PR 253](https://github.com/vultr/vultr-cloud-controller-manager/pull/253)

### New Contributors
* @mondragonfx made their first contribution in [PR 248](https://github.com/vultr/vultr-cloud-controller-manager/pull/248)

## [v0.12.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.11.0...v0.12.0) (2024-04-02)
### Enhancements
* Add ability for load balancer labels to be updated [PR 230](https://github.com/vultr/vultr-cloud-controller-manager/pull/230)

### Documentation
* Document https value option to LB protocol annotation [PR 231](https://github.com/vultr/vultr-cloud-controller-manager/pull/231)

### Dependencies
* Bump golang.org/x/crypto from 0.13.0 to 0.17.0 [PR 223](https://github.com/vultr/vultr-cloud-controller-manager/pull/223)
* Bump golang.org/x/oauth2 from 0.12.0 to 0.18.0 [PR 232](https://github.com/vultr/vultr-cloud-controller-manager/pull/232)
* Bump google.golang.org/grpc from 1.58.0 to 1.58.3 [PR 216](https://github.com/vultr/vultr-cloud-controller-manager/pull/216)
* Bump k8s.io/klog/v2 from 2.100.1 to 2.120.1 [PR 233](https://github.com/vultr/vultr-cloud-controller-manager/pull/233)
* Bump google.golang.org/protobuf from 1.31.0 to 1.33.0 [PR 234](https://github.com/vultr/vultr-cloud-controller-manager/pull/234)
* Bump github.com/vultr/govultr/v3 from 3.6.1 to 3.6.4 [PR 237](https://github.com/vultr/vultr-cloud-controller-manager/pull/237)

## [v0.11.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.10.0...v0.11.0) (2024-01-16)
### Enhancements
* Add annotation to create a load balancer on a service [PR 213](https://github.com/vultr/vultr-cloud-controller-manager/pull/213)
* Add check for firewall type of IPv6 [PR 212](https://github.com/vultr/vultr-cloud-controller-manager/pull/212)
* Add support for bare metal nodes [PR 228](https://github.com/vultr/vultr-cloud-controller-manager/pull/228)

## [v0.10.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.9.0...v0.10.0) (2023-09-18)
* Update to Go v1.20 [PR #186](https://github.com/vultr/vultr-cloud-controller-manager/pull/186)
* Update to Go-Vultr v3.1.0 [PR #191](https://github.com/vultr/vultr-cloud-controller-manager/pull/191)
* Add VLB node-count [PR #206](https://github.com/vultr/vultr-cloud-controller-manager/pull/206)
* Add Secret watcher to update services when TLS cert is renewed for VLB [PR #209](https://github.com/vultr/vultr-cloud-controller-manager/pull/209)

## [v0.9.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.8.2...v0.9.0) (2023-03-18)
* Added IPv6 support for VLB and worker nodes [PR #163](https://github.com/vultr/vultr-cloud-controller-manager/pull/163)

## [v0.8.2](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.8.1...v0.8.2) (2023-02-08)
* Increase page size for instance lookup and add additional error check during instance exists check [PR #154](https://github.com/vultr/vultr-cloud-controller-manager/pull/154)

## [v0.8.1](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.8.0...v0.8.1) (2023-02-02)
* Update instanceV2 to fix instance lookup [PR #152](https://github.com/vultr/vultr-cloud-controller-manager/pull/152)

## [v0.8.0](https://github.com/vultr/vultr-cloud-controller-manager/compare/v0.7.0...v0.8.0) (2023-02-01)
* Add support for cloudflare source in loadbalaner firewall rules [PR #139](https://github.com/vultr/vultr-cloud-controller-manager/pull/139)
* implemented instancesv2 [PR #150](https://github.com/vultr/vultr-cloud-controller-manager/pull/150)

### New Contributors
* @reubit made their first contribution in [PR #139](https://github.com/vultr/vultr-cloud-controller-manager/pull/139)
* @happytreees made their first contribution in [PR #150](https://github.com/vultr/vultr-cloud-controller-manager/pull/150)

## v0.7.0 (2022-09-30)
* Adds goreleaser: [PR #126](https://github.com/vultr/vultr-cloud-controller-manager/pull/126)
* Updates various `k8s.io` components from `v0.22.9` to `v0.24.6`: [PR #120](https://github.com/vultr/vultr-cloud-controller-manager/pull/120)
* Moves from Go 1.17 to Go 1.19: [PR #120](https://github.com/vultr/vultr-cloud-controller-manager/pull/120)
* Adds `golangci-lint`: [PR #120](https://github.com/vultr/vultr-cloud-controller-manager/pull/120)
* Updates `klog` version: [PR #120](https://github.com/vultr/vultr-cloud-controller-manager/pull/120)
* Removes `github.com/pkg/errors` in favor of built-in `fmt`: [PR #120](https://github.com/vultr/vultr-cloud-controller-manager/pull/120)
* Updates `golang.org/x/oauth2`: [PR #120](https://github.com/vultr/vultr-cloud-controller-manager/pull/120)
* Fixes a lot of `golangci-lint` issues: [PR #120](https://github.com/vultr/vultr-cloud-controller-manager/pull/120)
* Adds annotation to define "backend" protocol: [PR #118](https://github.com/vultr/vultr-cloud-controller-manager/pull/118)
* Updates to `govultr` from `2.16.0` to `2.17.2`: [PR #103](https://github.com/vultr/vultr-cloud-controller-manager/pull/103)

[CCM Container v0.7.0](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)

## v0.6.0 (2021-05-13)
* VPC changes by @ddymko in https://github.com/vultr/vultr-cloud-controller-manager/pull/87
* bumping the 22 patches from 7 to 9 by @ddymko in https://github.com/vultr/vultr-cloud-controller-manager/pull/89
* Allow changing the base URL the ccm uses for API calls by @ddymko in https://github.com/vultr/vultr-cloud-controller-manager/pull/88
* bump go to 1.17 by @ddymko in https://github.com/vultr/vultr-cloud-controller-manager/pull/91
* Bump github.com/vultr/metadata from 1.0.3 to 1.1.0 by @dependabot in https://github.com/vultr/vultr-cloud-controller-manager/pull/93
* updating all go 1.16 references to 1.17 by @ddymko in https://github.com/vultr/vultr-cloud-controller-manager/pull/94


[CCM Container v0.6.0](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)

## v0.5.0 (2021-03-11)
* Updated GoVultr dependencies from 2.11.1 to 2.14.1.
* Updated Klog to v2.9.0
* Updated Kubernetes dependencies from v0.21 to v0.22

[CCM Container v0.5.0](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)

## v0.4.0 (2021-11-29)
* Updated GoVultr dependencies from 2.5.1 to 2.11.1. This fixes LB issues with setting SSL secret data
* Read and store `kubeconfig` when it is passed in as a flag. This fixes and issue when the CCM runs to grab secret data for LB SSL

[CCM Container v0.4.0](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)


## v0.3.0 (2021-08-24)
* Updated Kubernetes dependencies from 1.19 to 1.20

[CCM Container v0.3.0](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)


## v0.2.1 (2021-06-15)
* Adding ability to change UserAgent
* Include binary builds on tagged releases

[CCM Container v0.2.1](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)


## v0.2.0 (2021-05-19)
* Bumped GoVultr from v2.2.0 - v2.5.1
* Load Balancers updates - support for firewalls + private networks
* Bump to build with go version 1.16

[CCM Container v0.2.0](https://hub.docker.com/repository/docker/vultr/vultr-cloud-controller-manager)

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

