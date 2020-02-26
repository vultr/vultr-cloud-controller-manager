package vultr

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/vultr/govultr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
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

	annoVultrStickySessionEnabled    = "service.beta.kubernetes.io/vultr-loadbalancer-sticky-session-enabled"
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
	return getDefaultLBName(service)
}

func getDefaultLBName(service *v1.Service) string {
	return cloudprovider.DefaultLoadBalancerName(service)
}

func (l *loadbalancers) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {

	_, exists, err := l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}

	// if exists is false and the err above was nil then this is errLbNotFound
	if !exists {
		// create the LB
		// return from here

	}

	panic("implement me")
}

func (l *loadbalancers) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	panic("implement me")
}

func (l *loadbalancers) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	_, exists, err := l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	lbName := l.GetLoadBalancerName(ctx, clusterName, service)

	lb, err := l.lbByName(ctx, lbName)
	if err != nil {
		return err
	}

	err = l.client.LoadBalancer.Delete(ctx, lb.ID)
	if err != nil {
		return err
	}

	return nil
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

func (l *loadbalancers) buildLoadBalancerRequest(ctx context.Context, service *v1.Service, nodes []*v1.Node) (*govultr.LBConfig, error) {

	//lbName := getDefaultLBName(service)
	// make each section of the LB and add it part fo a global at the bottom

	genericInfo, err := buildGenericInfo(service)
	if err != nil {
		return nil, err
	}

	return &govultr.LBConfig{
		GenericInfo:     *genericInfo,
		HealthCheck:     govultr.HealthCheck{},
		SSLInfo:         false,
		ForwardingRules: govultr.ForwardingRules{},
		InstanceList:    govultr.InstanceList{},
	}, nil
}

func buildGenericInfo(service *v1.Service) (*govultr.GenericInfo, error) {
	// balancing algorithm
	algo := getAlgorithm(service)

	// ssl redirect
	redirect := getSSLRedirect(service)

	// stickSession
	stickySession, err := buildStickySession(service)
	if err != nil {
		return nil, err
	}
	return &govultr.GenericInfo{
		BalancingAlgorithm: algo,
		SSLRedirect:        &redirect,
		StickySessions:     stickySession,
	}, nil
}

// getAlgorithm returns the algorithm to be used for load balancer service
// defaults to round_robin if no algorithm is provided.
func getAlgorithm(service *v1.Service) string {
	algorithm := service.Annotations[annoVultrAlgorithm]

	if algorithm == "least_connections" {
		return "least_connections"
	} else {
		return "round_robin"
	}
}

// getSSLRedirect returns if traffic should be redirected to https
// default to false if not specified
func getSSLRedirect(service *v1.Service) bool {
	redirect, ok := service.Annotations[annoVultrSSLRedirect]
	if !ok {
		return false
	}

	redirectBool, err := strconv.ParseBool(redirect)
	if err != nil {
		return false
	}

	return redirectBool
}

func buildStickySession(service *v1.Service) (*govultr.StickySessions, error) {

	enabled := getStickySessionEnabled(service)

	if enabled == "off" {
		return &govultr.StickySessions{
			StickySessionsEnabled: "off",
		}, nil
	}

	cookieName, err := getCookieName(service)
	if err != nil {
		return nil, err
	}

	return &govultr.StickySessions{
		StickySessionsEnabled: enabled,
		CookieName:            cookieName,
	}, nil
}

// getStickySessionEnabled returns whether or not sticky sessions should be enabled
// default is off
func getStickySessionEnabled(service *v1.Service) string {
	enabled, ok := service.Annotations[annoVultrStickySessionEnabled]
	if !ok {
		return "off"
	}

	if enabled == "off" {
		return "off"
	} else if enabled == "on" {
		return "on"
	} else {
		return "off"
	}
}

// getCookieName returns the cookie name for loadbalancer sticky sessions
func getCookieName(service *v1.Service) (string, error) {
	name, ok := service.Annotations[annoVultrStickySessionCookieName]
	if !ok || name == "" {
		return "", errors.New("sticky session cookie name name not supplied but is required")
	}

	return name, nil
}
