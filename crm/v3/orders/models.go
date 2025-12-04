package orders

// Order represents a HubSpot order object
type Order struct {
	ID                    string                           `json:"id"`
	Properties            map[string]string                `json:"properties"`
	PropertiesWithHistory map[string][]PropertyWithHistory `json:"propertiesWithHistory"`
	CreatedAt             string                           `json:"createdAt"`
	UpdatedAt             string                           `json:"updatedAt"`
	Archived              bool                             `json:"archived"`
	ArchivedAt            string                           `json:"archivedAt"`
}

// PropertyWithHistory represents a property with its historical values
type PropertyWithHistory struct {
	Value           string `json:"value"`
	Timestamp       string `json:"timestamp"`
	SourceType      string `json:"sourceType"`
	SourceID        string `json:"sourceId"`
	SourceLabel     string `json:"sourceLabel"`
	UpdatedByUserID int    `json:"updatedByUserId"`
}

// CreateOrderInput represents the input for creating an order
type CreateOrderInput struct {
	Properties map[string]string `json:"properties"`
}

// UpdateOrderInput represents the input for updating an order
type UpdateOrderInput struct {
	Properties map[string]string `json:"properties"`
}

// ListOrdersResponse represents the response from listing orders
type ListOrdersResponse struct {
	Results []Order `json:"results"`
	Paging  *Paging `json:"paging"`
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

// BatchReadOrdersInput represents input for batch read
type BatchReadOrdersInput struct {
	Properties            []string `json:"properties"`
	PropertiesWithHistory []string `json:"propertiesWithHistory"`
	IDProperty            string   `json:"idProperty"`
	Inputs                []struct {
		ID string `json:"id"`
	} `json:"inputs"`
}

// BatchCreateOrdersInput represents input for batch create
type BatchCreateOrdersInput struct {
	Inputs []CreateOrderInput `json:"inputs"`
}

// BatchUpdateOrdersInput represents input for batch update
type BatchUpdateOrdersInput struct {
	Inputs []struct {
		ID         string            `json:"id"`
		Properties map[string]string `json:"properties"`
	} `json:"inputs"`
}

// BatchArchiveOrdersInput represents input for batch archive
type BatchArchiveOrdersInput struct {
	Inputs []struct {
		ID string `json:"id"`
	} `json:"inputs"`
}

// BatchOrdersResponse represents response from batch operations
type BatchOrdersResponse struct {
	Status      string  `json:"status"`
	Results     []Order `json:"results"`
	StartedAt   string  `json:"startedAt"`
	CompletedAt string  `json:"completedAt"`
}

// SearchOrdersInput represents input for searching orders
type SearchOrdersInput struct {
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
	PropertyName string `json:"propertyName"`
	Operator     string `json:"operator"`
	Value        any    `json:"value"`
}

// SearchOrdersResponse represents response from search
type SearchOrdersResponse struct {
	Total   int     `json:"total"`
	Results []Order `json:"results"`
	Paging  *Paging `json:"paging"`
}
