package vultr

import (
	"context"
	"github.com/pkg/errors"
	"github.com/vultr/govultr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"strconv"
	"strings"
)

const (
	annoVultrLoadBalancerID = "kubernetes.vultr.com/load-balancer-id"

	annoVultrProtocol = "service.beta.kubernetes.io/vultr-loadbalancer-protocol"

	annoVultrHealthCheckPath               = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-path"
	annoVultrHealthCheckProtocol           = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-protocol"
	annoVultrHealthCheckPort               = "service.beta.kubernetes.io/vultr-loadbalancer-port-healthcheck-port"
	annoVultrHealthCheckInterval           = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-check-interval"
	annoVultrHealthCheckResponseTimeout    = "service.beta.kubernetes.io/vultr-loadbalancer-port-healthcheck-response-timeout"
	annoVultrHealthCheckUnhealthyThreshold = "service.beta.kubernetes.io/vultr-loadbalancer-port-healthcheck-unhealthy-threshold"
	annoVultrHealthCheckHealthyThreshold   = "service.beta.kubernetes.io/vultr-loadbalancer-port-healthcheck-healthy-threshold"

	annoVultrAlgorithm     = "service.beta.kubernetes.io/vultr-loadbalancer-algorithm"
	annoVultrSSLRedirect   = "service.beta.kubernetes.io/vultr-loadbalancer-ssl-redirect"
	annoVultrProxyProtocol = "service.beta.kubernetes.io/vultr-loadbalancer-proxy-protocol"

	annoVultrStickySessionCookieName = "service.beta.kubernetes.io/vultr-loadbalancer-sticky-session-cookie-name"
)

var errLbNotFound = errors.New("loadbalancer not found")

type loadbalancers struct {
	client *govultr.Client
	zone   string

	kubeClient kubernetes.Interface
}

func newLoadbalancers(client *govultr.Client, zone string) cloudprovider.LoadBalancer {
	return &loadbalancers{client: client, zone: zone}
}

func (l *loadbalancers) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	lbName := l.GetLoadBalancerName(ctx, clusterName, service)

	lb, err := l.lbByName(ctx, lbName)
	if err != nil {
		if err == errLbNotFound {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP:       strconv.Itoa(lb.ID),
				Hostname: lb.Label,
			},
		},
	}, true, nil
}

func (l *loadbalancers) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	ret := clusterName + string(service.UID)
	ret = strings.Replace(ret, "-", "", -1)
	return ret
}

func (l *loadbalancers) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	panic("implement me")
}

func (l *loadbalancers) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	panic("implement me")
}

func (l *loadbalancers) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	panic("implement me")
}

func getLoadBalancerID(service *v1.Service) string {
	return service.ObjectMeta.Annotations[annoVultrLoadBalancerID]
}

func (l *loadbalancers) lbByName(ctx context.Context, lbName string) (*govultr.LoadBalancers, error) {
	lbs, err := l.client.LoadBalancer.List(ctx)
	if err != nil {
		return nil, err
	}

	// go through the list and find the matching LB
	if len(lbs) > 0 {
		for _, v := range lbs {
			if v.Label == lbName {
				return &v, nil
				// grab the full config
				//lb, err := l.client.LoadBalancer.GetFullConfig(ctx, v.ID)
				//if err != nil {
				//	return nil, err
				//}
				//
				//return lb, nil
			}
		}
	}

	return nil, errLbNotFound
}
