package vultr

import (
	"github.com/vultr/govultr"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

const (
	providerName   = "vultr"
	accessTokenENV = "VULTR_API_KEY"
)

type cloud struct {
	client        *govultr.Client
	instances     cloudprovider.Instances
	zones         cloudprovider.Zones
	loadbalancers cloudprovider.LoadBalancer
}

func (c *cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return c.loadbalancers, true
}

func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	return c.instances, true
}

func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	return c.zones, true
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	klog.V(5).Info("called Clusters")
	panic("implement me")
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	klog.V(5).Info("called Routes")
	return nil, false
}

func (c *cloud) ProviderName() string {
	klog.V(5).Info("called ProviderName")
	return providerName
}

func (c *cloud) HasClusterID() bool {
	klog.V(5).Info("called HasClusterID")
	panic("implement me")
}
