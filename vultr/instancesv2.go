// Package vultr is vultr cloud specific implementation
package vultr

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/vultr/govultr/v3"
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

// InstanceExists return bool whether the instance exists
func (i *instancesv2) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	bm := false

	if label, ok := node.Labels["vultr.com/baremetal"]; ok {
		if label == "true" {
			bm = true
		}
	}

	if bm {
		newNode, err := i.getVultrBareMetal(ctx, node)
		if err != nil {
			log.Printf("baremetal(%s) exists check failed: %e", node.Spec.ProviderID, err) //nolint
			if strings.Contains(err.Error(), "invalid baremetal ID") {
				return false, nil
			}
			if strings.Contains(err.Error(), "baremetal not found") {
				return false, nil
			}
			if strings.Contains(err.Error(), "Invalid server") {
				return false, nil
			}
			return false, err
		}
		if newNode.Status == ACTIVE || newNode.Status == PENDING {
			log.Printf("baremetal(%s) status is: %s", newNode.Label, newNode.Status) //nolint
			return true, nil
		}
	} else {
		newNode, err := i.getVultrInstance(ctx, node)
		if err != nil {
			log.Printf("instance(%s) exists check failed: %e", node.Spec.ProviderID, err) //nolint
			if strings.Contains(err.Error(), "invalid instance ID") {
				return false, nil
			}
			if strings.Contains(err.Error(), "instance not found") {
				return false, nil
			}
			return false, err
		}
		if newNode.Status == ACTIVE || newNode.Status == PENDING || newNode.Status == RESIZING {
			log.Printf("instance(%s) status is: %s", newNode.Label, newNode.Status) //nolint
			return true, nil
		}
	}

	return false, nil
}

// InstanceShutdown returns bool whether the instance is running or powered off
func (i *instancesv2) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	bm := false

	if label, ok := node.Labels["vultr.com/baremetal"]; ok {
		if label == "true" {
			bm = true
		}
	}

	if bm {
		newNode, err := i.getVultrBareMetal(ctx, node)
		if err != nil {
			log.Printf("baremetal(%s) shutdown check failed: %e", node.Spec.ProviderID, err) //nolint
			return false, err
		}
		if newNode.Status == ACTIVE || newNode.Status == PENDING { //nolint
			return false, nil
		}
	} else {
		newNode, err := i.getVultrInstance(ctx, node)
		if err != nil {
			log.Printf("instance(%s) shutdown check failed: %e", node.Spec.ProviderID, err) //nolint
			return false, err
		}
		if newNode.PowerStatus != "running" {
			return true, nil
		}
	}

	return false, nil
}

// InstanceMetadata returns a struct of type InstanceMetadata containing the node information
func (i *instancesv2) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	bm := false

	if label, ok := node.Labels["vultr.com/baremetal"]; ok {
		if label == "true" {
			bm = true
		}
	}

	if bm {
		newNode, err := i.getVultrBareMetal(ctx, node)
		if err != nil {
			log.Printf("baremetal(%s) metadata check failed: %e", node.Spec.ProviderID, err) //nolint
			return nil, err
		}
		nodeAddress, err := i.nodeBareMetalAddresses(newNode)
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
	newNode, err := i.getVultrInstance(ctx, node)
	if err != nil {
		log.Printf("instance(%s) metadata check failed: %e", node.Spec.ProviderID, err) //nolint
		return nil, err
	}
	nodeAddress, err := i.nodeInstanceAddresses(newNode)
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

// nodeInstanceAddresses gathers public/private IP addresses and returns a []v1.NodeAddress .
func (i *instancesv2) nodeInstanceAddresses(instance *govultr.Instance) ([]v1.NodeAddress, error) {
        var addresses []v1.NodeAddress

        if reflect.DeepEqual(instance, *&govultr.Instance{}) { //nolint
                return nil, fmt.Errorf("instance is empty %v", instance)
        }

        addresses = append(addresses, v1.NodeAddress{
                Type:    v1.NodeHostName,
                Address: instance.Label,
        })

        // Check conditions for internal and main IP
        if instance.InternalIP == "" && instance.MainIP == "" {
                return nil, fmt.Errorf("require at least one of internal or public IP")
        }

        // Handle the case where both IPs are provided
        if instance.InternalIP != "" && instance.MainIP != "" {
                addresses = append(addresses,
                        v1.NodeAddress{Type: v1.NodeInternalIP, Address: instance.InternalIP}, // private IP
                        v1.NodeAddress{Type: v1.NodeExternalIP, Address: instance.MainIP},     // public IP
                )
        } else if instance.InternalIP == "" && instance.MainIP != "" {
                // If internal IP is empty but main IP is not, use main IP for both
                addresses = append(addresses,
                        v1.NodeAddress{Type: v1.NodeInternalIP, Address: instance.MainIP}, // treat main IP as internal IP
                        v1.NodeAddress{Type: v1.NodeExternalIP, Address: instance.MainIP}, // public IP
                )
        } else if instance.InternalIP != "" && instance.MainIP == "" {
                // If main IP is empty but internal IP is not, use internal IP for both
                addresses = append(addresses,
                        v1.NodeAddress{Type: v1.NodeInternalIP, Address: instance.InternalIP}, // private IP
                        v1.NodeAddress{Type: v1.NodeExternalIP, Address: instance.InternalIP}, // treat internal IP as external IP
                )
        }

        if instance.V6MainIP != "" {
                addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: instance.V6MainIP}) // IPv6
        }

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

		newNode, err := vultrByInstanceID(ctx, i.client, id)
		if err != nil {
			log.Printf("instance(%s) by ID failed: %e", node.Spec.ProviderID, err) //nolint
			return nil, err
		}
		return newNode, nil
	}
	newNode, err := vultrByInstanceName(ctx, i.client, types.NodeName(node.Name))
	if err != nil {
		log.Printf("instance(%s) by name failed: %e", node.Name, err) //nolint
		return nil, err
	}
	return newNode, nil
}
