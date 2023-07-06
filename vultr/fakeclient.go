// Package vultr is vultr cloud specific implementation
package vultr

import (
	"context"

	"github.com/vultr/govultr/v2"
)

func newFakeClient() *govultr.Client {
	fakeInstance := FakeInstance{client: nil}
	fakeLoadBalancer := fakeLB{client: nil}
	return &govultr.Client{
		Instance:     &fakeInstance,
		LoadBalancer: &fakeLoadBalancer,
	}
}

// FakeInstance creates a fake instance to govultr.Client
type FakeInstance struct {
	client *govultr.Client
}

// ListVPCInfo returns VPC info (not implemented, yet)
func (f *FakeInstance) ListVPCInfo(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.VPCInfo, *govultr.Meta, error) {
	panic("implement me")
}

// AttachVPC attaches VPC (not implemented, yet)
func (f *FakeInstance) AttachVPC(_ context.Context, _, _ string) error {
	panic("implement me")
}

// DetachVPC detaches VPC (not implemented, yet)
func (f *FakeInstance) DetachVPC(_ context.Context, _, _ string) error {
	panic("implement me")
}

// Create creates an instance (not implemented, yet)
func (f *FakeInstance) Create(_ context.Context, _ *govultr.InstanceCreateReq) (*govultr.Instance, error) {
	panic("implement me")
}

// Get returns instance
func (f *FakeInstance) Get(_ context.Context, _ string) (*govultr.Instance, error) {
	return &govultr.Instance{
		ID:           "75b95d83-47e2-4c0f-b273-cc9ce2b456f8",
		MainIP:       "149.28.225.110",
		VCPUCount:    4,
		Region:       "ewr",
		Status:       "running",
		NetmaskV4:    "255.255.254.0",
		GatewayV4:    "149.28.224.1",
		ServerStatus: "",
		Plan:         "vc2-4c-8gb",
		Label:        "ccm-test",
		InternalIP:   "10.1.95.4",
	}, nil
}

// Update updates and instance
func (f *FakeInstance) Update(_ context.Context, _ string, _ *govultr.InstanceUpdateReq) (*govultr.Instance, error) {
	panic("implement me")
}

// Delete deletes an instance
func (f *FakeInstance) Delete(_ context.Context, _ string) error {
	panic("implement me")
}

// List lists instances
func (f *FakeInstance) List(_ context.Context, _ *govultr.ListOptions) ([]govultr.Instance, *govultr.Meta, error) {
	return []govultr.Instance{
			{
				ID:           "75b95d83-47e2-4c0f-b273-cc9ce2b456f8",
				MainIP:       "149.28.225.110",
				VCPUCount:    4,
				Region:       "ewr",
				Status:       "running",
				NetmaskV4:    "255.255.254.0",
				GatewayV4:    "149.28.224.1",
				ServerStatus: "",
				Plan:         "vc2-4c-8gb",
				Label:        "ccm-test",
				InternalIP:   "10.1.95.4",
			},
		}, &govultr.Meta{
			Total: 0,
			Links: &govultr.Links{
				Next: "",
				Prev: "",
			},
		}, nil
}

// Start starts an instance
func (f *FakeInstance) Start(_ context.Context, _ string) error {
	panic("implement me")
}

// Halt halts/stops an instance
func (f *FakeInstance) Halt(_ context.Context, _ string) error {
	panic("implement me")
}

// Reboot reboots an instance
func (f *FakeInstance) Reboot(_ context.Context, _ string) error {
	panic("implement me")
}

// Reinstall reinstalls an instance
func (f *FakeInstance) Reinstall(_ context.Context, _ string, _ *govultr.ReinstallReq) (*govultr.Instance, error) {
	panic("implement me")
}

// MassStart bulk starts instances
func (f *FakeInstance) MassStart(_ context.Context, _ []string) error {
	panic("implement me")
}

// MassHalt bulk stops instances
func (f *FakeInstance) MassHalt(_ context.Context, _ []string) error {
	panic("implement me")
}

// MassReboot bulk reboots instances
func (f *FakeInstance) MassReboot(_ context.Context, _ []string) error {
	panic("implement me")
}

// Restore restores an instance
func (f *FakeInstance) Restore(_ context.Context, _ string, _ *govultr.RestoreReq) error {
	panic("implement me")
}

// GetBandwidth gets bandwidth for an instance
func (f *FakeInstance) GetBandwidth(_ context.Context, _ string) (*govultr.Bandwidth, error) {
	panic("implement me")
}

// GetNeighbors gets neighors for an instance
func (f *FakeInstance) GetNeighbors(_ context.Context, _ string) (*govultr.Neighbors, error) {
	panic("implement me")
}

// ListPrivateNetworks gets private networks
func (f *FakeInstance) ListPrivateNetworks(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.PrivateNetwork, *govultr.Meta, error) {
	panic("implement me")
}

// AttachPrivateNetwork attches private networks
func (f *FakeInstance) AttachPrivateNetwork(_ context.Context, _, _ string) error {
	panic("implement me")
}

// DetachPrivateNetwork detaches private network from instance
func (f *FakeInstance) DetachPrivateNetwork(_ context.Context, _, _ string) error {
	panic("implement me")
}

// ISOStatus gets ISO status from instance
func (f *FakeInstance) ISOStatus(_ context.Context, _ string) (*govultr.Iso, error) {
	panic("implement me")
}

// AttachISO attaches ISO to instance
func (f *FakeInstance) AttachISO(_ context.Context, _, _ string) error {
	panic("implement me")
}

// DetachISO detaches ISO from instance
func (f *FakeInstance) DetachISO(_ context.Context, _ string) error {
	panic("implement me")
}

