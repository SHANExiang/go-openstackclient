package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

var _ Client = (*RestClient)(nil)

type RestClient struct {
	client           *fasthttp.Client
}

func newRestClient() (client *fasthttp.Client) {
	readTimeout, _ := time.ParseDuration("5000000ms")
	writeTimeout, _ := time.ParseDuration("5000000ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	client = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
	return client
}

func (r *RestClient) NewRequest(endpoint, method, contentType string, headers map[string]string, req interface{}) Request {
	return newRequest(endpoint, method, contentType, headers, req)
}

func (r *RestClient) Call(ctx context.Context, req Request, resp interface{}) error {
	res := resp.(*fasthttp.Response)
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.SetRequestURI(req.Endpoint())
	request.Header.SetContentType(req.ContentType())
	request.Header.SetMethod(req.Method())
	if req.Body() != nil {
		request.SetBody([]byte(req.Body().(string)))
	}
	for k, v := range req.Headers() {
		request.Header.Set(k, v)
	}

	log.Printf("Starting to %s request %s, body==%+v", req.Method(), req.Endpoint(), req.Body())
	err := r.client.Do(request, res)
	if err != nil {
		errorInfo := fmt.Sprintf("Failed to make request %s", err)
		log.Println(errorInfo)
		return errors.New(errorInfo)
	}
	if res.StatusCode() > 204 {
		errorInfo := fmt.Sprintf("%s request %s failed, resp==%s", req.Method(), req.Endpoint(), res.Body())
		log.Printf(errorInfo)
		return errors.New(errorInfo)
	}
	return nil
}

func (r *RestClient) String() string {
    return "rest"
}




