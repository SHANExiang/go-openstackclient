package entity

import (
	"fmt"
	"go-openstackclient/consts"
	"time"
)

type Rule struct {
	MaxKbps      int    `json:"max_kbps"`
	Direction    string `json:"direction"`
	QosPolicyId  string `json:"qos_policy_id"`
	Type         string `json:"type"`
	Id           string `json:"id"`
	MaxBurstKbps int    `json:"max_burst_kbps"`
}

type Policy struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Rules          []Rule    `json:"rules"`
	Id             string    `json:"id"`
	IsDefault      bool      `json:"is_default"`
	ProjectId      string    `json:"project_id"`
	RevisionNumber int       `json:"revision_number"`
	TenantId       string    `json:"tenant_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Shared         bool      `json:"shared"`
	Tags           []string  `json:"tags"`
}

type QosPolicyMap struct {
	Policy `json:"policy"`
}

type QosPolicies struct {
	Qps                []Policy `json:"policies"`
	Count              int   `json:"count"`
}

type CreateQosPolicyOpts struct {
	// Name is the human-readable name of the QoS policy.
	Name string `json:"name"`

	// TenantID is the id of the Identity project.
	TenantID string `json:"tenant_id,omitempty"`

	// ProjectID is the id of the Identity project.
	ProjectID string `json:"project_id,omitempty"`

	// Shared indicates whether this QoS policy is shared across all projects.
	Shared bool `json:"shared,omitempty"`

	// Description is the human-readable description for the QoS policy.
	Description string `json:"description,omitempty"`

	// IsDefault indicates if this QoS policy is default policy or not.
	IsDefault bool `json:"is_default,omitempty"`
}


func (opts *CreateQosPolicyOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.POLICY)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type BandwidthLimitRule struct {
	Id           string `json:"id"`
	MaxKbps      int    `json:"max_kbps"`
	MaxBurstKbps int    `json:"max_burst_kbps"`
	Direction    string `json:"direction"`
}

type BandwidthLimitRuleMap struct {
    BandwidthLimitRule  	 `json:"bandwidth_limit_rule"`
}

// CreateBandwidthLimitRuleOpts specifies parameters of a new BandwidthLimitRule.
type CreateBandwidthLimitRuleOpts struct {
	// MaxKBps is a maximum kilobits per second. It's a required parameter.
	MaxKBps int `json:"max_kbps"`

	// MaxBurstKBps is a maximum burst size in kilobits.
	MaxBurstKBps int `json:"max_burst_kbps,omitempty"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

func (opts *CreateBandwidthLimitRuleOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.BANDWIDTH_LIMIT_RULE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

// CreateMinimumBandwidthRuleOpts specifies parameters of a new MinimumBandwidthRule.
type CreateMinimumBandwidthRuleOpts struct {
	// MaxKBps is a minimum kilobits per second. It's a required parameter.
	MinKBps int `json:"min_kbps"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// CreateDSCPMarkingRuleOpts specifies parameters of a new DSCPMarkingRule.
type CreateDSCPMarkingRuleOpts struct {
	// DSCPMark contains DSCP mark value.
	DSCPMark int `json:"dscp_mark"`
}

