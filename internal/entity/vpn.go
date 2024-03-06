package entity


type VpnService struct {
	RouterId     string      `json:"router_id"`
	Status       string      `json:"status"`
	Name         string      `json:"name"`
	ExternalV6Ip string      `json:"external_v6_ip"`
	AdminStateUp bool        `json:"admin_state_up"`
	SubnetId     interface{} `json:"subnet_id"`
	ProjectId    string      `json:"project_id"`
	TenantId     string      `json:"tenant_id"`
	ExternalV4Ip string      `json:"external_v4_ip"`
	Id           string      `json:"id"`
	Description  string      `json:"description"`
	FlavorId     interface{} `json:"flavor_id"`
}

type VpnServiceMap struct {
	 VpnService `json:"vpnservice"`
}

type VpnServices struct {
	Vss              []VpnService `json:"vpnservices"`
	Count            int          `json:"count"`
}

type EndpointGroup struct {
	Description string   `json:"description"`
	ProjectId   string   `json:"project_id"`
	TenantId    string   `json:"tenant_id"`
	Endpoints   []string `json:"endpoints"`
	Type        string   `json:"type"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
}

type EndpointGroupMap struct {
	EndpointGroup `json:"endpoint_group"`
}

type EndpointGroups struct {
	Egs              []EndpointGroup `json:"endpoint_groups"`
	Count              int           `json:"count"`
}

type Ikepolicy struct {
	Name                  string `json:"name"`
	ProjectId             string `json:"project_id"`
	TenantId              string `json:"tenant_id"`
	AuthAlgorithm         string `json:"auth_algorithm"`
	EncryptionAlgorithm   string `json:"encryption_algorithm"`
	Pfs                   string `json:"pfs"`
	Phase1NegotiationMode string `json:"phase1_negotiation_mode"`
	Lifetime              struct {
		Units string `json:"units"`
		Value int    `json:"value"`
	} `json:"lifetime"`
	IkeVersion  string `json:"ike_version"`
	Id          string `json:"id"`
	Description string `json:"description"`
}

type IkePolicyMap struct {
	Ikepolicy `json:"ikepolicy"`
}

type IkePolicies struct {
	Ips               []Ikepolicy `json:"ikepolicies"`
	Count              int        `json:"count"`
}

type Ipsecpolicy struct {
	Name                string `json:"name"`
	TransformProtocol   string `json:"transform_protocol"`
	AuthAlgorithm       string `json:"auth_algorithm"`
	EncapsulationMode   string `json:"encapsulation_mode"`
	EncryptionAlgorithm string `json:"encryption_algorithm"`
	Pfs                 string `json:"pfs"`
	ProjectId           string `json:"project_id"`
	TenantId            string `json:"tenant_id"`
	Lifetime            struct {
		Units string `json:"units"`
		Value int    `json:"value"`
	} `json:"lifetime"`
	Id          string `json:"id"`
	Description string `json:"description"`
}

type IpsecPolicyMap struct {
	Ipsecpolicy `json:"ipsecpolicy"`
}

type IpsecPolicies struct {
	Ips             []Ipsecpolicy `json:"ipsecpolicies"`
	Count              int        `json:"count"`
}

type IpsecSiteConnection struct {
	Status        string        `json:"status"`
	Psk           string        `json:"psk"`
	Initiator     string        `json:"initiator"`
	Name          string        `json:"name"`
	AdminStateUp  bool          `json:"admin_state_up"`
	ProjectId     string        `json:"project_id"`
	TenantId      string        `json:"tenant_id"`
	AuthMode      string        `json:"auth_mode"`
	PeerCidrs     []interface{} `json:"peer_cidrs"`
	Mtu           int           `json:"mtu"`
	PeerEpGroupId string        `json:"peer_ep_group_id"`
	IkepolicyId   string        `json:"ikepolicy_id"`
	VpnserviceId  string        `json:"vpnservice_id"`
	Dpd           struct {
		Action   string `json:"action"`
		Interval int    `json:"interval"`
		Timeout  int    `json:"timeout"`
	} `json:"dpd"`
	RouteMode      string `json:"route_mode"`
	IpsecpolicyId  string `json:"ipsecpolicy_id"`
	LocalEpGroupId string `json:"local_ep_group_id"`
	PeerAddress    string `json:"peer_address"`
	PeerId         string `json:"peer_id"`
	Id             string `json:"id"`
	Description    string `json:"description"`
}

type IpsecConnectionMap struct {
	IpsecSiteConnection `json:"ipsec_site_connection"`
}

type IpsecConnections struct {
	ICs          []IpsecSiteConnection `json:"ipsec_site_connections"`
	Count              int             `json:"count"`
}
