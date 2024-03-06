package entity

import (
	"fmt"
	"go-openstackclient/consts"
)


type FirewallRuleV1 struct {
	Protocol             interface{} `json:"protocol"`
	Description          string      `json:"description"`
	SourcePort           interface{} `json:"source_port"`
	SourceIpAddress      interface{} `json:"source_ip_address"`
	DestinationIpAddress interface{} `json:"destination_ip_address"`
	FirewallPolicyId     interface{} `json:"firewall_policy_id"`
	Position             interface{} `json:"position"`
	DestinationPort      interface{} `json:"destination_port"`
	Id                   string      `json:"id"`
	Name                 string      `json:"name"`
	TenantId             string      `json:"tenant_id"`
	Enabled              bool        `json:"enabled"`
	Action               string      `json:"action"`
	IpVersion            int         `json:"ip_version"`
	Shared               bool        `json:"shared"`
	ProjectId            string      `json:"project_id"`
}

type FirewallRuleV1Map struct {
	FirewallRuleV1 `json:"firewall_rule"`
}

type FirewallRules struct {
	Frs               []FirewallRuleV1 `json:"firewall_rules"`
}

type FirewallPolicyV1 struct {
	Name          string        `json:"name"`
	FirewallRules []interface{} `json:"firewall_rules"`
	TenantId      string        `json:"tenant_id"`
	Id            string        `json:"id"`
	Shared        bool          `json:"shared"`
	ProjectId     string        `json:"project_id"`
	Audited       bool          `json:"audited"`
	Description   string        `json:"description"`
}

type FirewallPolicyV1Map struct {
	FirewallPolicyV1 `json:"firewall_policy"`
}

type FirewallPolicies struct {
	Fps             []FirewallPolicyV1 `json:"firewall_policies"`
}

type Firewall struct {
	Status           string        `json:"status"`
	RouterIds        []interface{} `json:"router_ids"`
	Name             string        `json:"name"`
	AdminStateUp     bool          `json:"admin_state_up"`
	TenantId         string        `json:"tenant_id"`
	FirewallPolicyId string        `json:"firewall_policy_id"`
	ProjectId        string        `json:"project_id"`
	Id               string        `json:"id"`
	Description      string        `json:"description"`
}

type FirewallV1Map struct {
	Firewall `json:"firewall"`
}

type FirewallV1s struct {
	Fs             []Firewall `json:"firewalls"`
}


type CreateFirewallOpts struct {
	PolicyID string `json:"firewall_policy_id" required:"true"`
	// TenantID specifies a tenant to own the firewall. The caller must have
	// an admin role in order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID     string `json:"tenant_id,omitempty"`
	ProjectID    string `json:"project_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	AdminStateUp *bool  `json:"admin_state_up,omitempty"`
	Shared       *bool  `json:"shared,omitempty"`
	RouterIDs    []string `json:"router_ids"`
}

func (opts *CreateFirewallOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.FIREWALL)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type UpdateFirewallOpts struct {
	PolicyID     string  `json:"firewall_policy_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	AdminStateUp *bool   `json:"admin_state_up,omitempty"`
	Shared       *bool   `json:"shared,omitempty"`
	RouterIDs    []string `json:"router_ids,omitempty"`
}

func (opts *UpdateFirewallOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.FIREWALL)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type CreateFirewallPolicyOpts struct {
	// TenantID specifies a tenant to own the firewall. The caller must have
	// an admin role in order to set this. Otherwise, this field is left unset
	// and the caller will be the owner.
	TenantID    string   `json:"tenant_id,omitempty"`
	ProjectID   string   `json:"project_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Shared      *bool    `json:"shared,omitempty"`
	Audited     *bool    `json:"audited,omitempty"`
	Rules       []string `json:"firewall_rules,omitempty"`
}

func (opts *CreateFirewallPolicyOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.FIREWALLPOLICY)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type UpdateFirewallPolicyOpts struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Shared      *bool    `json:"shared,omitempty"`
	Audited     *bool    `json:"audited,omitempty"`
	Rules       []string `json:"firewall_rules,omitempty"`
}

func (opts *UpdateFirewallPolicyOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.FIREWALLPOLICY)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type CreateFirewallRuleOpts struct {
	Protocol             string                `json:"protocol" required:"true"`
	Action               string                `json:"action" required:"true"`
	TenantID             string                `json:"tenant_id,omitempty"`
	ProjectID            string                `json:"project_id,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Description          string                `json:"description,omitempty"`
	IPVersion            int                   `json:"ip_version,omitempty"`
	SourceIPAddress      string                `json:"source_ip_address,omitempty"`
	DestinationIPAddress string                `json:"destination_ip_address,omitempty"`
	SourcePort           string                `json:"source_port,omitempty"`
	DestinationPort      string                `json:"destination_port,omitempty"`
	Shared               *bool                 `json:"shared,omitempty"`
	Enabled              *bool                 `json:"enabled,omitempty"`
}

func (opts *CreateFirewallRuleOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.FIREWALLRULE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type UpdateFirewallRuleOpts struct {
	Protocol             *string                `json:"protocol,omitempty"`
	Action               *string                `json:"action,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	IPVersion            int                    `json:"ip_version,omitempty"`
	SourceIPAddress      *string                `json:"source_ip_address,omitempty"`
	DestinationIPAddress *string                `json:"destination_ip_address,omitempty"`
	SourcePort           *string                `json:"source_port,omitempty"`
	DestinationPort      *string                `json:"destination_port,omitempty"`
	Shared               *bool                  `json:"shared,omitempty"`
	Enabled              *bool                  `json:"enabled,omitempty"`
}

func (opts *UpdateFirewallRuleOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.FIREWALLRULE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}
