package vultr

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestInstances_InstanceExistsByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	expected := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "ccm-test",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.1.95.4",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "149.28.225.110",
		},
	}

	actual, err := instances.NodeAddressesByProviderID(context.TODO(), "vultr://576965")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expcted %+v got %+v", expected, actual)
	}
}

func TestInstances_InstanceShutdownByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceShutdownByProviderID(context.TODO(), "vultr://576965")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual {
		t.Errorf("expcted %+v got %+v", "false", "true")
	}
}

func TestInstances_InstanceTypeByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceTypeByProviderID(context.TODO(), "vultr://576965")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "204" {
		t.Errorf("expcted %+v got %+v", "204", actual)
	}
}

func TestInstances_NodeAddressesByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	expected := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "ccm-test",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.1.95.4",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "149.28.225.110",
		},
	}

	actual, err := instances.NodeAddresses(context.TODO(), "576965")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expcted %+v got %+v", expected, actual)
	}
}

func TestInstances_NodeAddresses(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	expected := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "ccm-test",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.1.95.4",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "149.28.225.110",
		},
	}

	actual, err := instances.NodeAddresses(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expcted %+v got %+v", expected, actual)
	}
}

func TestInstances_InstanceType(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceType(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "204" {
		t.Errorf("expcted %+v got %+v", "204", actual)
	}
}

func TestInstances_InstanceID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceID(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "576965" {
		t.Errorf("expcted %+v got %+v", "204", actual)
	}
}

func TestInstances_CurrentNodeName(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.CurrentNodeName(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "ccm-test" {
		t.Errorf("expcted %+v got %+v", "ccm-test", actual)
	}
}
