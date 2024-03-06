package entity

import (
	"fmt"
	"go-openstackclient/consts"
)

type VpcConnection struct {
	Status              string   `json:"status"`
	Id                  string   `json:"id"`
	TenantId            string   `json:"tenant_id"`
	LocalRouter         string   `json:"local_router"`
	LocalSubnets        []string `json:"local_subnets"`
	LocalCidrs          []string `json:"local_cidrs"`
	FwEnabled           bool     `json:"fw_enabled"`
	LocalFirewallEnable bool     `json:"local_firewall_enable"`
	PeerRouter          string   `json:"peer_router"`
	PeerSubnets         []string `json:"peer_subnets"`
	PeerCidrs           []string `json:"peer_cidrs"`
	PeerFirewallEnable  bool     `json:"peer_firewall_enable"`
	Mode                int      `json:"mode"`
	Priority            int      `json:"priority"`
	Name                string   `json:"name"`
}

type VpcConnectionMap struct {
	VpcConnection `json:"vpc_connection"`
}

type VpcConnections struct {
	Vcs                    []VpcConnection `json:"vpc_connections"`
}

type CreateVpcConnectionOpts struct {
	Name                string   `json:"name,omitempty"`
	TenantId            string   `json:"tenant_id,omitempty"`
	LocalRouter         string   `json:"local_router"`
	LocalSubnets        []string `json:"local_subnets,omitempty"`
	LocalCidrs          []string `json:"local_cidrs,omitempty"`
	FwEnabled           bool     `json:"fw_enabled,omitempty"`
	LocalFirewallEnable bool     `json:"local_firewall_enable,omitempty"`
	PeerRouter          string   `json:"peer_router"`
	PeerSubnets         []string `json:"peer_subnets,omitempty"`
	PeerCidrs           []string `json:"peer_cidrs,omitempty"`
	PeerFirewallEnable  bool     `json:"peer_firewall_enable,omitempty"`
	Mode                int      `json:"mode,omitempty"`
	Priority            int      `json:"priority,omitempty"`
}

func (opts *CreateVpcConnectionOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.VpcConnection)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}


type UpdateVpcConnectionOpts struct {
	Name                string   `json:"name,omitempty"`
	TenantId            string   `json:"tenant_id,omitempty"`
	LocalRouter         string   `json:"local_router,omitempty"`
	LocalSubnets        []string `json:"local_subnets,omitempty"`
	LocalCidrs          []string `json:"local_cidrs,omitempty"`
	FwEnabled           bool     `json:"fw_enabled,omitempty"`
	LocalFirewallEnable bool     `json:"local_firewall_enable,omitempty"`
	PeerRouter          string   `json:"peer_router,omitempty"`
	PeerSubnets         []string `json:"peer_subnets,omitempty"`
	PeerCidrs           []string `json:"peer_cidrs,omitempty"`
	PeerFirewallEnable  bool     `json:"peer_firewall_enable,omitempty"`
	Mode                int      `json:"mode,omitempty"`
	Priority            int      `json:"priority,omitempty"`
}

func (opts *UpdateVpcConnectionOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.VpcConnection)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}