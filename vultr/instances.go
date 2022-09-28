// Package vultr is vultr cloud specific implementation
package vultr

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vultr/govultr/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

const instanceShutdownStatus = "stopped"

var _ cloudprovider.Instances = &instances{}

type instances struct {
	client *govultr.Client
}

func newInstances(client *govultr.Client) cloudprovider.Instances {
	return &instances{client}
}

// NodeAddresses return all IPv4 addresses associated to a instance by nodeName.
func (i *instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	instance, err := vultrByName(ctx, i.client, name)
	if err != nil {
		return nil, err
	}

	addresses, err := i.nodeAddresses(instance)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

// NodeAddressesByProviderID return all IPv4 addresses associated to a instance by ProviderID.
func (i *instances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	id, err := vultrIDFromProviderID(providerID)
	if err != nil {
		return nil, err
	}

	instance, err := vultrByID(ctx, i.client, id)
	if err != nil {
		return nil, err
	}

	addresses, err := i.nodeAddresses(instance)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

// InstanceID returns the instance ID of the droplet identified by nodeName.
func (i *instances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	instance, err := vultrByName(ctx, i.client, nodeName)
	if err != nil {
		return "", err
	}

	return instance.ID, nil
}

// InstanceType returns the type of instance for given name.
func (i *instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	instance, err := vultrByName(ctx, i.client, name)
	if err != nil {
		return "", err
	}

	return instance.Plan, nil
}

// InstanceTypeByProviderID returns the instance type for given providerID.
func (i *instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	id, err := vultrIDFromProviderID(providerID)
	if err != nil {
		return "", err
	}

	instance, err := vultrByID(ctx, i.client, id)
	if err != nil {
		return "", err
	}

	return instance.Plan, nil
}

// AddSSHKeyToAllInstances is not implemented.
func (i *instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return cloudprovider.NotImplemented
}

// CurrentNodeName returns the hostname as a NodeName.
func (i *instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

// InstanceExistsByProviderID returns true if the instance with the providerID is running
func (i *instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	id, err := vultrIDFromProviderID(providerID)
	if err != nil {
		return false, err
	}

	_, err = vultrByID(ctx, i.client, id)
	if err == nil {
		return true, nil
	}

	return false, nil
}

// InstanceShutdownByProviderID returns true if the instance is turned off
func (i *instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	id, err := vultrIDFromProviderID(providerID)
	if err != nil {
		return false, err
	}

	instance, err := vultrByID(ctx, i.client, id)
	if err != nil {
		return false, err
	}

	return instance.PowerStatus == instanceShutdownStatus, nil
}

// nodeAddresses gathers public/private IP addresses and returns a []v1.NodeAddress .
func (i *instances) nodeAddresses(instance *govultr.Instance) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress

	addresses = append(addresses, v1.NodeAddress{
		Type:    v1.NodeHostName,
		Address: instance.Label,
	})

	// make sure we have either pubic and private ip
	if instance.InternalIP == "" || instance.MainIP == "" {
		return nil, errors.New("require both public and private IP")
	}

	addresses = append(addresses,
		v1.NodeAddress{Type: v1.NodeInternalIP, Address: instance.InternalIP}, // private IP
		v1.NodeAddress{Type: v1.NodeExternalIP, Address: instance.MainIP},     // public IP
	)

	return addresses, nil
}

// vultrIDFromProviderID returns a vultr instance ID from providerID.
func vultrIDFromProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("providerID cannot be an empty string")
	}

	split := strings.Split(providerID, "://")
	if len(split) != 2 { //nolint
		return "", fmt.Errorf("unexpected providerID format %s, expected format to be: vultr://abc123", providerID)
	}

	if split[0] != ProviderName {
		return "", fmt.Errorf("provider scheme from providerID should be 'vultr://', %s", providerID)
	}
	return split[1], nil
}

// vultrByID returns a vultr instance for the given id.
func vultrByID(ctx context.Context, client *govultr.Client, id string) (*govultr.Instance, error) {
	instance, err := client.Instance.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return instance, err
}

// vultrByName returns a vultr instance for a given NodeName.
// Note that if multiple nodes with the same name exist and error will be thrown.
func vultrByName(ctx context.Context, client *govultr.Client, nodeName types.NodeName) (*govultr.Instance, error) {
	listOptions := &govultr.ListOptions{}
	var instances []govultr.Instance
	for {
		i, meta, err := client.Instance.List(ctx, listOptions)
		if err != nil {
			return nil, err
		}

		for _, v := range i { //nolint
			if v.Label == string(nodeName) {
				instances = append(instances, v)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			listOptions.Cursor = meta.Links.Next
			continue
		}
	}

	if len(instances) == 0 {
		return nil, cloudprovider.InstanceNotFound
	} else if len(instances) > 1 {
		return nil, errors.New(fmt.Sprintf("Multiple instances found with name %v", nodeName))
	}

	return &instances[0], nil
}
