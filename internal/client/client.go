package client

import (
    "context"
)

type Request interface {
    Method()           string
    Endpoint()         string
    ContentType()      string
    Body()             interface{}
    Headers()          map[string]string
}

type Client interface {
    NewRequest(endpoint, method, contextType string, headers map[string]string, body interface{}) Request
    Call(ctx context.Context, req Request, resp interface{}) error
    String()  string
}

func NewClient() Client {
    return &RestClient{
        client: newRestClient(),
    }
}
