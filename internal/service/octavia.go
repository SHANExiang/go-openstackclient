package service

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
)

var supportedOctaviaResourceTypes = map[string]struct{}{
	consts.LOADBALANCER: struct{}{}, consts.LISTENER: struct{}{}, consts.POOL: struct{}{},
	consts.MEMBER: struct{}{}, consts.HEALTHMONITOR: struct{}{}, consts.L7POLICY: struct{}{},
	consts.L7RULE: struct{}{},
}

//func initOctaviaOutputChannels() map[string]chan Output {
//	outputChannel := make(map[string]chan Output)
//	for _, resourceType := range supportedOctaviaResourceTypes {
//		outputChannel[resourceType] = make(chan Output, 0)
//	}
//	return outputChannel
//}


type Octavia struct {
	client          client.Client
	httpPrefix      string
}

func newOctavia() *Octavia {
	return &Octavia{
		client: client.NewClient(),
		httpPrefix: fmt.Sprintf("http://%s:%d/v2.0/", configs.CONF.Host, consts.OctaviaPort),
	}
}

func (o *Octavia) SupportedResources() map[string]struct{} {
	return supportedOctaviaResourceTypes
}

func (o *Octavia) HttpPrefix() string {
	return fmt.Sprintf("http://%s:%d/v2.0/", configs.CONF.Host, consts.OctaviaPort)
}

func (o *Octavia) Client() client.Client {
	if o.client == nil {
		o.client = client.NewClient()
	}
	return o.client
}

func (o *Octavia) Call(req client.Request) *fasthttp.Response {
	resp := fasthttp.AcquireResponse()
	err := o.client.Call(context.Background(), req, resp)
	if err != nil {
		panic(err)
	}
	return resp
}
