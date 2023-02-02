// Package vultr is vultr cloud specific implementation
package vultr

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/vultr/govultr/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

var _ cloudprovider.InstancesV2 = &instancesv2{}

type instancesv2 struct {
	client *govultr.Client
}

const (
	PENDING  = "pending"  //nolint
	ACTIVE   = "active"   //nolint
	RESIZING = "resizing" //nolint
)

func newInstancesV2(client *govultr.Client) cloudprovider.InstancesV2 {
	return &instancesv2{client}
}

// InstanceExists return bool whether or not the instance exists
func (i *instancesv2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	newNode, err := i.getVultrInstance(ctx, node)
	if err != nil {
		log.Printf("instance(%s) exists check failed: %e", node.Spec.ProviderID, err) //nolint
		if strings.Contains(err.Error(), "invalid instance ID") {
			return false, nil
		}
		return false, err
	}

	if newNode.Status == ACTIVE || newNode.Status == PENDING || newNode.Status == RESIZING {
		log.Printf("instance(%s) status is: %s", newNode.Label, newNode.Status) //nolint
		return true, nil
	}
	return false, nil
}

// InstanceShutdown returns bool whether or not the instance is running or powered off
func (i *instancesv2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	newNode, err := i.getVultrInstance(ctx, node)
	if err != nil {
		log.Printf("instance(%s) shutdown check failed: %e", node.Spec.ProviderID, err) //nolint
		return false, err
	}

	if newNode.PowerStatus != "running" {
		return true, nil
	}
	return false, nil
}

// InstanceMetadata returns a struct of type InstanceMetadata containing the node information
func (i *instancesv2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	newNode, err := i.getVultrInstance(ctx, node)
	if err != nil {
		log.Printf("instance(%s) metadata check failed: %e", node.Spec.ProviderID, err) //nolint
		return nil, err
	}

	nodeAddress, err := i.nodeAddresses(newNode)
	if err != nil {
		return nil, err
	}

	vultrNode := cloudprovider.InstanceMetadata{
		InstanceType:  newNode.Plan,
		ProviderID:    fmt.Sprintf("vultr://%s", newNode.ID),
		Region:        newNode.Region,
		NodeAddresses: nodeAddress,
	}

	log.Printf("returned node metadata: %v", vultrNode) //nolint
	return &vultrNode, nil
}

// nodeAddresses gathers public/private IP addresses and returns a []v1.NodeAddress .
func (i *instancesv2) nodeAddresses(instance *govultr.Instance) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress

	if reflect.DeepEqual(instance, *&govultr.Instance{}) { //nolint
		return nil, fmt.Errorf("instance is empty %v", instance)
	}

	addresses = append(addresses, v1.NodeAddress{
		Type:    v1.NodeHostName,
		Address: instance.Label,
	})

	// make sure we have either pubic and private ip
	if instance.InternalIP == "" || instance.MainIP == "" {
		return nil, fmt.Errorf("require both public and private IP")
	}

	addresses = append(addresses,
		v1.NodeAddress{Type: v1.NodeInternalIP, Address: instance.InternalIP}, // private IP
		v1.NodeAddress{Type: v1.NodeExternalIP, Address: instance.MainIP},     // public IP
	)
	return addresses, nil
}

// getVultrInstance attempts to obtain Vultr Instance from Vultr API
func (i *instancesv2) getVultrInstance(ctx context.Context, node *v1.Node) (*govultr.Instance, error) {
	skipID := false

	if node.Spec.ProviderID == "" {
		skipID = true
	}

	if node.Name == "" {
		return nil, fmt.Errorf("node name cannot be empty")
	}

	if !skipID {
		id, err := vultrIDFromProviderID(node.Spec.ProviderID)
		if err != nil {
			log.Printf("instance(%s) provider split failed: %e", node.Spec.ProviderID, err) //nolint
			return nil, err
		}

		newNode, err := vultrByID(ctx, i.client, id)
		if err != nil {
			log.Printf("instance(%s) by ID failed: %e", node.Spec.ProviderID, err) //nolint
			return nil, err
		}
		return newNode, nil
	}
	newNode, err := vultrByName(ctx, i.client, types.NodeName(node.Name))
	if err != nil {
		log.Printf("instance(%s) by name failed: %e", node.Name, err) //nolint
		return nil, err
	}
	return newNode, nil
}
