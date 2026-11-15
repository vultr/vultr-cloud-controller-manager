package vultr

import (
	"context"
	"reflect"
	"testing"

	"github.com/vultr/govultr/v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetHTTPSPorts(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[int32]bool
		wantErr bool
	}{
		{name: "single port", value: "443", want: map[int32]bool{443: true}},
		{name: "comma separated", value: "443,8443", want: map[int32]bool{443: true, 8443: true}},
		{name: "space after comma", value: "443, 8443", want: map[int32]bool{443: true, 8443: true}},
		{name: "not a number", value: "443,https", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &v1.Service{ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{annoVultrLBHTTPSPorts: tt.value},
			}}
			got, err := getHTTPSPorts(svc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getHTTPSPorts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHTTPSPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildFirewallRulesFromAnnotation(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    []govultr.LBFirewallRule
		wantErr string
	}{
		{
			name:  "documented format",
			value: "0.0.0.0/0,80;0.0.0.0/0,90",
			want: []govultr.LBFirewallRule{
				{Source: "0.0.0.0/0", IPType: "v4", Port: 80},
				{Source: "0.0.0.0/0", IPType: "v4", Port: 90},
			},
		},
		{
			name:  "spaces around separators",
			value: "10.0.0.0/8, 80; 2001:db8::/32, 443",
			want: []govultr.LBFirewallRule{
				{Source: "10.0.0.0/8", IPType: "v4", Port: 80},
				{Source: "2001:db8::/32", IPType: "v6", Port: 443},
			},
		},
		{
			name:  "trailing semicolon",
			value: "0.0.0.0/0,80;",
			want: []govultr.LBFirewallRule{
				{Source: "0.0.0.0/0", IPType: "v4", Port: 80},
			},
		},
		{
			name:    "invalid port is reported",
			value:   "0.0.0.0/0,http",
			wantErr: "loadbalancer fw rules : port http is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLoadbalancers(newFakeClient(), "ewr").(*loadbalancers)
			svc := &v1.Service{ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{annoVultrFirewallRules: tt.value},
			}}
			got, err := l.buildFirewallRules(context.Background(), svc)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("buildFirewallRules() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("buildFirewallRules() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildFirewallRules() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
