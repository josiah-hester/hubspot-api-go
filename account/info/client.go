// Package info specifies the client methods for the HubSpot Accounts Info API
//
// # This API package can be used to get details about the account and the private-app API usages
//
// Scope Requirements:
// oauth
package info

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

func (c *Client) GetAccountDetails(ctx context.Context) (*AccountDetails, error) {
	req := client.NewRequest("GET", "/account-info/v3/details")
	req.WithContext(ctx)
	req.WithResourceType("accounts")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var details AccountDetails
	if err := json.Unmarshal(resp.Body, &details); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account details response: %w", err)
	}

	return &details, nil
}

// RetrievePrivateAppDailyAPIUsage returns the daily rate limits and usage for legacy private-apps
func (c *Client) RetrievePrivateAppDailyAPIUsage(ctx context.Context) ([]PrivateAppAPIUsage, error) {
	req := client.NewRequest("GET", "/account-info/v3/api-usage/daily/private-apps")
	req.WithContext(ctx)
	req.WithResourceType("accounts")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var results struct {
		Results []PrivateAppAPIUsage `json:"results"`
	}
	if err := json.Unmarshal(resp.Body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal daily rate limits and usage for legacy private-apps: %w", err)
	}

	return results.Results, nil
}
