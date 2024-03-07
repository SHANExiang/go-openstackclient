package service

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
)


var supportedGlanceResourceTypes = map[string]struct{}{
	consts.Image: struct{}{},
}

type Glance struct {
	client          client.Client
	httpPrefix      string
}

func newGlance() *Glance {
	return &Glance{
		client: client.NewClient(),
		httpPrefix: fmt.Sprintf("http://%s:%d/v2", configs.CONF.Host, consts.GlancePort),
	}
}

func (g *Glance) SupportedResources() map[string]struct{} {
	return supportedGlanceResourceTypes
}

func (g *Glance) HttpPrefix() string {
	return fmt.Sprintf("http://%s:%d/v2", configs.CONF.Host, consts.GlancePort)
}

func (g *Glance) Client() client.Client {
	if g.client == nil {
		g.client = client.NewClient()
	}
	return g.client
}

func (g *Glance) Call(req client.Request) *fasthttp.Response {
	resp := fasthttp.AcquireResponse()
	err := g.client.Call(context.Background(), req, resp)
	if err != nil {
		panic(err)
	}
	return resp
}
