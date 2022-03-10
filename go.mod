module github.com/vultr/vultr-cloud-controller-manager

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/vultr/govultr/v2 v2.14.1
	github.com/vultr/metadata v1.0.3
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	k8s.io/cloud-provider v0.22.7
	k8s.io/component-base v0.22.7
	k8s.io/klog/v2 v2.30.0
)
