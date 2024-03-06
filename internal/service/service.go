package service

import (
	"github.com/valyala/fasthttp"
	"go-openstackclient/internal/client"
)

type Service interface {
	HttpPrefix()                  string
	SupportedResources()          map[string]struct{}
	Call(req client.Request)      *fasthttp.Response
	Client()                      client.Client
}
