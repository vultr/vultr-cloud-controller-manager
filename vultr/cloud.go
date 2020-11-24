package vultr

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vultr/govultr/v2"
	"github.com/vultr/metadata"
	"golang.org/x/oauth2"
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

	mClient := metadata.NewClient()
	meta, err := mClient.Metadata()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve region from metadata: %v", meta)
	}

	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: apiToken,
	})
	client := oauth2.NewClient(context.Background(), tokenSrc)

	vultr := govultr.NewClient(client)
	vultr.SetUserAgent(fmt.Sprintf("vultr-cloud-controller-manager %s", vultr.UserAgent))

	return &cloud{
		client:        vultr,
		instances:     newInstances(vultr),
		zones:         newZones(vultr, strings.ToLower(meta.Region.RegionCode)),
		loadbalancers: newLoadbalancers(vultr, strings.ToLower(meta.Region.RegionCode)),
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

func (c *cloud) InstancesV2() (cloudprovider.InstancesV2, bool) {
	// TODO we will need to implement this but for now it is not required and experimental
	return nil, false
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
