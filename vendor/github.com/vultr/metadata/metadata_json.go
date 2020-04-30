package metadata

type MetaData struct {
	Hostname   string `json:"hostname,omitempty"`
	InstanceID string `json:"instanceid,omitempty"`
	PublicKeys string `json:"public-keys,omitempty"`

	Region struct {
		RegionCode string `json:"regioncode,omitempty"`
	} `json:"region,omitempty"`

	BGP struct {
		IPv4 struct {
			MyAddress   string `json:"my-address,omitempty"`
			MyASN       string `json:"my-asn,omitempty"`
			PeerAddress string `json:"peer-address,omitempty"`
			PeerASN     string `json:"peer-asn,omitempty"`
		} `json:"ipv4,omitempty"`
		IPv6 struct {
			MyAddress   string `json:"my-address,omitempty"`
			MyASN       string `json:"my-asn,omitempty"`
			PeerAddress string `json:"peer-address,omitempty"`
			PeerASN     string `json:"peer-asn,omitempty"`
		} `json:"ipv6,omitempty"`
	} `json:"bgp,omitempty"`

	Interfaces []struct {
		IPv4 struct {
			Additional []struct {
				Address string `json:"address,omitempty"`
				Netmask string `json:"netmask,omitempty"`
			}
			Address string `json:"address,omitempty"`
			Gateway string `json:"gateway,omitempty"`
			Netmask string `json:"netmask,omitempty"`
		} `json:"ipv4,omitempty"`
		IPv6 struct {
			Additional []struct {
				Address string `json:"address,omitempty"`
				Prefix  string `json:"prefix,omitempty"`
			}
			Address string `json:"address,omitempty"`
			Gateway string `json:"gateway,omitempty"`
			Prefix  string `json:"prefix,omitempty"`
		} `json:"ipv6,omitempty"`
		Mac         string `json:"mac,omitempty"`
		NetworkType string `json:"network-type,omitempty"`
		NetworkID   string `json:"networkid,omitempty"`
	} `json:"interfaces,omitempty"`
}
