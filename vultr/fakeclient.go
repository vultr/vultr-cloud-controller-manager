package vultr

import (
	"context"

	"github.com/vultr/govultr"
)

func newFakeClient() *govultr.Client {
	fakeInstance := FakeInstance{client: nil}
	fakeLoadBalancer := fakeLB{nil}
	return &govultr.Client{
		Server:       &fakeInstance,
		LoadBalancer: &fakeLoadBalancer,
	}
}

func newFakeInstance() *govultr.Server {
	return &govultr.Server{
		InstanceID:  "576965",
		MainIP:      "149.28.225.110",
		VPSCpus:     "4",
		Location:    "New Jersey",
		RegionID:    "1",
		Status:      "running",
		NetmaskV4:   "255.255.254.0",
		GatewayV4:   "149.28.224.1",
		PowerStatus: "",
		ServerState: "",
		PlanID:      "204",
		Label:       "cluster-name",
		InternalIP:  "10.1.95.4",
	}
}

type FakeInstance struct {
	client *govultr.Client
}

func (f *FakeInstance) ChangeApp(ctx context.Context, instanceID, appID string) error {
	panic("implement me")
}

func (f *FakeInstance) ListApps(ctx context.Context, instanceID string) ([]govultr.Application, error) {
	panic("implement me")
}

func (f *FakeInstance) AppInfo(ctx context.Context, instanceID string) (*govultr.AppInfo, error) {
	panic("implement me")
}

func (f *FakeInstance) EnableBackup(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) DisableBackup(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) GetBackupSchedule(ctx context.Context, instanceID string) (*govultr.BackupSchedule, error) {
	panic("implement me")
}

func (f *FakeInstance) SetBackupSchedule(ctx context.Context, instanceID string, backup *govultr.BackupSchedule) error {
	panic("implement me")
}

func (f *FakeInstance) RestoreBackup(ctx context.Context, instanceID, backupID string) error {
	panic("implement me")
}

func (f *FakeInstance) RestoreSnapshot(ctx context.Context, instanceID, snapshotID string) error {
	panic("implement me")
}

func (f *FakeInstance) SetLabel(ctx context.Context, instanceID, label string) error {
	panic("implement me")
}

func (f *FakeInstance) SetTag(ctx context.Context, instanceID, tag string) error {
	panic("implement me")
}

func (f *FakeInstance) Neighbors(ctx context.Context, instanceID string) ([]int, error) {
	panic("implement me")
}

