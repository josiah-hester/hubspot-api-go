// Package objects specifies the client methods for the HubSpot CRM Objects API
package objects

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aacc-dev/go-hubspot-sdk/client"
	"github.com/aacc-dev/go-hubspot-sdk/internal/tools"
)

type Client struct {
	apiClient *client.Client
}

// NewClient creates a new objects client
func NewClient(apiClient *client.Client) *Client {
	return &Client{
		apiClient: apiClient,
	}
}

// -------- Basic Methods --------

// ListObjects returns a list of HubSpot objects
//
// opts:
// WithLimit
// WithAfter
// WithProperties
// WithPropertiesWithHistory
// WithAssociations
// WithArchived
func (c *Client) ListObjects(ctx context.Context, objectType string, opts ...ObjectsOption) ([]Object, string, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/objects/%s", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, "", ParseObjectError(err, objectType)
	}

	var objResp ListObjectsResponse
	if err := json.Unmarshal(resp.Body, &objResp); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal object response: %w", err)
	}

	if len(objResp.Results) > 0 {
		return objResp.Results, objResp.Paging.Next.After, nil
	}

	return nil, "", fmt.Errorf("no objects found for type %s", objectType)
}

// CreateObject creates a new HubSpot object
func (c *Client) CreateObject(ctx context.Context, input *CreateObjectInput, objectType string) (*Object, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var object CreateObjectResponse
	if err := json.Unmarshal(resp.Body, &object); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object response: %w", err)
	}

	return &object.Entity, nil
}

// ReadObject reads a HubSpot object by id or specified idProperty
//
// opts:
// WithProperties
// WithPropertiesWithHistory
// WithAssociations
// WithArchived
// WithIDProperty
func (c *Client) ReadObject(ctx context.Context, objectType string, id string, opts ...ObjectsOption) (*Object, error) {
	req := client.NewRequest("GET", fmt.Sprintf("/crm/v3/objects/%s/%s", objectType, id))
	req.WithContext(ctx)
	req.WithResourceType("objects")

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj Object
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object response: %w", err)
	}

	return &obj, nil
}

// UpdateObject updates a HubSpot object by id or specified idProperty
//
// opts:
// WithIDProperty
func (c *Client) UpdateObject(ctx context.Context, objectType string, id string, input *UpdateObjectInput, opts ...ObjectsOption) (*Object, error) {
	req := client.NewRequest("PATCH", fmt.Sprintf("/crm/v3/objects/%s/%s", objectType, id))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj Object
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object response: %w", err)
	}

	return &obj, nil
}

// ArchiveObject archives a HubSpot object by id
func (c *Client) ArchiveObject(ctx context.Context, objectType string, id string) error {
	req := client.NewRequest("DELETE", fmt.Sprintf("/crm/v3/objects/%s/%s", objectType, id))
	req.WithContext(ctx)
	req.WithResourceType("objects")

	_, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return ParseObjectError(err, objectType)
	}

	return nil
}

// MergeObjects merges two HubSpot objects by id
func (c *Client) MergeObjects(ctx context.Context, objectType string, input *MergeObjectsInput) (*Object, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s/merge", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj Object
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object response: %w", err)
	}

	return &obj, nil
}

// -------- Batch Methods --------

// BatchReadObjects reads a batch of HubSpot objects by id or unique idProperty
//
// opts:
// WithArchived
func (c *Client) BatchReadObjects(ctx context.Context, objectType string, input *BatchReadObjectsInput, opts ...ObjectsOption) (*BatchResponse, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s/batch/read", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj BatchResponse
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object response: %w", err)
	}

	if len(obj.Errors) > 0 {
		for _, err := range obj.Errors {
			fmt.Printf("some errors occurred in the batch request: %v", ParseObjectError(&err, objectType))
		}
	}

	return &obj, nil
}

// BatchCreateObjects creates a batch of HubSpot objects
func (c *Client) BatchCreateObjects(ctx context.Context, objectType string, input *BatchCreateObjectsInput) (*BatchResponse, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s/batch/create", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj BatchResponse
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to ubmarshal object response: %w", err)
	}

	if len(obj.Errors) > 0 {
		for _, err := range obj.Errors {
			fmt.Printf("some errors occurred in the batch request: %v", ParseObjectError(&err, objectType))
		}
	}

	return &obj, nil
}

// BatchUpdateObjects updates a batch of HubSpot objects
func (c *Client) BatchUpdateObjects(ctx context.Context, objectType string, input *BatchUpdateObjectsInput) (*BatchResponse, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s/batch/update", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj BatchResponse
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to ubmarshal object response: %w", err)
	}

	if len(obj.Errors) > 0 {
		for _, err := range obj.Errors {
			fmt.Printf("some errors occurred in the batch request: %v", ParseObjectError(&err, objectType))
		}
	}

	return &obj, nil
}

// BatchCreateOrUpdateObjects creates or updates a batch of HubSpot objects
func (c *Client) BatchCreateOrUpdateObjects(ctx context.Context, objectType string, input *BatchCreateOrUpdateObjectsInput) (*BatchResponse, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s/batch/upsert", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj BatchResponse
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to ubmarshal object response: %w", err)
	}

	if len(obj.Errors) > 0 {
		for _, err := range obj.Errors {
			fmt.Printf("some errors occurred in the batch request: %v", ParseObjectError(&err, objectType))
		}
	}

	return &obj, nil
}

// BatchArchiveObjects archives a batch of HubSpot objects
func (c *Client) BatchArchiveObjects(ctx context.Context, objectType string, input *BatchArchiveObjectsInput) (*BatchResponse, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s/batch/archive", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj BatchResponse
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		return nil, fmt.Errorf("failed to ubmarshal object response: %w", err)
	}

	if len(obj.Errors) > 0 {
		for _, err := range obj.Errors {
			fmt.Printf("some errors occurred in the batch request: %v", ParseObjectError(&err, objectType))
		}
	}
	return &obj, nil
}

// -------- Search Methods --------

// SearchObjects searches for HubSpot objects
func (c *Client) SearchObjects(ctx context.Context, objectType string, input *SearchObjectsInput) (*SearchObjectsResponse, error) {
	req := client.NewRequest("POST", fmt.Sprintf("/crm/v3/objects/%s/search", objectType))
	req.WithContext(ctx)
	req.WithResourceType("objects")
	req.WithBody(input)

	resp, err := c.apiClient.Do(ctx, req)
	if err != nil {
		return nil, ParseObjectError(err, objectType)
	}

	var obj SearchObjectsResponse
	if err := tools.NewRequiredTagStruct(obj).UnmarhsalJSON(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to ubmarshal object response: %w", err)
	}

	if len(obj.Results) == 0 {
		return nil, fmt.Errorf("no results found in search request")
	}

	return &obj, nil
}
