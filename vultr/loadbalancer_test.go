package vultr

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/govultr/v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
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

	lb.(*loadbalancers).kubeClient = &fake.Clientset{}

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

func TestLoadbalancers_UpdateLoadBalancer_SharedLabelMergesForwardingRules(t *testing.T) {
	for _, tc := range []struct {
		name             string
		annotations      map[string]string
		frontendProtocol string
		backendProtocol  string
		existingRuleID   string
		existingFrontend int
		existingBackend  int
		desiredFrontend  int32
		desiredBackend   int32
	}{
		{
			name:             "tcp",
			annotations:      map[string]string{annoVultrLBProtocol: protocolTCP},
			frontendProtocol: protocolTCP,
			backendProtocol:  protocolTCP,
			existingRuleID:   "rule-50001",
			existingFrontend: 50001,
			existingBackend:  30001,
			desiredFrontend:  50002,
			desiredBackend:   30002,
		},
		{
			name:             "udp",
			annotations:      map[string]string{annoVultrLBProtocol: protocolUDP},
			frontendProtocol: protocolUDP,
			backendProtocol:  protocolUDP,
			existingRuleID:   "rule-50001",
			existingFrontend: 50001,
			existingBackend:  30001,
			desiredFrontend:  50002,
			desiredBackend:   30002,
		},
		{
			name:             "http",
			annotations:      map[string]string{annoVultrLBProtocol: protocolHTTP},
			frontendProtocol: protocolHTTP,
			backendProtocol:  protocolHTTP,
			existingRuleID:   "rule-80",
			existingFrontend: 80,
			existingBackend:  30080,
			desiredFrontend:  8080,
			desiredBackend:   30081,
		},
		{
			name: protocolHTTPS,
			annotations: map[string]string{
				annoVultrLBProtocol:   protocolHTTP,
				annoVultrLBHTTPSPorts: "8443",
			},
			frontendProtocol: protocolHTTPS,
			backendProtocol:  protocolHTTPS,
			existingRuleID:   "rule-443",
			existingFrontend: 443,
			existingBackend:  30443,
			desiredFrontend:  8443,
			desiredBackend:   30444,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testSharedLabelMergeForwardingRules(t, tc.annotations, govultr.ForwardingRule{
				RuleID:           tc.existingRuleID,
				FrontendProtocol: tc.frontendProtocol,
				FrontendPort:     tc.existingFrontend,
				BackendProtocol:  tc.backendProtocol,
				BackendPort:      tc.existingBackend,
			}, govultr.ForwardingRule{
				FrontendProtocol: tc.frontendProtocol,
				FrontendPort:     int(tc.desiredFrontend),
				BackendProtocol:  tc.backendProtocol,
				BackendPort:      int(tc.desiredBackend),
			}, tc.desiredFrontend, tc.desiredBackend)
		})
	}
}

func testSharedLabelMergeForwardingRules(t *testing.T, annotations map[string]string, existingRule, expectedRule govultr.ForwardingRule, port, nodePort int32) {
	t.Helper()

	fakeLoadBalancer := &fakeLB{
		forwardingRules: []govultr.ForwardingRule{existingRule},
	}
	lb := &loadbalancers{
		client: &govultr.Client{LoadBalancer: fakeLoadBalancer},
		zone:   "ewr",
	}
	svcAnnotations := map[string]string{
		annoVultrLoadBalancerID:    "6334f227-6d96-4cbd-9bcb-5be0759354fa",
		annoVultrLoadBalancerLabel: "shared-load-balancer",
	}
	for key, value := range annotations {
		svcAnnotations[key] = value
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "shared-service-b",
			Namespace:   v1.NamespaceDefault,
			UID:         "shared-service-b",
			Annotations: svcAnnotations,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "service-port",
					Port:     port,
					NodePort: nodePort,
				},
			},
		},
	}

	err := lb.UpdateLoadBalancer(context.Background(), "cluster-name", svc, nil)
	if err != nil {
		t.Fatalf("expected nil got %s", err.Error())
	}

	if fakeLoadBalancer.updatedReq == nil {
		t.Fatal("expected load balancer update request")
	}
	if fakeLoadBalancer.updatedReq.ForwardingRules != nil {
		t.Fatalf("expected shared load balancer update to omit forwarding rules, got %+v", fakeLoadBalancer.updatedReq.ForwardingRules)
	}
	if len(fakeLoadBalancer.deletedRules) != 0 {
		t.Fatalf("expected no forwarding rule deletions, got %+v", fakeLoadBalancer.deletedRules)
	}
	if len(fakeLoadBalancer.createdRules) != 1 {
		t.Fatalf("expected one created forwarding rule, got %+v", fakeLoadBalancer.createdRules)
	}

	createdRule := fakeLoadBalancer.createdRules[0]
	if !forwardingRulesEqual(createdRule, expectedRule) {
		t.Fatalf("unexpected created forwarding rule: %+v", createdRule)
	}
}

