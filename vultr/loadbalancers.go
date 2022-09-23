package vultr

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/vultr/govultr/v2"
	"github.com/vultr/metadata"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

const (
	annoVultrLoadBalancerID = "kubernetes.vultr.com/load-balancer-id"

	// annoVultrLBProtocol is the annotation used to specify
	// which protocol should be used for a Load Balancer.
	// Note that if annoVultrLBHTTPSPorts is defined then this will be overridden to HTTPS
	annoVultrLBProtocol = "service.beta.kubernetes.io/vultr-loadbalancer-protocol"

	// annoVultrLBHTTPSPorts is the annotation used to specify
	// which ports should be used for HTTPS.
	// You can pass in a comma seperated list: 443,8443
	annoVultrLbHttpsPorts = "service.beta.kubernetes.io/vultr-loadbalancer-https-ports"

	// annoVultrLBSSLPassthrough is the annotation used to specify
	// whether or not you do not wish to have SSL termination on the load balancer
	// default is false to enable set to true

	annoVultrLBSSLPassthrough = "service.beta.kubernetes.io/vultr-loadbalancer-ssl-pass-through"

	// annoVultrLBSSL is the annotation used to specify
	// which TLS secret you want to be used for your load balancers SSL
	annoVultrLBSSL = "service.beta.kubernetes.io/vultr-loadbalancer-ssl"

	// annoVultrLBBackendProtocol backend protocol
	annoVultrLBBackendProtocol = "service.beta.kubernetes.io/vultr-loadbalancer-backend-protocol"

	annoVultrHealthCheckPath               = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-path"
	annoVultrHealthCheckProtocol           = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-protocol"
	annoVultrHealthCheckPort               = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-port"
	annoVultrHealthCheckInterval           = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-check-interval"
	annoVultrHealthCheckResponseTimeout    = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-response-timeout"
	annoVultrHealthCheckUnhealthyThreshold = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-unhealthy-threshold"
	annoVultrHealthCheckHealthyThreshold   = "service.beta.kubernetes.io/vultr-loadbalancer-healthcheck-healthy-threshold"

	annoVultrAlgorithm     = "service.beta.kubernetes.io/vultr-loadbalancer-algorithm"
	annoVultrSSLRedirect   = "service.beta.kubernetes.io/vultr-loadbalancer-ssl-redirect"
	annoVultrProxyProtocol = "service.beta.kubernetes.io/vultr-loadbalancer-proxy-protocol"

	annoVultrStickySessionEnabled    = "service.beta.kubernetes.io/vultr-loadbalancer-sticky-session-enabled"
	annoVultrStickySessionCookieName = "service.beta.kubernetes.io/vultr-loadbalancer-sticky-session-cookie-name"

	annoVultrFirewallRules  = "service.beta.kubernetes.io/vultr-loadbalancer-firewall-rules"
	annoVultrPrivateNetwork = "service.beta.kubernetes.io/vultr-loadbalancer-private-network"
	annoVultrVPC            = "service.beta.kubernetes.io/vultr-loadbalancer-vpc"

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

	lbStatusActive = "active"
)

var errLbNotFound = errors.New("loadbalancer not found")
var _ cloudprovider.LoadBalancer = &loadbalancers{}

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
				IP:       lb.IPV4,
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
		lbReq, err := l.buildLoadBalancerRequest(service, nodes)
		if err != nil {
			return nil, err
		}

		lbReq.Region = l.zone
		l, err := l.client.LoadBalancer.Create(ctx, lbReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create load-balancer: %s", err)
		}

		if l.Status != lbStatusActive {
			return nil, fmt.Errorf("load-balancer is not yet active - current status: %s", l.Status)
		}

		return &v1.LoadBalancerStatus{
			Ingress: []v1.LoadBalancerIngress{
				{
					IP:       l.IPV4,
					Hostname: l.Label,
				},
			},
		}, nil
	}

	lbName := l.GetLoadBalancerName(ctx, clusterName, service)
	lb, err := l.lbByName(ctx, lbName)
	if err != nil {
		if err == errLbNotFound {
			return nil, errLbNotFound
		}

		return nil, err
	}

	if lb.Status != lbStatusActive {
		return nil, fmt.Errorf("load-balancer is not yet active - current status: %s", lb.Status)
	}

	if err := l.UpdateLoadBalancer(ctx, clusterName, service, nodes); err != nil {
		return nil, err
	}

	lbStatus, _, err := l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}

	return lbStatus, nil
}

func (l *loadbalancers) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	klog.V(3).Info("Called UpdateLoadBalancers")
	if _, _, err := l.GetLoadBalancer(ctx, clusterName, service); err != nil {
		return err
	}

	lbName := l.GetLoadBalancerName(ctx, clusterName, service)
	lb, err := l.lbByName(ctx, lbName)
	if err != nil {
		return err
	}

	lbReq, err := l.buildLoadBalancerRequest(service, nodes)
	if err != nil {
		return fmt.Errorf("failed to create load balancer request: %s", err)
	}

	if err := l.client.LoadBalancer.Update(ctx, lb.ID, lbReq); err != nil {
		return fmt.Errorf("failed to update LB: %s", err)
	}

	return nil
}

