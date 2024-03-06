package service

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
)

var supportedNovaResourceTypes = map[string]struct{}{
	consts.SERVER: struct{}{},
}

type Nova struct {
	client          client.Client
	httpPrefix      string
}

func newNova() *Nova {
	return &Nova{
		client: client.NewClient(),
		httpPrefix: fmt.Sprintf("http://%s:%d/v2.0/", configs.CONF.Host, consts.NovaPort),
	}
}

func (n *Nova) SupportedResources() map[string]struct{} {
	return supportedNovaResourceTypes
}


func (n *Nova) HttpPrefix() string {
	return fmt.Sprintf("http://%s:%d/v2.1", configs.CONF.Host, consts.NovaPort)
}

func (n *Nova) Client() client.Client {
	if n.client == nil {
		n.client = client.NewClient()
	}
	return n.client
}

func (n *Nova) Call(req client.Request) *fasthttp.Response {
    resp := fasthttp.AcquireResponse()
    req.Headers()["OpenStack-API-Version"] = "compute 2.74"
	err := n.client.Call(context.Background(), req, resp)
	if err != nil {
		panic(err)
	}
	return resp
}