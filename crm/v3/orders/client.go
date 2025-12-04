// Package orders provides client methods for the HubSpot CRM Orders API
package orders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// Client represents the Orders API client
type Client struct {
	apiClient *client.Client
}

// NewClient creates a new orders client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// CreateOrder creates a new order
func (c *Client) CreateOrder(ctx context.Context, input *CreateOrderInput) (*Order, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/orders")
	req.WithContext(ctx)
	req.WithResourceType("orders")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp.Body, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	return &order, nil
}

// GetOrder retrieves an order by ID
func (c *Client) GetOrder(ctx context.Context, orderID string, opts ...OrderOption) (*Order, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/objects/orders/%s", orderID))
	req.WithContext(ctx)
	req.WithResourceType("orders")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp.Body, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	return &order, nil
}

// UpdateOrder updates an order
func (c *Client) UpdateOrder(ctx context.Context, orderID string, input *UpdateOrderInput) (*Order, error) {
	req := client.NewRequest("PATCH", fmt.Sprintf("/crm/v3/objects/orders/%s", orderID))
	req.WithContext(ctx)
	req.WithResourceType("orders")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp.Body, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order response: %w", err)
	}

	return &order, nil
}

// ArchiveOrder archives (deletes) an order
func (c *Client) ArchiveOrder(ctx context.Context, orderID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/objects/orders/%s", orderID))
	req.WithContext(ctx)
	req.WithResourceType("orders")

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// ListOrders lists orders with optional filters
func (c *Client) ListOrders(ctx context.Context, opts ...OrderOption) (*ListOrdersResponse, error) {
	req := client.NewRequest("GET", "/crm/v3/objects/orders")
	req.WithContext(ctx)
	req.WithResourceType("orders")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var listResp ListOrdersResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal orders list response: %w", err)
	}

	return &listResp, nil
}

// BatchReadOrders retrieves multiple orders by ID
func (c *Client) BatchReadOrders(ctx context.Context, input *BatchReadOrdersInput) (*BatchOrdersResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/orders/batch/read")
	req.WithContext(ctx)
	req.WithResourceType("orders")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchOrdersResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchCreateOrders creates multiple orders
func (c *Client) BatchCreateOrders(ctx context.Context, input *BatchCreateOrdersInput) (*BatchOrdersResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/orders/batch/create")
	req.WithContext(ctx)
	req.WithResourceType("orders")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchOrdersResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchUpdateOrders updates multiple orders
func (c *Client) BatchUpdateOrders(ctx context.Context, input *BatchUpdateOrdersInput) (*BatchOrdersResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/orders/batch/update")
	req.WithContext(ctx)
	req.WithResourceType("orders")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var batchResp BatchOrdersResponse
	if err := json.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
	}

	return &batchResp, nil
}

// BatchArchiveOrders archives multiple orders
func (c *Client) BatchArchiveOrders(ctx context.Context, input *BatchArchiveOrdersInput) error {
	req := client.NewRequest("POST", "/crm/v3/objects/orders/batch/archive")
	req.WithContext(ctx)
	req.WithResourceType("orders")
	req.WithBody(input)

	_, err := c.apiClient.Do(ctx, req)
	return err
}

// SearchOrders searches for orders
func (c *Client) SearchOrders(ctx context.Context, input *SearchOrdersInput) (*SearchOrdersResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/orders/search")
	req.WithContext(ctx)
	req.WithResourceType("orders")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var searchResp SearchOrdersResponse
	if err := json.Unmarshal(resp.Body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &searchResp, nil
}
