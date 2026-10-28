package vultr

import (
	"reflect"
	"testing"

	"github.com/vultr/govultr/v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoadBalancerCNAMEIngress(t *testing.T) {
	for _, tt := range []struct {
		name        string
		annotations map[string]string
		cname       string
		dualStack   bool
		want        []v1.LoadBalancerIngress
		wantError   bool
	}{
		{
			name: "managed CNAME", annotations: map[string]string{annoVultrUseCNAME: "true"},
			cname: "vlb-test.vultrlb.com", want: []v1.LoadBalancerIngress{{Hostname: "vlb-test.vultrlb.com"}},
		},
		{
			name: "dual stack CNAME", annotations: map[string]string{annoVultrUseCNAME: "true"}, dualStack: true,
			cname: "vlb-test.vultrlb.com", want: []v1.LoadBalancerIngress{{Hostname: "vlb-test.vultrlb.com"}},
		},
		{
			name:        "explicit hostname takes precedence",
			annotations: map[string]string{annoVultrUseCNAME: "true", annoVultrHostname: "app.example.com"},
			want:        []v1.LoadBalancerIngress{{Hostname: "app.example.com"}},
		},
		{
			name:        "empty hostname allows CNAME",
			annotations: map[string]string{annoVultrUseCNAME: "true", annoVultrHostname: ""},
			cname:       "vlb-test.vultrlb.com", want: []v1.LoadBalancerIngress{{Hostname: "vlb-test.vultrlb.com"}},
		},
		{
			name: "disabled", annotations: map[string]string{annoVultrUseCNAME: "false"}, cname: "vlb-test.vultrlb.com",
			want: []v1.LoadBalancerIngress{{Hostname: "test", IP: "192.0.2.1"}},
		},
		{
			name: "default", cname: "vlb-test.vultrlb.com",
			want: []v1.LoadBalancerIngress{{Hostname: "test", IP: "192.0.2.1"}},
		},
		{
			name: "dual stack IP default", dualStack: true,
			want: []v1.LoadBalancerIngress{{Hostname: "test", IP: "192.0.2.1"}, {Hostname: "test", IP: "2001:db8::1"}},
		},
		{
			name: "invalid boolean", annotations: map[string]string{annoVultrUseCNAME: "invalid"}, wantError: true,
		},
		{
			name:        "invalid explicit hostname",
			annotations: map[string]string{annoVultrUseCNAME: "true", annoVultrHostname: "not a hostname"}, wantError: true,
		},
		{
			name: "not ready", annotations: map[string]string{annoVultrUseCNAME: "true"}, wantError: true,
		},
		{
			name: "invalid CNAME", annotations: map[string]string{annoVultrUseCNAME: "true"},
			cname: "not a hostname", wantError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := &loadbalancers{}
			svc := &v1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: tt.annotations}}
			if tt.dualStack {
				svc.Spec.IPFamilies = []v1.IPFamily{v1.IPv4Protocol, v1.IPv6Protocol}
			}
			lb := &govultr.LoadBalancer{ID: "test-id", Label: "test", IPV4: "192.0.2.1", IPV6: "2001:db8::1", CNAME: tt.cname}
			got, err := l.buildLoadBalancerIngress(svc, lb)
			if (err != nil) != tt.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ingress = %+v, want %+v", got, tt.want)
			}
		})
	}
}
