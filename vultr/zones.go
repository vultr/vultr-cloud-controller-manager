package vultr

import (
	"context"

	"github.com/vultr/govultr"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cloud-provider"
)

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
	vultr, err := vultrByID(ctx, z.client, id)
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{Region: vultr.RegionID}, nil
}

func (z zones) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	vultr, err := vultrByName(ctx, z.client, nodeName)
	if err != nil {
		return cloudprovider.Zone{}, nil
	}

	return cloudprovider.Zone{Region: vultr.RegionID}, nil
}
