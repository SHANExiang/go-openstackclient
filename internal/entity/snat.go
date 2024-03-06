package entity

import (
    "fmt"
    "go-openstackclient/consts"
)

type Snat struct {
    Id            string `json:"id,omitempty"`
    Name          string `json:"name,omitempty"`
    TenantId      string `json:"tenant_id" required:"true"`
    RouterId      string `json:"router_id" required:"true"`
    SnatNetworkId string `json:"snat_network_id" required:"true"`
    SnatIpAddress string `json:"snat_ip_address"`
    SnatIpPool    []struct {
        BeginIp string `json:"begin_ip"`
        EndIp   string `json:"end_ip"`
    } `json:"snat_ip_pool,omitempty"`
    OriginalCidrs []string `json:"original_cidrs" required:"true"`
}

type SnatMap struct {
     Snat        `json:"snat"`
}

func (opts *Snat) ToRequestBody() string {
    reqBody, err := BuildRequestBody(opts, consts.Snat)
    if err != nil {
        panic(fmt.Sprintf("Failed to build request body %s", err))
    }
    return reqBody
}

type Snats struct {
    Ss           []Snat         `json:"snats"`
}

type Dnat struct {
    Id                string `json:"id,omitempty"`
    Name              string `json:"name,omitempty"`
    TenantId          string `json:"tenant_id,omitempty"`
    FloatingNetworkId string `json:"floating_network_id,omitempty"`
    PortId            string `json:"port_id,omitempty"`
    FixedIpAddress    string `json:"fixed_ip_address,omitempty"`
    FloatingIpAddress string `json:"floating_ip_address,omitempty"`
    RouterId          string `json:"router_id,omitempty"`
    FloatingIpPort    int    `json:"floating_ip_port,omitempty"`
    Protocol          string `json:"protocol,omitempty"`
    FloatingipId      string `json:"floating_ip_id,omitempty"`
    FixedIpPort       int    `json:"fixed_ip_port,omitempty"`
}

type DnatMap struct {
     Dnat        `json:"dnat"`
}

func (opts *Dnat) ToRequestBody() string {
    reqBody, err := BuildRequestBody(opts, consts.Dnat)
    if err != nil {
        panic(fmt.Sprintf("Failed to build request body %s", err))
    }
    return reqBody
}

type Dnats struct {
    Ds           []Dnat         `json:"dnats"`
}
