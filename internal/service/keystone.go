package service

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
)

var supportedKeystoneResourceTypes = map[string]struct{}{
	consts.PROJECT: struct{}{}, consts.USER: struct {}{},
	consts.TOKEN: struct{}{},
}

type Keystone struct {
	client          client.Client
	httpPrefix      string
}

func newKeystone() *Keystone {
	return &Keystone{
		client: client.NewClient(),
		httpPrefix: fmt.Sprintf("http://%s:%d/v3/", configs.CONF.Host, consts.KeystonePort),
	}
}

func (k *Keystone) SupportedResources() map[string]struct{} {
	return supportedKeystoneResourceTypes
}

func (k *Keystone) HttpPrefix() string {
	return fmt.Sprintf("http://%s:%d/v3", configs.CONF.Host, consts.KeystonePort)
}

func (k *Keystone) Client() client.Client {
	if k.client == nil {
		k.client = client.NewClient()
	}
	return k.client
}

func (k *Keystone) Call(req client.Request) *fasthttp.Response {
    resp := fasthttp.AcquireResponse()
	err := k.client.Call(context.Background(), req, resp)
	if err != nil {
		panic(err)
	}
	return resp
}