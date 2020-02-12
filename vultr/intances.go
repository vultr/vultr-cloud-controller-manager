package vultr

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type instances struct{
}

func (i *instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	panic("implement me")
}

func (i *instances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	panic("implement me")
}

func (i *instances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	panic("implement me")
}

func (i *instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	panic("implement me")
}

func (i *instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	panic("implement me")
}

func (i *instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	panic("implement me")
}

func (i *instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	panic("implement me")
}

func (i *instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	panic("implement me")
}

func (i *instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	panic("implement me")
}
