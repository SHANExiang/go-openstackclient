package entity



type FirewallGroup struct {
	Ports               []string         `json:"ports"`
}

type FirewallPolicy struct {
	Name                  string        `json:"name"`
	FirewallRules         []string      `json:"firewall_rules"`
}

