package main

import (
	_ "k8s.io/cloud-provider"

	_ "github.com/vultr/vultr-cloud-controller-manager/vultr"
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version" // for version metrics registration
)

func main() {
}
