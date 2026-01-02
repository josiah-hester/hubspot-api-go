// Package contacts specifies the client methods for the HubSpot CRM Contacts API
package contacts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

type Client struct {
	apiClient *client.Client
}

// NewClient creates a new contacts client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

func (c *Client) GetContact(ctx context.Context, contactID string, opts ...GetContactOption) (*Contact, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/objects/contacts/%s", contactID))
	req.WithContext(ctx)
	req.WithResourceType("contacts")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseContactError(err, contactID)
	}

	var contact ContactResponse
	if err := json.Unmarshal(resp.Body, &contact); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact response: %w", err)
	}

	return &Contact{
		ID:         contact.ID,
		Properties: contact.Properties,
		Archived:   contact.Archived,
	}, nil
}

func (c *Client) CreateContact(ctx context.Context, input *CreateContactInput) (*Contact, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/contacts")
	req.WithContext(ctx)
	req.WithResourceType("contacts")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseContactError(err, "")
	}

	var contact ContactResponse
	if err := json.Unmarshal(resp.Body, &contact); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact response: %w", err)
	}

	return &Contact{
		ID:         contact.ID,
		Properties: contact.Properties,
		Archived:   contact.Archived,
	}, nil
}

func (c *Client) UpdateContact(ctx context.Context, contactID string, input *UpdateContactInput) (*Contact, error) {
	req := client.NewRequest("PATCH", fmt.Sprintf("/crm/v3/objects/contacts/%s", contactID))
	req.WithContext(ctx)
	req.WithResourceType("contacts")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseContactError(err, contactID)
	}

	var contact ContactResponse
	if err := json.Unmarshal(resp.Body, &contact); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact response: %w", err)
	}

	return &Contact{
		ID:         contact.ID,
		Properties: contact.Properties,
		Archived:   contact.Archived,
	}, nil
}

func (c *Client) DeleteContact(ctx context.Context, contactID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/objects/contacts/%s", contactID))
	req.WithContext(ctx)
	req.WithResourceType("contacts")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return ParseContactError(err, contactID)
	}

	return nil
}

func (c *Client) ListContacts(ctx context.Context, opts ...ListContactsOption) ([]Contact, string, error) {
	req := client.NewRequest("GET", "/crm/v3/objects/contacts")
	req.WithContext(ctx)
	req.WithResourceType("contacts")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, "", err
	}

	var listResp struct {
		Results []Contact `json:"results"`
		Paging  struct {
			Next struct {
				After string `json:"after"`
				Link  string `json:"link"`
			} `json:"next"`
		} `json:"paging"`
	}

	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal contacts list response: %w", err)
	}

	return listResp.Results, listResp.Paging.Next.After, nil
}

func (c *Client) SearchContacts(ctx context.Context, input *SearchContactsInput) (*SearchContactsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/contacts/search")
	req.WithContext(ctx)
	req.WithResourceType("contacts")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var searchResp SearchContactsResponse
	if err := json.Unmarshal(resp.Body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &searchResp, nil
}
