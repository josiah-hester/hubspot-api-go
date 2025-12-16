// Package appflags specifies the client methods for the HubSpot App Management App Flags API
package appflags

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

type Client struct {
	apiClient *client.Client
}

// NewClient creates a new app flags client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// RetrieveAppFeatureFlags Retrieve the current status of the app’s feature flags. No request body is included.
func (c *Client) RetrieveAppFeatureFlags(ctx context.Context, appID int, flagName string) (*FlagInfo, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/app-management/v3/apps/%d/feature-flags/%s", appID, flagName))
	req.WithContext(ctx)
	req.WithResourceType("feature-flags")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve app feature flags: %w", err)
	}

	var flagInfo FlagInfo
	if err := json.Unmarshal(resp.Body, &flagInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal app feature flags: %w", err)
	}

	return &flagInfo, nil
}

// RetrieveAccountsWithSetFlagState Retrieve a list of HubSpot accounts with an account-level flag setting for the specified app. No request body is included.
func (c *Client) RetrieveAccountsWithSetFlagState(ctx context.Context, appID int, flagName string, opts ...AppFlagOption) ([]FlagState, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/feature-flags/v3/%d/flags/%s/portals", appID, flagName))
	req.WithContext(ctx)
	req.WithResourceType("feature-flags")

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve app feature flags: %w", err)
	}

	var response struct {
		PortalFlagStates []FlagState `json:"portalFlagStates"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal app feature flags: %w", err)
	}

	return response.PortalFlagStates, nil
}
