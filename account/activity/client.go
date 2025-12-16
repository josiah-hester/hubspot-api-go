// Package activity specifies the client methods for the HubSpot Account Activity API
//
// # This API package can be used to get details about the account and the private-app API usages
//
// Scope Requirements:
// account-info.security.read
package activity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

type Client struct {
	apiClient *client.Client
}

// NewClient creates a new accounts client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// RetrieveAuditLogs Retrieve activity history for user actions related to approvals, content updates, CRM object updates, security activity, and more (Enterprise only). Learn more about activities included in audit log exports.
func (c *Client) RetrieveAuditLogs(ctx context.Context, opts ...AccountActivityOption) ([]AuditLog, *Paging, error) {
	req := client.NewRequest("GET", "/account-info/v3/activity/audit-logs")
	req.WithContext(ctx)
	req.WithResourceType("accounts")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve audit logs: %w", err)
	}

	var response struct {
		Results []AuditLog `json:"results"`
		Paging  *Paging    `json:"paging"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal audit logs: %w", err)
	}

	return response.Results, response.Paging, nil
}

// RetrieveLoginActivity Retrieve logs of user actions related to login activity.
func (c *Client) RetrieveLoginActivity(ctx context.Context, opts ...AccountActivityOption) ([]LoginActivity, *Paging, error) {
	req := client.NewRequest("GET", "/account-info/v3/activity/login")
	req.WithContext(ctx)
	req.WithResourceType("accounts")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve login activity: %w", err)
	}

	var response struct {
		Results []LoginActivity `json:"results"`
		Paging  *Paging         `json:"paging"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal login activity: %w", err)
	}

	return response.Results, response.Paging, nil
}

// RetrieveSecurityHistory Retrieve logs of user actions related to security activity.
func (c *Client) RetrieveSecurityHistory(ctx context.Context, opts ...AccountActivityOption) ([]SecurityHistory, *Paging, error) {
	req := client.NewRequest("GET", "/account-info/v3/activity/security")
	req.WithContext(ctx)
	req.WithResourceType("accounts")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve security history: %w", err)
	}

	var response struct {
		Results []SecurityHistory `json:"results"`
		Paging  *Paging           `json:"paging"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal security history: %w", err)
	}

	return response.Results, response.Paging, nil
}
