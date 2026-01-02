package contacts

type FilterOperator string

const (
	EQ               FilterOperator = "EQ"
	NEQ              FilterOperator = "NEQ"
	LT               FilterOperator = "LT"
	LTE              FilterOperator = "LTE"
	GT               FilterOperator = "GT"
	GTE              FilterOperator = "GTE"
	Between          FilterOperator = "BETWEEN"
	In               FilterOperator = "IN"
	NotIn            FilterOperator = "NOT_IN"
	HasProperty      FilterOperator = "HAS_PROPERTY"
	NotHasProperty   FilterOperator = "NOT_HAS_PROPERTY"
	ContainsToken    FilterOperator = "CONTAINS_TOKEN"
	NotContainsToken FilterOperator = "NOT_CONTAINS_TOKEN"
)

// Contact represents a HubSpot contact object
type Contact struct {
	ID         string         `json:"id"`
	Properties map[string]any `json:"properties"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
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

// Paging represents pagination information
type Paging struct {
	Next *PagingLink `json:"next"`
	Prev *PagingLink `json:"prev"`
}

// PagingLink represents a pagination link
type PagingLink struct {
	After string `json:"after"`
	Link  string `json:"link"`
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

// SearchContactsInput represents input for searching contacts
type SearchContactsInput struct {
	FilterGroups []FilterGroup `json:"filterGroups"`
	Sorts        []string      `json:"sorts"`
	Query        string        `json:"query"`
	Properties   []string      `json:"properties"`
	Limit        int           `json:"limit"`
	After        string        `json:"after"`
}

// FilterGroup represents a group of filters
type FilterGroup struct {
	Filters []Filter `json:"filters"`
}

// Filter represents a single filter
type Filter struct {
	PropertyName string         `json:"propertyName"`
	Operator     FilterOperator `json:"operator"`
	HighValue    string         `json:"highValue,omitempty"`
	Value        string         `json:"value,omitempty"`
	Values       []string       `json:"values,omitempty"`
}

// SearchContactsResponse represents a response from search
type SearchContactsResponse struct {
	Total   int       `json:"total"`
	Results []Contact `json:"results"`
	Paging  *Paging   `json:"paging"`
}
