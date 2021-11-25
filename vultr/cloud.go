package vultr

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/metadata"
	"golang.org/x/oauth2"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

const (
	ProviderName   = "vultr"
	accessTokenEnv = "VULTR_API_KEY"
	userAgent      = "CCM_USER_AGENT"
)

// Options currently stores the Kubeconfig that was passed in.
// We can use this to extend any other flags that may have been passed in that we require
var Options struct {
	KubeconfigFlag *pflag.Flag
}

type cloud struct {
	client        *govultr.Client
	instances     cloudprovider.Instances
	zones         cloudprovider.Zones
	loadbalancers cloudprovider.LoadBalancer
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (i cloudprovider.Interface, err error) {
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
		return nil, fmt.Errorf("failed to retrieve metadata: %v", err)
	}

	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: apiToken,
	})
	client := oauth2.NewClient(context.Background(), tokenSrc)

	vultr := govultr.NewClient(client)

	ua := os.Getenv(userAgent)
	if ua != "" {
		vultr.SetUserAgent(fmt.Sprintf("vultr-cloud-controller-manager:%s", ua))
	} else {
		vultr.SetUserAgent(fmt.Sprintf("vultr-cloud-controller-manager:%s", vultr.UserAgent))
	}

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
	return ProviderName
}

func (c *cloud) HasClusterID() bool {
	klog.V(5).Info("called HasClusterID")
	return false
}
