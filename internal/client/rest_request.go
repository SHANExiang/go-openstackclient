package client

var _ Request = (*RestRequest)(nil)

type RestRequest struct {
	method      string
	endpoint    string
	contentType string
	body        interface{}
	headers     map[string]string
}

func newRequest(endpoint, method string, contentType string, headers map[string]string, body interface{}) Request {
	return &RestRequest{
		method:      method,
		endpoint:    endpoint,
		body:        body,
		contentType: contentType,
		headers: headers,
	}
}

func (r *RestRequest) Method() string {
	return r.method
}

func (r *RestRequest) Endpoint() string {
	return r.endpoint
}

func (r *RestRequest) ContentType() string {
	return r.contentType
}

func (r *RestRequest) Body() interface{} {
	return r.body
}

func (r *RestRequest) Headers() map[string]string {
	return r.headers
}
