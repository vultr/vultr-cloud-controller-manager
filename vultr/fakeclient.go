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

func newFakeInstance() *govultr.Instance {
	return &govultr.Instance{
		ID:           "5b95d83-47e2-4c0f-b273-cc9ce2b456f8",
		MainIP:       "149.28.225.110",
		VCPUCount:    4,
		Region:       "ewr",
		Status:       "running",
		NetmaskV4:    "255.255.254.0",
		GatewayV4:    "149.28.224.1",
		PowerStatus:  "",
		ServerStatus: "",
		Plan:         "vc2-4c-8gb",
		Label:        "cluster-name",
		InternalIP:   "10.1.95.4",
	}
}

type FakeInstance struct {
	client *govultr.Client
}

func (f *FakeInstance) Create(ctx context.Context, instanceReq *govultr.InstanceCreateReq) (*govultr.Instance, error) {
	panic("implement me")
}

func (f *FakeInstance) Get(ctx context.Context, instanceID string) (*govultr.Instance, error) {
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

func (f *FakeInstance) Update(ctx context.Context, instanceID string, instanceReq *govultr.InstanceUpdateReq) error {
	panic("implement me")
}

func (f *FakeInstance) Delete(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Instance, *govultr.Meta, error) {
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
		},
		nil, nil
}

func (f *FakeInstance) Start(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) Halt(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) Reboot(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) Reinstall(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) MassStart(ctx context.Context, instanceList []string) error {
	panic("implement me")
}

func (f *FakeInstance) MassHalt(ctx context.Context, instanceList []string) error {
	panic("implement me")
}

func (f *FakeInstance) MassReboot(ctx context.Context, instanceList []string) error {
	panic("implement me")
}

func (f *FakeInstance) Restore(ctx context.Context, instanceID string, restoreReq *govultr.RestoreReq) error {
	panic("implement me")
}

func (f *FakeInstance) GetBandwidth(ctx context.Context, instanceID string) (*govultr.Bandwidth, error) {
	panic("implement me")
}

func (f *FakeInstance) GetNeighbors(ctx context.Context, instanceID string) (*govultr.Neighbors, error) {
	panic("implement me")
}

func (f *FakeInstance) ListPrivateNetworks(ctx context.Context, instanceID string) ([]govultr.PrivateNetwork, *govultr.Meta, error) {
	panic("implement me")
}

func (f *FakeInstance) AttachPrivateNetwork(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

func (f *FakeInstance) DetachPrivateNetwork(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

func (f *FakeInstance) ISOStatus(ctx context.Context, instanceID string) (*govultr.Iso, error) {
	panic("implement me")
}

func (f *FakeInstance) AttachISO(ctx context.Context, instanceID, isoID string) error {
	panic("implement me")
}

func (f *FakeInstance) DetachISO(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) GetBackupSchedule(ctx context.Context, instanceID string) (*govultr.BackupSchedule, error) {
	panic("implement me")
}

func (f *FakeInstance) SetBackupSchedule(ctx context.Context, instanceID string, backup *govultr.BackupScheduleReq) error {
	panic("implement me")
}

func (f *FakeInstance) CreateIPv4(ctx context.Context, instanceID string, reboot bool) (*govultr.IPv4, error) {
	panic("implement me")
}

func (f *FakeInstance) ListIPv4(ctx context.Context, instanceID string, option *govultr.ListOptions) ([]govultr.IPv4, *govultr.Meta, error) {
	panic("implement me")
}

func (f *FakeInstance) DeleteIPv4(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

func (f *FakeInstance) ListIPv6(ctx context.Context, instanceID string, option *govultr.ListOptions) ([]govultr.IPv6, *govultr.Meta, error) {
	panic("implement me")
}

func (f *FakeInstance) CreateReverseIPv6(ctx context.Context, instanceID string, reverseReq *govultr.ReverseIP) error {
	panic("implement me")
}

func (f *FakeInstance) ListReverseIPv6(ctx context.Context, instanceID string) ([]govultr.ReverseIP, error) {
	panic("implement me")
}

func (f *FakeInstance) DeleteReverseIPv6(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

func (f *FakeInstance) CreateReverseIPv4(ctx context.Context, instanceID string, reverseReq *govultr.ReverseIP) error {
	panic("implement me")
}

func (f *FakeInstance) DefaultReverseIPv4(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

func (f *FakeInstance) GetUserData(ctx context.Context, instanceID string) (*govultr.UserData, error) {
	panic("implement me")
}

func (f *FakeInstance) GetUpgrades(ctx context.Context, instanceID string) (*govultr.Upgrades, error) {
	panic("implement me")
}

type fakeLB struct {
	client *govultr.Client
}

func (f *fakeLB) Create(ctx context.Context, createReq *govultr.LoadBalancerReq) (*govultr.LoadBalancer, error) {
	panic("implement me")
}

func (f *fakeLB) Get(ctx context.Context, ID string) (*govultr.LoadBalancer, error) {
	return &govultr.LoadBalancer{
		ID:        "6334f227-6d96-4cbd-9bcb-5be0759354fa",
		Region:    "ewr",
		Label:     "lbname",
		Status:    "active",
		IPV4:      "192.168.0.1",
		Instances: []string{"0c51cc3d-529e-4e03-ad86-fd0af47467ba", "ca9a74cb-2d9f-4786-9bb0-094398c593a2"},
	}, nil
}

func (f *fakeLB) Update(ctx context.Context, ID string, updateReq *govultr.LoadBalancerReq) error {
	panic("implement me")
}

func (f *fakeLB) Delete(ctx context.Context, ID string) error {
	panic("implement me")
}

func (f *fakeLB) List(ctx context.Context, options *govultr.ListOptions) ([]govultr.LoadBalancer, *govultr.Meta, error) {
	return []govultr.LoadBalancer{
		{
			ID:     "6334f227-6d96-4cbd-9bcb-5be0759354fa",
			Region: "ewr",
			Label:  "lbname",
			Status: "active",
			IPV4:   "192.168.0.1",
		},
	}, nil, nil
}

func (f *fakeLB) CreateForwardingRule(ctx context.Context, ID string, rule *govultr.ForwardingRule) (*govultr.ForwardingRule, error) {
	panic("implement me")
}

func (f *fakeLB) GetForwardingRule(ctx context.Context, ID string, ruleID string) (*govultr.ForwardingRule, error) {
	panic("implement me")
}

func (f *fakeLB) DeleteForwardingRule(ctx context.Context, ID string, RuleID string) error {
	panic("implement me")
}

func (f *fakeLB) ListForwardingRules(ctx context.Context, ID string, options *govultr.ListOptions) ([]govultr.ForwardingRule, *govultr.Meta, error) {
	return []govultr.ForwardingRule{{
		RuleID:           "1234",
		FrontendProtocol: "tcp",
		FrontendPort:     80,
		BackendProtocol:  "tco",
		BackendPort:      80,
	}}, nil, nil
}
