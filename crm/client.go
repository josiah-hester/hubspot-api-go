package crm

import (
	"github.com/aacc-dev/go-hubspot-sdk/client"
	v3 "github.com/aacc-dev/go-hubspot-sdk/crm/v3"
	v4 "github.com/aacc-dev/go-hubspot-sdk/crm/v4"
)

// CRMClient provides access to all CRM API versions
type CRMClient struct {
	V3 *v3.ClientV3
	V4 *v4.ClientV4
}

// NewCRMClient creates a new CRM client
func NewCRMClient(apiClient *client.Client) *CRMClient {
	return &CRMClient{
		V3: v3.NewClientV3(apiClient),
		V4: v4.NewClientV4(apiClient),
	}
}
