package vultr

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vultr/govultr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"strings"
)

type instances struct {
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

func vultrIDFromProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("providerID cannot be an empty string")
	}

	split := strings.Split(providerID, "://")
	if len(split) != 2 {
		return "", fmt.Errorf("unexpected providerID format %s, expected format to be: vultr://abc123", providerID)
	}

	if split[0] != providerName {
		return "", fmt.Errorf("provider scheme from providerID should be 'vultr://', %s", providerID)
	}
	return split[1], nil
}

func vultrByID(ctx context.Context, client *govultr.Client, id string) (*govultr.Server, error) {
	instance, err := client.Server.GetServer(ctx, id)
	if err != nil {
		return nil, err
	}
	return instance, err
}

func vultrByName(ctx context.Context, client *govultr.Client, nodeName types.NodeName) (*govultr.Server, error) {
	instances, err := client.Server.ListByLabel(ctx, string(nodeName))
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, cloudprovider.InstanceNotFound
	} else if len(instances) > 1 {
		return nil, errors.New(fmt.Sprintf("Multiple instances found with name %v", nodeName))
	}

	return &instances[0], nil
}
