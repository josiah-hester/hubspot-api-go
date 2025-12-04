// Package tickets specifies the client methods for the HubSpot CRM Tickets API
package tickets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/josiah-hester/go-hubspot-sdk/internal/tools"
)

type Client struct {
	apiClient *client.Client
}

// NewClient creates a new tickets client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// -------- Basic Methods --------

// ListTickets returns a list of tickets
//
// opts:
// WithLimit
// WithAfter
// WithProperties
// WithPropertiesWithHistory
// WithAssociations
// WithArchived
func (c *Client) ListTickets(ctx context.Context, opts ...TicketOption) (*ListTicketsResponse, error) {
	req := client.NewRequest("GET", "/crm/v3/objects/tickets")
	req.WithContext(ctx)
	req.WithResourceType("tickets")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var tickets ListTicketsResponse
	if err := json.Unmarshal(resp.Body, &tickets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tickets response: %w", err)
	}

	if len(tickets.Results) == 0 {
		return nil, fmt.Errorf("no tickets found")
	}

	return &tickets, nil
}

// CreateTicket creates a new ticket
func (c *Client) CreateTicket(ctx context.Context, input *CreateTicketInput) (*CreateTicketResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var ticket CreateTicketResponse
	if err := json.Unmarshal(resp.Body, &ticket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket response: %w", err)
	}

	return &ticket, nil
}

// ReadTicket returns a single ticket
//
// opts:
// WithProperties
// WithPropertiesWithHistory
// WithAssociations
// WithArchived
// WithIDProperty
func (c *Client) ReadTicket(ctx context.Context, ticketID string, opts ...TicketOption) (*Ticket, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/objects/tickets/%s", ticketID))
	req.WithContext(ctx)
	req.WithResourceType("tickets")

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var ticket Ticket
	if err := json.Unmarshal(resp.Body, &ticket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket response: %w", err)
	}

	return &ticket, nil
}

// UpdateTicket updates a ticket
//
// opts:
// WithIDProperty
func (c *Client) UpdateTicket(ctx context.Context, ticketID string, input *UpdateTicketInput, opts ...TicketOption) (*Ticket, error) {
	req := client.NewRequest("PATCH", fmt.Sprintf("/crm/v3/objects/tickets/%s", ticketID))
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var ticket Ticket
	if err := json.Unmarshal(resp.Body, &ticket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket response: %w", err)
	}

	return &ticket, nil
}

// ArchiveTicket archives a ticket
func (c *Client) ArchiveTicket(ctx context.Context, ticketID string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/objects/tickets/%s", ticketID))
	req.WithContext(ctx)
	req.WithResourceType("tickets")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// MergeTwoTickets merges two tickets together
func (c *Client) MergeTwoTickets(ctx context.Context, input *MergeTwoTicketsInput) error {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets/merge")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return err
	}

	var ticket Ticket
	if err := json.Unmarshal(resp.Body, &ticket); err != nil {
		return fmt.Errorf("failed to unmarshal ticket response: %w", err)
	}

	return nil
}

// -------- Batch Methods --------

// BatchReadTickets Retrieve a batch of tickets by ID (ticketId) or unique property value (idProperty)
//
// opts:
// WithArchived
func (c *Client) BatchReadTickets(ctx context.Context, input *BatchReadTicketsInput, opts ...TicketOption) (*BatchTicketsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets/batch/read")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var obj BatchTicketsResponse
	if err := tools.NewRequiredTagStruct(&obj).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch tickets response: %w", err)
	}

	return &obj, nil
}

// BatchCreateTickets Create a batch of tickets. The inputs array can contain a properties object to define property values for the ticket, along with an associations array to define associations with other CRM records.
func (c *Client) BatchCreateTickets(ctx context.Context, input *BatchCreateTicketsInput) (*BatchTicketsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets/batch/create")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var obj BatchTicketsResponse
	if err := tools.NewRequiredTagStruct(&obj).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch tickets response: %w", err)
	}

	return &obj, nil
}

// BatchUpdateTickets Update a batch of tickets by ID (ticketId) or unique property value (idProperty). Provided property values will be overwritten. Read-only and non-existent properties will result in an error. Properties values can be cleared by passing an empty string.
func (c *Client) BatchUpdateTickets(ctx context.Context, input *BatchUpdateTicketsInput) (*BatchTicketsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets/batch/update")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var obj BatchTicketsResponse
	if err := tools.NewRequiredTagStruct(&obj).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch tickets response: %w", err)
	}

	return &obj, nil
}

// BatchCreateOrUpdateTickets Create or update records identified by a unique property value as specified by the idProperty query param. idProperty query param refers to a property whose values are unique for the object.
func (c *Client) BatchCreateOrUpdateTickets(ctx context.Context, input *BatchCreateOrUpdateTicketsInput) (*BatchTicketsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets/batch/createOrUpdate")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var obj BatchTicketsResponse
	if err := tools.NewRequiredTagStruct(&obj).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch tickets response: %w", err)
	}

	return &obj, nil
}

// BatchArchiveTickets Delete a batch of tickets by ID. Deleted tickets can be restored within 90 days of deletion.
func (c *Client) BatchArchiveTickets(ctx context.Context, input *BatchArchiveTicketsInput) (*BatchTicketsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets/batch/archive")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var obj BatchTicketsResponse
	if err := tools.NewRequiredTagStruct(&obj).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch tickets response: %w", err)
	}

	return &obj, nil
}

// -------- Search Methods --------

// SearchTickets Search for tickets by filtering on properties, searching through associations, and sorting results.
func (c *Client) SearchTickets(ctx context.Context, input *SearchTicketsInput) (*SearchTicketsResponse, error) {
	req := client.NewRequest("POST", "/crm/v3/objects/tickets/search")
	req.WithContext(ctx)
	req.WithResourceType("tickets")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var search SearchTicketsResponse
	if err := tools.NewRequiredTagStruct(&search).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &search, nil
}
