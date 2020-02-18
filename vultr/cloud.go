package vultr

import (
	"fmt"
	"io"
	"os"

	"github.com/vultr/govultr"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

const (
	providerName   = "vultr"
	accessTokenEnv = "VULTR_API_KEY"
	regionEnv      = "VULTR_REGION"
)

type cloud struct {
	client        *govultr.Client
	instances     cloudprovider.Instances
	zones         cloudprovider.Zones
	loadbalancers cloudprovider.LoadBalancer
}

func init() {
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (i cloudprovider.Interface, err error) {
		return newCloud()
	})
}

func newCloud() (cloudprovider.Interface, error) {
	apiToken := os.Getenv(accessTokenEnv)
	if apiToken == "" {
		return nil, fmt.Errorf("%s must be set in the environment (use a k8s secret)", accessTokenEnv)
	}

	region := os.Getenv(regionEnv)
	if region == "" {
		return nil, fmt.Errorf("%s must be set in the environment (use a k8s secret)", regionEnv)
	}

	vultr := govultr.NewClient(nil, apiToken)
	vultr.SetUserAgent(fmt.Sprintf("vultr-cloud-controller-manager %s", vultr.UserAgent))

	return &cloud{
		client:    vultr,
		instances: newInstances(vultr),
		zones:     newZones(vultr, region),
		//loadbalancers: nil,
	}, nil
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
	klog.V(5).Info("called Zones")
	return c.zones, true
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	klog.V(5).Info("called Clusters")
	return nil, false
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
	return false
}
