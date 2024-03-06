package entity

import (
	"fmt"
)

type NetworkQuota struct {
	Subnet              int `json:"subnet,omitempty"`
	Ikepolicy           int `json:"ikepolicy,omitempty"`
	Subnetpool          int `json:"subnetpool,omitempty"`
	FirewallRule        int `json:"firewall_rule,omitempty"`
	Network             int `json:"network,omitempty"`
	IpsecSiteConnection int `json:"ipsec_site_connection,omitempty"`
	EndpointGroup       int `json:"endpoint_group,omitempty"`
	Firewall            int `json:"firewall,omitempty"`
	Ipsecpolicy         int `json:"ipsecpolicy,omitempty"`
	FirewallPolicy      int `json:"firewall_policy,omitempty"`
	SecurityGroupRule   int `json:"security_group_rule,omitempty"`
	Vpnservice          int `json:"vpnservice,omitempty"`
	Floatingip          int `json:"floatingip,omitempty"`
	SecurityGroup       int `json:"security_group,omitempty"`
	Router              int `json:"router,omitempty"`
	RbacPolicy          int `json:"rbac_policy,omitempty"`
	Port                int `json:"port,omitempty"`
}

type NetworkQuotaMap struct {
	NetworkQuota `json:"quota"`
}

func (opts *NetworkQuotaMap) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}