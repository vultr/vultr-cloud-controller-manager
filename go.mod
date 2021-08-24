module github.com/vultr/vultr-cloud-controller-manager

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/vultr/govultr/v2 v2.5.1
	github.com/vultr/metadata v1.0.3
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	k8s.io/cloud-provider v0.21.3
	k8s.io/component-base v0.21.3
	k8s.io/klog/v2 v2.8.0
)