// GetBackupSchedule gets instance backup stchedule
func (f *FakeInstance) GetBackupSchedule(_ context.Context, _ string) (*govultr.BackupSchedule, error) {
	panic("implement me")
}

// SetBackupSchedule sets instance backup schedule
func (f *FakeInstance) SetBackupSchedule(_ context.Context, _ string, _ *govultr.BackupScheduleReq) error {
	panic("implement me")
}

// CreateIPv4 creates an IPv4 association to instance
func (f *FakeInstance) CreateIPv4(_ context.Context, _ string, _ *bool) (*govultr.IPv4, error) {
	panic("implement me")
}

// ListIPv4 gets IPv4 addresses associated with instance
func (f *FakeInstance) ListIPv4(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.IPv4, *govultr.Meta, error) {
	panic("implement me")
}

// DeleteIPv4 deletes IPv4 address associated with instance
func (f *FakeInstance) DeleteIPv4(_ context.Context, _, _ string) error {
	panic("implement me")
}

// ListIPv6 lists IPv6 addresses associated with instance
func (f *FakeInstance) ListIPv6(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.IPv6, *govultr.Meta, error) {
	panic("implement me")
}

// CreateReverseIPv6 adds reverse IP to IPv6
func (f *FakeInstance) CreateReverseIPv6(_ context.Context, _ string, _ *govultr.ReverseIP) error {
	panic("implement me")
}

// ListReverseIPv6 gets reverse IP for IPv6 on instance
func (f *FakeInstance) ListReverseIPv6(_ context.Context, _ string) ([]govultr.ReverseIP, error) {
	panic("implement me")
}

// DeleteReverseIPv6 deletes IPv6 reverse for instance
func (f *FakeInstance) DeleteReverseIPv6(_ context.Context, _, _ string) error {
	panic("implement me")
}

// CreateReverseIPv4 creates reverse IPv4 for instance
func (f *FakeInstance) CreateReverseIPv4(_ context.Context, _ string, _ *govultr.ReverseIP) error {
	panic("implement me")
}

// DefaultReverseIPv4 sets default for IPv4 on instance
func (f *FakeInstance) DefaultReverseIPv4(_ context.Context, _, _ string) error {
	panic("implement me")
}

// GetUserData returns instance userdata
func (f *FakeInstance) GetUserData(_ context.Context, _ string) (*govultr.UserData, error) {
	panic("implement me")
}

// GetUpgrades gets instance upgade
func (f *FakeInstance) GetUpgrades(_ context.Context, _ string) (*govultr.Upgrades, error) {
	panic("implement me")
}

type fakeLB struct {
	client *govultr.Client
}

// Create creates loadbalancer
func (f *fakeLB) Create(_ context.Context, _ *govultr.LoadBalancerReq) (*govultr.LoadBalancer, error) {
	panic("implement me")
}

// Get gets loadbalancer
func (f *fakeLB) Get(_ context.Context, _ string) (*govultr.LoadBalancer, error) {
	return &govultr.LoadBalancer{
		ID:        "6334f227-6d96-4cbd-9bcb-5be0759354fa",
		Region:    "ewr",
		Label:     "albname",
		Status:    "active",
		IPV4:      "192.168.0.1",
		Instances: []string{"0c51cc3d-529e-4e03-ad86-fd0af47467ba", "ca9a74cb-2d9f-4786-9bb0-094398c593a2"},
	}, nil
}

// Update updates loadbalancer
func (f *fakeLB) Update(_ context.Context, _ string, _ *govultr.LoadBalancerReq) error {
	return nil
}

// Delete deletes loadbalancer
func (f *fakeLB) Delete(_ context.Context, _ string) error {
	panic("implement me")
}

// List gets loadbalancers
func (f *fakeLB) List(_ context.Context, _ *govultr.ListOptions) ([]govultr.LoadBalancer, *govultr.Meta, error) {
	return []govultr.LoadBalancer{
			{
				ID:     "6334f227-6d96-4cbd-9bcb-5be0759354fa",
				Region: "ewr",
				Label:  "albname",
				Status: "active",
				IPV4:   "192.168.0.1",
			},
		}, &govultr.Meta{
			Total: 0,
			Links: &govultr.Links{
				Next: "",
				Prev: "",
			},
		}, nil
}

// CreateForwardingRule adds forwarding rule
func (f *fakeLB) CreateForwardingRule(_ context.Context, _ string, _ *govultr.ForwardingRule) (*govultr.ForwardingRule, error) {
	panic("implement me")
}

// GetForwardingRule returns forwarding rule
func (f *fakeLB) GetForwardingRule(_ context.Context, _, _ string) (*govultr.ForwardingRule, error) {
	panic("implement me")
}

// DeleteForwardingRule deletes forwarding rule
func (f *fakeLB) DeleteForwardingRule(_ context.Context, _, _ string) error {
	panic("implement me")
}

// ListForwardingRules gets forwarding rules
func (f *fakeLB) ListForwardingRules(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.ForwardingRule, *govultr.Meta, error) {
	return []govultr.ForwardingRule{{
			RuleID:           "1234",
			FrontendProtocol: "tcp",
			FrontendPort:     80,
			BackendProtocol:  "tcp",
			BackendPort:      80,
		}}, &govultr.Meta{
			Total: 0,
			Links: &govultr.Links{
				Next: "",
				Prev: "",
			},
		}, nil
}

// ListFirewallRules gets forwarding rules
func (f *fakeLB) ListFirewallRules(_ context.Context, _ string, _ *govultr.ListOptions) ([]govultr.LBFirewallRule, *govultr.Meta, error) {
	return nil, nil, nil
}

// GetFirewallRule gets firewall rules
func (f *fakeLB) GetFirewallRule(_ context.Context, _, _ string) (*govultr.LBFirewallRule, error) {
	return nil, nil
}
