package service

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
)

var supportedNeutronResourceTypes = map[string]struct{}{
	consts.NETWORK: struct{}{}, consts.SUBNET: struct {}{},
	consts.PORT: struct{}{}, consts.SECURITYGROUPRULE: struct{}{},
	consts.SECURITYGROUP: struct{}{}, consts.BANDWIDTH_LIMIT_RULE: struct{}{},
	consts.DSCP_MARKING_RULE: struct{}{}, consts.MINIMUM_BANDWIDTH_RULE: struct{}{},
	consts.QOS_POLICY: struct{}{}, consts.ROUTER: struct{}{},
	consts.ROUTERINTERFACE: struct{}{}, consts.ROUTERGATEWAY: struct{}{},
	consts.ROUTERROUTE: struct{}{}, consts.FLOATINGIP: struct{}{},
	consts.PORTFORWARDING: struct{}{}, consts.FIREWALLRULE: struct{}{},
	consts.FIREWALLPOLICY: struct{}{}, consts.FIREWALL: struct{}{},
	consts.VpcConnection: struct{}{},
}

type Neutron struct {
	client          client.Client
	httpPrefix      string
}

func newNeutron() *Neutron {
	return &Neutron{
		client: client.NewClient(),
		httpPrefix: fmt.Sprintf("http://%s:%d/v2.0/", configs.CONF.Host, consts.NeutronPort),
	}
}

func (n *Neutron) SupportedResources() map[string]struct{} {
	return supportedNeutronResourceTypes
}


func (n *Neutron) HttpPrefix() string {
	return fmt.Sprintf("http://%s:%d/v2.0", configs.CONF.Host, consts.NeutronPort)
}

func (n *Neutron) Client() client.Client {
	if n.client == nil {
		n.client = client.NewClient()
	}
	return n.client
}

func (n *Neutron) Call(req client.Request) *fasthttp.Response {
    resp := fasthttp.AcquireResponse()
	err := n.client.Call(context.Background(), req, resp)
	if err != nil {
		panic(err)
	}
	return resp
}