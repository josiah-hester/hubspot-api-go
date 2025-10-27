package client

import "context"

type Request struct {
	Method      string
	Path        string
	Body        any
	QueryParams map[string]string
	Headers     map[string]string

	// Metadata for middleware
	ResourceType string
	RetryCount   int

	// Context for timeouts/cancellation
	Context context.Context
}

func NewRequest(method, path string) *Request {
	return &Request{
		Method:      method,
		Path:        path,
		QueryParams: make(map[string]string),
		Headers:     make(map[string]string),
		Context:     context.Background(),
	}
}

func (r *Request) WithContext(ctx context.Context) *Request {
	r.Context = ctx
	return r
}

func (r *Request) WithResourceType(resourceType string) *Request {
	r.ResourceType = resourceType
	return r
}

func (r *Request) WithBody(body any) *Request {
	r.Body = body
	return r
}

func (r *Request) AddQueryParam(key, value string) *Request {
	r.QueryParams[key] = value
	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	r.Headers[key] = value
	return r
}
