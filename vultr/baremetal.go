// Package vultr is vultr cloud specific implementation
package vultr

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/vultr/govultr/v3"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

func (i *instancesv2) getVultrBareMetal(ctx context.Context, node *v1.Node) (*govultr.BareMetalServer, error) {
	skipID := node.Spec.ProviderID == ""

	if node.Name == "" {
		return nil, fmt.Errorf("node name cannot be empty")
	}

	if !skipID {
		id, err := vultrIDFromProviderID(node.Spec.ProviderID)
		if err != nil {
			log.Printf("baremetal(%s) provider split failed: %e", node.Spec.ProviderID, err) //nolint
			return nil, err
		}
		bm, err := vultrByBareMetalID(ctx, i.client, id)
		if err != nil {
			log.Printf("baremetal(%s) could not be found: %e", id, err) //nolint
			return nil, err
		}

		return bm, nil
	}

	newNode, err := vultrByBareMetalName(ctx, i.client, types.NodeName(node.Name))
	if err != nil {
		log.Printf("baremetal(%s) by name failed: %e", node.Name, err) //nolint
		return nil, err
	}
	return newNode, nil
}

// vultrByBareMetalID returns a vultr baremetal for the given id.
func vultrByBareMetalID(ctx context.Context, client *govultr.Client, id string) (*govultr.BareMetalServer, error) {
	bm, _, err := client.BareMetalServer.Get(ctx, id) //nolint:bodyclose
	if err != nil {
		return nil, err
	}
	return bm, err
}

// vultrByBareMetalName returns a vultr bare metal for a given NodeName.
// Note that if multiple nodes with the same name exist and error will be thrown.
func vultrByBareMetalName(ctx context.Context, client *govultr.Client, nodeName types.NodeName) (*govultr.BareMetalServer, error) {
	listOptions := &govultr.ListOptions{PerPage: 300}

	var baremetals []govultr.BareMetalServer
	for {
		bm, meta, _, err := client.BareMetalServer.List(ctx, listOptions) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		for _, v := range bm { //nolint
			if v.Label == string(nodeName) {
				baremetals = append(baremetals, v)
			}
		}

		if meta.Links.Next == "" {
			break
		}

		listOptions.Cursor = meta.Links.Next
	}

	if len(baremetals) == 0 {
		return nil, cloudprovider.InstanceNotFound
	} else if len(baremetals) > 1 {
		return nil, fmt.Errorf("multiple baremetals found with name %v", nodeName)
	}

	return &baremetals[0], nil
}

// nodeBareMetalAddresses gathers public/private IP addresses and returns a []v1.NodeAddress .
func (i *instancesv2) nodeBareMetalAddresses(baremetal *govultr.BareMetalServer) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress

	if reflect.DeepEqual(baremetal, *&govultr.BareMetalServer{}) { //nolint
		return nil, fmt.Errorf("baremetal is empty %v", baremetal)
	}

	addresses = append(addresses, v1.NodeAddress{
		Type:    v1.NodeHostName,
		Address: baremetal.Label,
	})

	vpc2, _, err := i.client.BareMetalServer.ListVPC2Info(context.Background(), baremetal.ID) //nolint:bodyclose
	if err != nil {
		return nil, fmt.Errorf("error getting VPC2 info for bm %s", baremetal.Label)
	}

	for _, vpc := range vpc2 {
		addresses = append(addresses,
			v1.NodeAddress{Type: v1.NodeInternalIP, Address: vpc.IPAddress})
	}

	vpc1, _, err := i.client.BareMetalServer.ListVPCInfo(context.Background(), baremetal.ID) //nolint:bodyclose
	if err != nil {
		return nil, fmt.Errorf("error getting VPC1 info for bm %s", baremetal.Label)
	}

	for _, vpc := range vpc1 {
		addresses = append(addresses,
			v1.NodeAddress{Type: v1.NodeInternalIP, Address: vpc.IPAddress})
	}

	// make sure we have public ip
	if baremetal.MainIP == "" {
		return nil, fmt.Errorf("require both public IP")
	}

	addresses = append(addresses,
		v1.NodeAddress{Type: v1.NodeExternalIP, Address: baremetal.MainIP})

	if baremetal.V6MainIP != "" {
		addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: baremetal.V6MainIP}) // IPv6
	}

	return addresses, nil
}
