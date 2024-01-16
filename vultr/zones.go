// Package vultr is vultr cloud specific implementation
package vultr

import (
	"context"

	"github.com/vultr/govultr/v3"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

var _ cloudprovider.Zones = &zones{}

type zones struct {
	client *govultr.Client
	region string
}

func newZones(client *govultr.Client, zone string) cloudprovider.Zones {
	return zones{client, zone}
}

func (z zones) GetZone(_ context.Context) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{Region: z.region}, nil
}

func (z zones) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	id, err := vultrIDFromProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, nil
	}
	vultr, err := vultrByInstanceID(ctx, z.client, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: vultr.Region}, nil
}

func (z zones) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	vultr, err := vultrByInstanceName(ctx, z.client, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, nil
	}

	return cloudprovider.Zone{Region: vultr.Region}, nil
}
