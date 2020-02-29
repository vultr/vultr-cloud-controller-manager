package vultr

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vultr/govultr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"strconv"
)

const (
	annoVultrLoadBalancerID = "kubernetes.vultr.com/load-balancer-id"
	annoVultrLBProtocol     = "service.beta.kubernetes.io/vultr-loadbalancer-protocol"
	// add in support for LB ports
	// add in support for LB https

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

	// Supported Protocols
	protocolHTTP  = "http"
	protocolHTTPs = "https"
	protocolTCP   = "tcp"

	portProtocolTCP = "TCP"
	portProtocolUDP = "UDP"

	healthCheckInterval  = 15
	healthCheckResponse  = 5
	healthCheckUnhealthy = 5
	healthCheckHealthy   = 5
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
		lb, lbName, err := l.buildLoadBalancerRequest(service, nodes)
		if err != nil {
			return nil, err
		}
		zone, err := strconv.Atoi(l.zone)
		if err != nil {
			return nil, err
		}
		lbID, err := l.client.LoadBalancer.Create(ctx, zone, &lb.GenericInfo, &lb.HealthCheck, lb.ForwardingRules.ForwardRuleList)
		if err != nil {
			return nil, fmt.Errorf("failed to create load-balancer: %s", err)
		}

		err = l.client.LoadBalancer.SetLabel(ctx, lbID.ID, lbName)
		if err != nil {
			return nil, fmt.Errorf("failed to set load-balancer name: %s", err)
		}

		for _, node := range lb.InstanceList.InstanceList {
			n, err := strconv.Atoi(node)
			if err != nil {
				return nil, err
			}
			err = l.client.LoadBalancer.AttachInstance(ctx, lbID.ID, n)
			if err != nil {
				return nil, fmt.Errorf("failed attach nodes to lb %s", err)
			}
		}


		list, _ := l.client.LoadBalancer.List(ctx)
		var l govultr.LoadBalancers
		for _, v := range list {
			if v.ID == lbID.ID {
				l = v
			}
		}

		return &v1.LoadBalancerStatus{
			Ingress: []v1.LoadBalancerIngress{
				{
					IP: l.IPV4,
					Hostname: l.Label,
				},
			},
		}, nil
	}

	// update LB
	lbStatus,_, err := l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}

	return lbStatus, nil
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

func (l *loadbalancers) buildLoadBalancerRequest(service *v1.Service, nodes []*v1.Node) (*govultr.LBConfig, string, error) {

	lbName := getDefaultLBName(service)

	genericInfo, err := buildGenericInfo(service)
	if err != nil {
		return nil, "", err
	}

	healthCheck, err := buildHealthChecks(service)
	if err != nil {
		return nil, "", err
	}

	rules, err := buildForwardingRules(service)
	if err != nil {
		return nil, "", err
	}

	instances, err := buildInstanceList(nodes)
	if err != nil {
		return nil, "", err
	}

	return &govultr.LBConfig{
		GenericInfo:     *genericInfo,
		HealthCheck:     *healthCheck,
		SSLInfo:         false,
		ForwardingRules: *rules,
		InstanceList:    *instances,
	}, lbName, nil
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

func buildHealthChecks(service *v1.Service) (*govultr.HealthCheck, error) {

	healthCheckProtocol, err := getHealthCheckProtocol(service)
	if err != nil {
		return nil, err
	}

	port, err := getHealthCheckPort(service)
	if err != nil {
		return nil, err
	}

	path := getHealthCheckPath(service)

	interval, err := getHealthCheckInterval(service)
	if err != nil {
		return nil, err
	}

	response, err := getHealthCheckResponse(service)
	if err != nil {
		return nil, err
	}

	unhealthy, err := getHealthCheckUnhealthy(service)
	if err != nil {
		return nil, err
	}

	healthy, err := getHealthCheckHealthy(service)
	if err != nil {
		return nil, err
	}

	return &govultr.HealthCheck{
		Protocol:           healthCheckProtocol,
		Port:               port,
		Path:               path,
		CheckInterval:      interval,
		ResponseTimeout:    response,
		UnhealthyThreshold: unhealthy,
		HealthyThreshold:   healthy,
	}, nil
}

// getHealthCheckProtocol returns the protocol for a health check
// default is TCP
func getHealthCheckProtocol(service *v1.Service) (string, error) {
	protocol := service.Annotations[annoVultrHealthCheckProtocol]

	// add in https
	if protocol != "" {
		if getHealthCheckPath(service) != "" {
			return protocolHTTP, nil
		}
		return protocolTCP, nil
	}

	if protocol != protocolHTTP && protocol != protocolTCP {
		return "", fmt.Errorf("invalid protocol : %s given in the anootation : %s", protocol, annoVultrHealthCheckProtocol)
	}

	return protocol, nil
}

// getHealthCheckPath returns the path for a health check
func getHealthCheckPath(service *v1.Service) string {
	path, ok := service.Annotations[annoVultrHealthCheckPath]
	if !ok {
		return ""
	}

	return path
}

func getHealthCheckPort(service *v1.Service) (int, error) {
	port, ok := service.Annotations[annoVultrHealthCheckPort]
	if !ok {
		return int(service.Spec.Ports[0].NodePort), nil
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}

	for _, v := range service.Spec.Ports {
		if int(v.Port) == portInt {
			return int(v.Port), nil
		}
		// The provided port does not exist
		return 0, fmt.Errorf("provided health check port %d does not exist for service %s/%s", portInt, service.Namespace, service.Name)
	}

	// need to default a return here
	return 0, nil
}

func getHealthCheckInterval(service *v1.Service) (int, error) {
	interval, ok := service.Annotations[annoVultrHealthCheckInterval]
	if !ok {
		return healthCheckInterval, nil
	}

	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		return 0, fmt.Errorf("failed to retireve health check interval %s - %s", annoVultrHealthCheckInterval, err)
	}

	return intervalInt, err
}