func (f *FakeInstance) EnablePrivateNetwork(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

func (f *FakeInstance) DisablePrivateNetwork(ctx context.Context, instanceID, networkID string) error {
	panic("implement me")
}

func (f *FakeInstance) ListPrivateNetworks(ctx context.Context, instanceID string) ([]govultr.PrivateNetwork, error) {
	panic("implement me")
}

func (f *FakeInstance) ListUpgradePlan(ctx context.Context, instanceID string) ([]int, error) {
	panic("implement me")
}

func (f *FakeInstance) UpgradePlan(ctx context.Context, instanceID, vpsPlanID string) error {
	panic("implement me")
}

func (f *FakeInstance) ListOS(ctx context.Context, instanceID string) ([]govultr.OS, error) {
	panic("implement me")
}

func (f *FakeInstance) ChangeOS(ctx context.Context, instanceID, osID string) error {
	panic("implement me")
}

func (f *FakeInstance) IsoAttach(ctx context.Context, instanceID, isoID string) error {
	panic("implement me")
}

func (f *FakeInstance) IsoDetach(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) IsoStatus(ctx context.Context, instanceID string) (*govultr.ServerIso, error) {
	panic("implement me")
}

func (f *FakeInstance) SetFirewallGroup(ctx context.Context, instanceID, firewallGroupID string) error {
	panic("implement me")
}

func (f *FakeInstance) GetUserData(ctx context.Context, instanceID string) (*govultr.UserData, error) {
	panic("implement me")
}

func (f *FakeInstance) SetUserData(ctx context.Context, instanceID, userData string) error {
	panic("implement me")
}

func (f *FakeInstance) IPV4Info(ctx context.Context, instanceID string, public bool) ([]govultr.IPV4, error) {
	panic("implement me")
}

func (f *FakeInstance) IPV6Info(ctx context.Context, instanceID string) ([]govultr.IPV6, error) {
	panic("implement me")
}

func (f *FakeInstance) AddIPV4(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) DestroyIPV4(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

func (f *FakeInstance) EnableIPV6(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) Bandwidth(ctx context.Context, instanceID string) ([]map[string]string, error) {
	panic("implement me")
}

func (f *FakeInstance) ListReverseIPV6(ctx context.Context, instanceID string) ([]govultr.ReverseIPV6, error) {
	panic("implement me")
}

func (f *FakeInstance) SetDefaultReverseIPV4(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

func (f *FakeInstance) DeleteReverseIPV6(ctx context.Context, instanceID, ip string) error {
	panic("implement me")
}

func (f *FakeInstance) SetReverseIPV4(ctx context.Context, instanceID, ipv4, entry string) error {
	panic("implement me")
}

func (f *FakeInstance) SetReverseIPV6(ctx context.Context, instanceID, ipv6, entry string) error {
	panic("implement me")
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

func (f *FakeInstance) Delete(ctx context.Context, instanceID string) error {
	panic("implement me")
}

func (f *FakeInstance) Create(ctx context.Context, regionID, vpsPlanID, osID int, options *govultr.ServerOptions) (*govultr.Server, error) {
	panic("implement me")
}

func (f *FakeInstance) List(ctx context.Context) ([]govultr.Server, error) {
	panic("implement me")
}

func (f *FakeInstance) ListByLabel(ctx context.Context, label string) ([]govultr.Server, error) {

	return []govultr.Server{
		{
			InstanceID:  "576965",
			MainIP:      "149.28.225.110",
			VPSCpus:     "4",
			Location:    "New Jersey",
			RegionID:    "1",
			Status:      "running",
			NetmaskV4:   "255.255.254.0",
			GatewayV4:   "149.28.224.1",
			PowerStatus: "",
			ServerState: "",
			PlanID:      "204",
			Label:       "ccm-test",
			InternalIP:  "10.1.95.4",
		},
	}, nil
}

func (f *FakeInstance) ListByMainIP(ctx context.Context, mainIP string) ([]govultr.Server, error) {
	panic("implement me")
}

func (f *FakeInstance) ListByTag(ctx context.Context, tag string) ([]govultr.Server, error) {
	panic("implement me")
}

func (f *FakeInstance) GetServer(ctx context.Context, instanceID string) (*govultr.Server, error) {
	return &govultr.Server{
		InstanceID:  "576965",
		MainIP:      "149.28.225.110",
		VPSCpus:     "4",
		Location:    "New Jersey",
		RegionID:    "1",
		Status:      "running",
		NetmaskV4:   "255.255.254.0",
		GatewayV4:   "149.28.224.1",
		PowerStatus: "",
		ServerState: "",
		PlanID:      "204",
		Label:       "ccm-test",
		InternalIP:  "10.1.95.4",
	}, nil
}

type fakeLB struct {
	client *govultr.Client
}

func (l *fakeLB) List(ctx context.Context) ([]govultr.LoadBalancers, error) {
	return []govultr.LoadBalancers{
		{
			ID:       12345,
			RegionID: 1,
			Label:    "albname",
			Status:   "active",
			IPV4:     "192.168.0.1",
		},
	}, nil
}

func (l *fakeLB) Delete(ctx context.Context, ID int) error {
	panic("implement me")
}

func (l *fakeLB) SetLabel(ctx context.Context, ID int, label string) error {
	panic("implement me")
}

func (l *fakeLB) AttachedInstances(ctx context.Context, ID int) (*govultr.InstanceList, error) {
	return &govultr.InstanceList{
		InstanceList: []int{
			123,
			124,
		},
	}, nil
}

func (l *fakeLB) AttachInstance(ctx context.Context, ID, backendNode int) error {
	panic("implement me")
}

func (l *fakeLB) DetachInstance(ctx context.Context, ID, backendNode int) error {
	panic("implement me")
}

func (l *fakeLB) GetHealthCheck(ctx context.Context, ID int) (*govultr.HealthCheck, error) {
	panic("implement me")
}

func (l *fakeLB) SetHealthCheck(ctx context.Context, ID int, healthConfig *govultr.HealthCheck) error {
	return nil
}

func (l *fakeLB) GetGenericInfo(ctx context.Context, ID int) (*govultr.GenericInfo, error) {
	panic("implement me")
}

func (l *fakeLB) ListForwardingRules(ctx context.Context, ID int) (*govultr.ForwardingRules, error) {
	return &govultr.ForwardingRules{
		ForwardRuleList: []govultr.ForwardingRule{
			{
				RuleID:           "1234",
				FrontendProtocol: "tcp",
				FrontendPort:     80,
				BackendProtocol:  "tcp",
				BackendPort:      8080,
			},
		},
	}, nil
}

func (l *fakeLB) DeleteForwardingRule(ctx context.Context, ID int, RuleID string) error {
	panic("implement me")
}

func (l *fakeLB) CreateForwardingRule(ctx context.Context, ID int, rule *govultr.ForwardingRule) (*govultr.ForwardingRule, error) {
	panic("implement me")
}

func (l *fakeLB) GetFullConfig(ctx context.Context, ID int) (*govultr.LBConfig, error) {
	panic("implement me")
}

func (l *fakeLB) HasSSL(ctx context.Context, ID int) (*struct {
	SSLInfo bool `json:"has_ssl"`
}, error) {
	panic("implement me")
}

func (l *fakeLB) Create(ctx context.Context, region int, label string, genericInfo *govultr.GenericInfo, healthCheck *govultr.HealthCheck, rules []govultr.ForwardingRule, ssl *govultr.SSL) (*govultr.LoadBalancers, error) {
	panic("implement me")
}

func (l *fakeLB) UpdateGenericInfo(ctx context.Context, ID int, label string, genericInfo *govultr.GenericInfo) error {
	return nil
}

func (l *fakeLB) AddSSL(ctx context.Context, ID int, ssl *govultr.SSL) error {
	panic("implement me")
}

func (l *fakeLB) RemoveSSL(ctx context.Context, ID int) error {
	panic("implement me")
}
