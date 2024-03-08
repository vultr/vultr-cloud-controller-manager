package vultr

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoadbalancers_GetLoadBalancer(t *testing.T) {
	client := newFakeClient()

	lb := newLoadbalancers(client, "ewr")

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "lb-name",
			Namespace: v1.NamespaceDefault,
			UID:       "lb-name",
			Annotations: map[string]string{
				annoVultrLoadBalancerID: "abc123",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(8080),
				},
			},
		},
	}

	expected := &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{
		{
			IP:       "192.168.0.1",
			Hostname: "albname",
		},
	}}

	actual, exists, err := lb.GetLoadBalancer(context.TODO(), "lbname", svc)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !exists {
		t.Error("expected true got false")
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v got %+v", expected, actual)
	}
}

func TestLoadbalancers_GetLoadBalancerName(t *testing.T) {
	client := newFakeClient()

	lb := newLoadbalancers(client, "1")

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "lb-name",
			Namespace: v1.NamespaceDefault,
			UID:       "lb-name",
			Annotations: map[string]string{
				annoVultrLoadBalancerID: "abc123",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(8080),
				},
			},
		},
	}

	actual := lb.GetLoadBalancerName(context.Background(), "cluster-name", svc)

	if actual != "albname" {
		t.Errorf("expected lbname got %s", actual)
	}
}

func TestLoadbalancers_EnsureLoadBalancer(t *testing.T) {
	client := newFakeClient()
	lb := newLoadbalancers(client, "1")

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "lb-name",
			Namespace: v1.NamespaceDefault,
			UID:       "lb-name",
			Annotations: map[string]string{
				annoVultrFirewallRules:  "cloudflare,80;10.0.0.0/8,80",
				annoVultrNodeCount:      "5",
				annoVultrLoadBalancerID: "abc123",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(8080),
				},
			},
		},
	}
	expected := &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{
		{
			IP:       "192.168.0.1",
			Hostname: "albname",
		},
	}}

	nodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node1",
			},
			Spec: v1.NodeSpec{
				ProviderID: "vultr://123",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node2",
			},
			Spec: v1.NodeSpec{
				ProviderID: "vultr://124",
			},
		},
	}

	actual, err := lb.EnsureLoadBalancer(context.Background(), "cluster-name", svc, nodes)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v got %+v", expected, actual)
	}
	// todo go through does not exist route
}

func TestLoadbalancers_UpdateLoadBalancer(t *testing.T) {
	client := newFakeClient()
	lb := newLoadbalancers(client, "1")

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "lb-name",
			Namespace: v1.NamespaceDefault,
			UID:       "lb-name",
			Annotations: map[string]string{
				annoVultrLoadBalancerID: "abc123",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(8080),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node1",
			},
			Spec: v1.NodeSpec{
				ProviderID: "vultr://123",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node2",
			},
			Spec: v1.NodeSpec{
				ProviderID: "vultr://124",
			},
		},
	}

	err := lb.UpdateLoadBalancer(context.Background(), "cluster-name", svc, nodes)
	if err != nil {
		t.Errorf("expected nil got %s", err.Error())
	}
}

func TestLoadbalancers_EnsureLoadBalancerDeleted(t *testing.T) {
	client := newFakeClient()
	lb := newLoadbalancers(client, "1")

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "lb-name-deleted",
			Namespace:   v1.NamespaceDefault,
			UID:         "lb-name-deleted",
			Annotations: nil,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(8080),
				},
			},
		},
	}

	err := lb.EnsureLoadBalancerDeleted(context.Background(), "cluster-name", svc)
	if err != nil {
		t.Errorf("expected nil got %s", err.Error())
	}
}
