package v3

import "github.com/aacc-dev/go-hubspot-sdk/client"

// ClientV3 provides access to CRM v3 API endpoints
type ClientV3 struct {
	client *client.Client
}

// NewClientV3 creates a new v3 CRM client
func NewClientV3(client *client.Client) *ClientV3 {
	return &ClientV3{
		client: client,
	}
}