func (l *loadbalancers) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	_, exists, err := l.GetLoadBalancer(ctx, clusterName, service)
	if err != nil {
		return err
	}
	// This is the same as if we were to check if err == errLbNotFound {
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

func (l *loadbalancers) lbByName(ctx context.Context, lbName string) (*govultr.LoadBalancer, error) {

	listOptions := &govultr.ListOptions{
		PerPage: 25,
	}

	for {
		lbs, meta, err := l.client.LoadBalancer.List(ctx, listOptions)
		if err != nil {
			return nil, err
		}

		for _, v := range lbs {
			if v.Label == lbName {
				return &v, nil
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			listOptions.Cursor = meta.Links.Next
		}
	}

	return nil, errLbNotFound
}

func (l *loadbalancers) buildLoadBalancerRequest(service *v1.Service, nodes []*v1.Node) (*govultr.LoadBalancerReq, error) {

	stickySession, err := buildStickySession(service)
	if err != nil {
		return nil, err
	}

	healthCheck, err := buildHealthChecks(service)
	if err != nil {
		return nil, err
	}

	rules, err := buildForwardingRules(service)
	if err != nil {
		return nil, err
	}

	instances, err := buildInstanceList(nodes)
	if err != nil {
		return nil, err
	}

	ssl := &govultr.SSL{}
	if secretName, ok := service.Annotations[annoVultrLBSSL]; ok {
		ssl, err = l.GetSSL(service, secretName)
		if err != nil {
			return nil, err
		}
	} else {
		ssl = nil
	}

	firewallRules, err := buildFirewallRules(service)
	if err != nil {
		return nil, err
	}
	vpc, err := getVPC(service)
	if err != nil {
		return nil, err
	}

	return &govultr.LoadBalancerReq{
		Label:              getDefaultLBName(service),                        // will always be set
		Instances:          instances,                                        // will always be set
		HealthCheck:        healthCheck,                                      // will always be set
		StickySessions:     stickySession,                                    // need to check
		ForwardingRules:    rules,                                            // all always be set
		SSL:                ssl,                                              // will always be set
		SSLRedirect:        govultr.BoolToBoolPtr(getSSLRedirect(service)),   // need to check
		ProxyProtocol:      govultr.BoolToBoolPtr(getProxyProtocol(service)), // need to check
		BalancingAlgorithm: getAlgorithm(service),                            // will always be set
		FirewallRules:      firewallRules,                                    // need to check
		VPC:                govultr.StringToStringPtr(vpc),                   // need to check
	}, nil
}

// getAlgorithm returns the algorithm to be used for load balancer service
// defaults to round_robin if no algorithm is provided.
func getAlgorithm(service *v1.Service) string {
	algo := service.Annotations[annoVultrAlgorithm]
	if algo == "least_connections" {
		return "leastconn"
	} else {
		return "roundrobin"
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
			CookieName: "",
		}, nil
	}

	cookieName, err := getCookieName(service)
	if err != nil {
		return nil, err
	}

	return &govultr.StickySessions{
		CookieName: cookieName,
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
	if protocol == "" {
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

// buildInstanceList create list of nodes to be attached to a load balancer
func buildInstanceList(nodes []*v1.Node) ([]string, error) {
	var list []string

	for _, node := range nodes {
		instanceID, err := vultrIDFromProviderID(node.Spec.ProviderID)
		if err != nil {
			return nil, fmt.Errorf("error getting the provider ID %s : %s", node.Spec.ProviderID, err)
		}

		list = append(list, instanceID)
	}

	return list, nil
}

func buildForwardingRules(service *v1.Service) ([]govultr.ForwardingRule, error) {
	var rules []govultr.ForwardingRule

	defaultProtocol := getLBProtocol(service)

	httpsPorts, err := getHttpsPorts(service)
	if err != nil {
		return nil, err
	}

	for _, port := range service.Spec.Ports {
		// default the port
		protocol := defaultProtocol
		backendProtocol := getBackendProtocol(service)

		if httpsPorts[port.Port] {
			if getSSLPassthrough(service) {
				protocol = protocolTCP
			} else {
				protocol = protocolHTTPs
			}
		}

		if backendProtocol == "" {
			backendProtocol = protocol
		}

		rule, err := buildForwardingRule(&port, protocol, backendProtocol)
		if err != nil {
			return nil, err
		}

		rules = append(rules, *rule)
	}

	return rules, nil
}

func buildForwardingRule(port *v1.ServicePort, protocol, backendProtocol string) (*govultr.ForwardingRule, error) {
	var rule govultr.ForwardingRule

	if port.Protocol == portProtocolUDP {
		return nil, fmt.Errorf("TCP protocol is only supported: recieved %s", port.Protocol)
	}

	rule.FrontendProtocol = protocol
	rule.BackendProtocol = backendProtocol

	rule.FrontendPort = int(port.Port)
	rule.BackendPort = int(port.NodePort)

	return &rule, nil
}

func getLBProtocol(service *v1.Service) string {
	protocol, ok := service.Annotations[annoVultrLBProtocol]
	if !ok {
		return protocolTCP
	}

	return protocol
}

func getHttpsPorts(service *v1.Service) (map[int32]bool, error) {
	ports, ok := service.Annotations[annoVultrLbHttpsPorts]
	if !ok {
		return nil, nil
	}

	portStrings := strings.Split(ports, ",")
	portInt := map[int32]bool{}

	for _, port := range portStrings {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		portInt[int32(p)] = true
	}
	return portInt, nil
}

func (l *loadbalancers) GetSSL(service *v1.Service, secretName string) (*govultr.SSL, error) {
	if err := l.GetKubeClient(); err != nil {
		return nil, err
	}

	secret, err := l.kubeClient.CoreV1().Secrets(service.Namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	cert := string(secret.Data[v1.TLSCertKey])
	cert = strings.TrimSpace(cert)

	key := string(secret.Data[v1.TLSPrivateKeyKey])
	key = strings.TrimSpace(key)

	ssl := govultr.SSL{
		PrivateKey:  key,
		Certificate: cert,
	}
	return &ssl, nil
}

func (l *loadbalancers) GetKubeClient() error {
	if l.kubeClient != nil {
		return nil
	}

	var (
		kubeConfig *rest.Config
		err        error
		config     string
	)

	// If no kubeconfig was passed in or set then we want to default to an empty string
	// This will have `clientcmd.BuildConfigFromFlags` default to `restclient.InClusterConfig()` which was existing behavior
	if Options.KubeconfigFlag == nil || Options.KubeconfigFlag.Value.String() == "" {
		config = ""
	} else {
		config = Options.KubeconfigFlag.Value.String()
	}

	kubeConfig, err = clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return err
	}

	l.kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

func getSSLPassthrough(service *v1.Service) bool {
	passThrough, ok := service.Annotations[annoVultrLBSSLPassthrough]
	if !ok {
		return false
	}

	pass, err := strconv.ParseBool(passThrough)
	if err != nil {
		return false
	}
	return pass
}

func getProxyProtocol(service *v1.Service) bool {
	proxy, ok := service.Annotations[annoVultrProxyProtocol]
	if !ok {
		return false
	}

	pass, err := strconv.ParseBool(proxy)
	if err != nil {
		return false
	}

	return pass
}

func buildFirewallRules(service *v1.Service) ([]govultr.LBFirewallRule, error) {
	lbFWRules := []govultr.LBFirewallRule{}
	fwRules := getFirewallRules(service)
	if fwRules == "" {
		return lbFWRules, nil
	}

	for _, v := range strings.Split(fwRules, ";") {
		fwRule := govultr.LBFirewallRule{}

		rules := strings.Split(v, ",")
		if len(rules) != 2 {
			return nil, fmt.Errorf("loadbalancer fw rules : %s invalid configuration", rules)
		}
		source := rules[0]
		_, _, err := net.ParseCIDR(source)
		if err != nil {
			return nil, fmt.Errorf("loadbalancer fw rules : source %s is invalid", source)
		}

		port, err := strconv.Atoi(rules[1])
		if err != nil {
			return nil, fmt.Errorf("loadbalancer fw rules : port %d is invalid", port)
		}

		fwRule.Source = source
		fwRule.IPType = "v4"
		fwRule.Port = port
		lbFWRules = append(lbFWRules, fwRule)
	}
	return lbFWRules, nil
}

func getFirewallRules(service *v1.Service) string {
	fwRules, ok := service.Annotations[annoVultrFirewallRules]
	if !ok {
		return ""
	}

	return fwRules
}

func getVPC(service *v1.Service) (string, error) {

	var vpc string
	pn, pnOk := service.Annotations[annoVultrPrivateNetwork]
	v, vpcOk := service.Annotations[annoVultrVPC]

	if vpcOk && pnOk {
		return "", fmt.Errorf("can not use private_network and vpc annotations. Please use VPC as private network is deprecated")
	} else if pnOk {
		vpc = pn
	} else if vpcOk {
		vpc = v
	} else {
		return "", nil
	}

	if strings.ToLower(vpc) == "false" {
		return "", nil
	}

	meta := metadata.NewClient()
	m, err := meta.Metadata()
	if err != nil {
		return "", fmt.Errorf("error getting metadata for private_network: %v", err.Error())
	}

	pnID := ""
	for _, v := range m.Interfaces {
		if v.NetworkV2ID != "" {
			pnID = v.NetworkV2ID
			break
		}
	}

	return pnID, nil
}

func getBackendProtocol(service *v1.Service) string {
	proto, ok := service.Annotations[annoVultrLBBackendProtocol]
	if !ok {
		return ""
	}

	switch proto {
	case "http":
		return "http"
	case "https":
		return "https"
	case "tcp":
		return "tcp"
	default:
		return ""
	}
}
