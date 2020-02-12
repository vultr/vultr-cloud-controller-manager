package vultr

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cloud-provider"
)

type zones struct{

}

func (z zones) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	panic("implement me")
}

func (z zones) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	panic("implement me")
}

func (z zones) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	panic("implement me")
}