func getHealthCheckResponse(service *v1.Service) (int, error) {
	response, ok := service.Annotations[annoVultrHealthCheckResponseTimeout]
	if !ok {
		return healthCheckResponse, nil
	}

	responseInt, err := strconv.Atoi(response)
	if err != nil {
		return 0, fmt.Errorf("failed to retireve health check response timeout %s - %s", annoVultrHealthCheckResponseTimeout, err)
	}

	return responseInt, err
}

func getHealthCheckUnhealthy(service *v1.Service) (int, error) {
	unhealthy, ok := service.Annotations[annoVultrHealthCheckUnhealthyThreshold]
	if !ok {
		return healthCheckUnhealthy, nil
	}

	unhealthyInt, err := strconv.Atoi(unhealthy)
	if err != nil {
		return 0, fmt.Errorf("failed to retireve health check unhealthy treshold %s - %s", annoVultrHealthCheckUnhealthyThreshold, err)
	}

	return unhealthyInt, err
}

func getHealthCheckHealthy(service *v1.Service) (int, error) {
	healthy, ok := service.Annotations[annoVultrHealthCheckHealthyThreshold]
	if !ok {
		return healthCheckHealthy, nil
	}

	healthyInt, err := strconv.Atoi(healthy)
	if err != nil {
		return 0, fmt.Errorf("failed to retireve health check healthy treshold %s - %s", annoVultrHealthCheckHealthyThreshold, err)
	}

	return healthyInt, err
}

func buildInstanceList(nodes []*v1.Node) (*govultr.InstanceList, error) {
	var list []string

	for _, node := range nodes {
		instanceID, err := vultrIDFromProviderID(node.Spec.ProviderID)
		if err != nil {
			return nil, fmt.Errorf("error getting the provider ID %s : %s", node.Spec.ProviderID, err)
		}

		list = append(list, instanceID)
	}

	return &govultr.InstanceList{InstanceList: list}, nil
}

func buildForwardingRules(service *v1.Service) (*govultr.ForwardingRules, error) {
	var rules govultr.ForwardingRules

	// flush this out with HTTPS and better handling
	protocol := getLBProtocol(service)

	for _, port := range service.Spec.Ports {
		rule, err := buildForwardingRule(&port, protocol)
		if err != nil {
			return nil, err
		}

		rules.ForwardRuleList = append(rules.ForwardRuleList, *rule)
	}

	return &rules, nil
}

func buildForwardingRule(port *v1.ServicePort, protocol string) (*govultr.ForwardingRule, error) {
	var rule govultr.ForwardingRule

	if port.Protocol == portProtocolUDP {
		return nil, fmt.Errorf("TCP protocol is only supported: recieved %s", port.Protocol)
	}

	rule.FrontendProtocol = protocol
	rule.BackendProtocol = protocol

	rule.FrontendPort = int(port.Port)
	rule.BackendPort = int(port.NodePort)

	return &rule, nil
}

func getLBProtocol(service *v1.Service) string {
	protocol, ok := service.Annotations[annoVultrLBProtocol]
	if !ok {
		return protocolHTTP
	}

	return protocol
}
