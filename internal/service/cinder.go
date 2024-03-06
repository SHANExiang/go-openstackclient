package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
	"go-openstackclient/internal/entity"
	"log"
	"math"
	"strconv"
	"strings"
)

var maxEfficientQOS string
var maxSSDEfficientQOS string
var supportedCinderResourceTypes = map[string]struct{}{
	consts.VOLUME: struct{}{}, consts.SNAPSHOT: struct{}{},
}

type Cinder struct {
	client          client.Client
	httpPrefix      string
}

func newCinder() *Cinder {
	return &Cinder{
		client: client.NewClient(),
		httpPrefix: fmt.Sprintf("http://%s:%d/v3/", configs.CONF.Host, consts.CINDER),
	}
}

func (c *Cinder) SupportedResources() map[string]struct{} {
	return supportedCinderResourceTypes
}


func (c *Cinder) HttpPrefix() string {
	return fmt.Sprintf("http://%s:%d/v3", configs.CONF.Host, consts.CinderPort)
}

func (c *Cinder) Client() client.Client {
	if c.client == nil {
		c.client = client.NewClient()
	}
	return c.client
}

func (c *Cinder) Call(req client.Request) *fasthttp.Response {
	resp := fasthttp.AcquireResponse()
	req.Headers()["Openstack-Api-Version"] = "volume 3.59"
	err := c.client.Call(context.Background(), req, resp)
	if err != nil {
		panic(err)
	}
	return resp
}
