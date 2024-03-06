package entity

import "time"

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