func TestLoadbalancers_EnsureLoadBalancerDeleted_SharedLabelRemovesOnlyServiceRules(t *testing.T) {
	fakeLoadBalancer := &fakeLB{
		forwardingRules: []govultr.ForwardingRule{
			{
				RuleID:           "rule-50001",
				FrontendProtocol: protocolUDP,
				FrontendPort:     50001,
				BackendProtocol:  protocolUDP,
				BackendPort:      30001,
			},
			{
				RuleID:           "rule-50002",
				FrontendProtocol: protocolUDP,
				FrontendPort:     50002,
				BackendProtocol:  protocolUDP,
				BackendPort:      30002,
			},
		},
	}
	lb := &loadbalancers{
		client:     &govultr.Client{LoadBalancer: fakeLoadBalancer},
		zone:       "ewr",
		kubeClient: fake.NewClientset(sharedLabelService("shared-service-b", "shared-service-b", 50002, 30002)),
	}

	deletingService := sharedLabelService("shared-service-a", "shared-service-a", 50001, 30001)
	err := lb.EnsureLoadBalancerDeleted(context.Background(), "cluster-name", deletingService)
	if err != nil {
		t.Fatalf("expected nil got %s", err.Error())
	}

	if fakeLoadBalancer.deletedLB {
		t.Fatal("expected shared load balancer to remain while another service references its label")
	}
	if !reflect.DeepEqual(fakeLoadBalancer.deletedRules, []string{"rule-50001"}) {
		t.Fatalf("expected only deleting service rule to be removed, got %+v", fakeLoadBalancer.deletedRules)
	}
}

func TestLoadbalancers_EnsureLoadBalancerDeleted_SharedLabelDeletesLBWhenLastReference(t *testing.T) {
	fakeLoadBalancer := &fakeLB{}
	lb := &loadbalancers{
		client:     &govultr.Client{LoadBalancer: fakeLoadBalancer},
		zone:       "ewr",
		kubeClient: fake.NewClientset(),
	}

	deletingService := sharedLabelService("shared-service-a", "shared-service-a", 50001, 30001)
	err := lb.EnsureLoadBalancerDeleted(context.Background(), "cluster-name", deletingService)
	if err != nil {
		t.Fatalf("expected nil got %s", err.Error())
	}

	if !fakeLoadBalancer.deletedLB {
		t.Fatal("expected shared load balancer to be deleted after last service reference is removed")
	}
	if len(fakeLoadBalancer.deletedRules) != 0 {
		t.Fatalf("expected no forwarding rule deletions when deleting the load balancer, got %+v", fakeLoadBalancer.deletedRules)
	}
}

func sharedLabelService(name, uid string, port, nodePort int32) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: v1.NamespaceDefault,
			UID:       typesUID(uid),
			Annotations: map[string]string{
				annoVultrLoadBalancerID:    "6334f227-6d96-4cbd-9bcb-5be0759354fa",
				annoVultrLoadBalancerLabel: "shared-load-balancer",
				annoVultrLBProtocol:        protocolUDP,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "service-port",
					Protocol: v1.ProtocolUDP,
					Port:     port,
					NodePort: nodePort,
				},
			},
		},
	}
}

func typesUID(uid string) types.UID {
	return types.UID(uid)
}
