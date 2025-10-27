package contacts

import "time"

// Contact represents a HubSpot contact object
type Contact struct {
	ID         string         `json:"id"`
	Properties map[string]any `json:"properties"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Archived   bool           `json:"archived"`
}

// ContactResponse is the API response format for contacts
type ContactResponse struct {
	ID         string         `json:"id"`
	Properties map[string]any `json:"properties"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
	Archived   bool           `json:"archived"`
}

// CreateContactInput is the input for creating a contact
type CreateContactInput struct {
	Properties map[string]string `json:"properties"`
}

// UpdateContactInput is the input for updating a contact
type UpdateContactInput struct {
	Properties map[string]string `json:"properties"`
}

// BatchCreateContactsInput is the input for batch creating contacts
type BatchCreateContactsInput struct {
	Inputs []struct {
		Properties map[string]string `json:"properties"`
	} `json:"inputs"`
}

// BatchReadContactsInput is the input for batch reading contacts
type BatchReadContactsInput struct {
	PropertiesToFetch []string `json:"propertiesToFetch"`
	Inputs            []struct {
		ID string `json:"id"`
	} `json:"inputs"`
}

