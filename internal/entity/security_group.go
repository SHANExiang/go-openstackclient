package entity

import (
	"fmt"
	"go-openstackclient/consts"
	"time"
)

type SecurityGroup struct {
	Description        string `json:"description"`
	Id                 string `json:"id"`
	Name               string `json:"name"`
	SecurityGroupRules []struct {
		Direction       string      `json:"direction"`
		Ethertype       string      `json:"ethertype"`
		Id              string      `json:"id"`
		PortRangeMax    interface{} `json:"port_range_max"`
		PortRangeMin    interface{} `json:"port_range_min"`
		Protocol        interface{} `json:"protocol"`
		RemoteGroupId   interface{} `json:"remote_group_id"`
		RemoteIpPrefix  interface{} `json:"remote_ip_prefix"`
		SecurityGroupId string      `json:"security_group_id"`
		ProjectId       string      `json:"project_id"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		RevisionNumber  int         `json:"revision_number"`
		RevisioNNumber  int         `json:"revisio[n_number,omitempty"`
		Tags            []string    `json:"tags"`
		TenantId        string      `json:"tenant_id"`
		Description     string      `json:"description"`
	} `json:"security_group_rules"`
	ProjectId      string    `json:"project_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	RevisionNumber int       `json:"revision_number"`
	Tags           []string  `json:"tags"`
	TenantId       string    `json:"tenant_id"`
	Stateful       bool      `json:"stateful"`
	Shared         bool      `json:"shared"`
}

type SecurityGroupRule struct {
	Direction       string      `json:"direction"`
	Ethertype       string      `json:"ethertype"`
	Id              string      `json:"id"`
	PortRangeMax    int         `json:"port_range_max"`
	PortRangeMin    int         `json:"port_range_min"`
	Protocol        string      `json:"protocol"`
	RemoteGroupId   string      `json:"remote_group_id"`
	RemoteIpPrefix  interface{} `json:"remote_ip_prefix"`
	SecurityGroupId string      `json:"security_group_id"`
	ProjectId       string      `json:"project_id"`
	RevisionNumber  int         `json:"revision_number"`
	TenantId        string      `json:"tenant_id"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Description     string      `json:"description"`
}

type Sg struct {
	SecurityGroup `json:"security_group"`
}

type SgRule struct {
	SecurityGroupRule `json:"security_group_rule"`
}

type SgRules struct {
	Srs         []SecurityGroupRule `json:"security_group_rules"`
}

type Sgs struct {
	Sgs            []SecurityGroup `json:"security_groups"`
	Count          int             `json:"count"`
}

type CreateSecurityGroupOpts struct {
	// Human-readable name for the Security Group. Does not have to be unique.
	Name string `json:"name" required:"true"`

	// TenantID is the UUID of the project who owns the Group.
	// Only administrative users can specify a tenant UUID other than their own.
	TenantID string `json:"tenant_id,omitempty"`

	// ProjectID is the UUID of the project who owns the Group.
	// Only administrative users can specify a tenant UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// Describes the security group.
	Description string `json:"description,omitempty"`
}

// CreateSecurityRuleOpts contains all the values needed to create a new security group
// rule.
type CreateSecurityRuleOpts struct {
	// Must be either "ingress" or "egress": the direction in which the security
	// group rule is applied.
	Direction string `json:"direction" required:"true"`

	// String description of each rule, optional
	Description string `json:"description,omitempty"`

	// Must be "IPv4" or "IPv6", and addresses represented in CIDR must match the
	// ingress or egress rules.
	EtherType string `json:"ethertype" required:"true"`

	// The security group ID to associate with this security group rule.
	SecGroupID string `json:"security_group_id" required:"true"`

	// The maximum port number in the range that is matched by the security group
	// rule. The PortRangeMin attribute constrains the PortRangeMax attribute. If
	// the protocol is ICMP, this value must be an ICMP type.
	PortRangeMax int `json:"port_range_max,omitempty"`

	// The minimum port number in the range that is matched by the security group
	// rule. If the protocol is TCP or UDP, this value must be less than or equal
	// to the value of the PortRangeMax attribute. If the protocol is ICMP, this
	// value must be an ICMP type.
	PortRangeMin int `json:"port_range_min,omitempty"`

	// The protocol that is matched by the security group rule. Valid values are
	// "tcp", "udp", "icmp" or an empty string.
	Protocol string `json:"protocol,omitempty"`

	// The remote group ID to be associated with this security group rule. You can
	// specify either RemoteGroupID or RemoteIPPrefix.
	RemoteGroupID string `json:"remote_group_id,omitempty"`

	// The remote IP prefix to be associated with this security group rule. You can
	// specify either RemoteGroupID or RemoteIPPrefix. This attribute matches the
	// specified IP prefix as the source IP address of the IP packet.
	RemoteIPPrefix string `json:"remote_ip_prefix,omitempty"`

	// TenantID is the UUID of the project who owns the Rule.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`
}


func (opts *CreateSecurityGroupOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.SECURITYGROUP)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}


func (opts *CreateSecurityRuleOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.SECURITYGROUPRULE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}
