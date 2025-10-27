package v4

import "github.com/aacc-dev/go-hubspot-sdk/client"

// ClientV4 provides access to CRM v4 API endpoints
type ClientV4 struct {
	client *client.Client
}

// NewClientV4 creates a new v4 CRM client
func NewClientV4(client *client.Client) *ClientV4 {
	return &ClientV4{
		client: client,
	}
}
