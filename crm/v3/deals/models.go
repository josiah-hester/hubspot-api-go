package deals

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

// Deal represents a HubSpot deal object
type Deal struct {
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

// CreateDealInput represents the input for creating a deal
type CreateDealInput struct {
	Properties map[string]string `json:"properties"`
}

// UpdateDealInput represents the input for updating a deal
type UpdateDealInput struct {
	Properties map[string]string `json:"properties"`
}

// ListDealsResponse represents the response from listing deals
type ListDealsResponse struct {
	Results []Deal  `json:"results"`
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

// BatchReadDealsInput represents input for batch read
type BatchReadDealsInput struct {
	Properties            []string `json:"properties"`
	PropertiesWithHistory []string `json:"propertiesWithHistory"`
	IDProperty            string   `json:"idProperty"`
	Inputs                []struct {
		ID string `json:"id"`
	} `json:"inputs"`
}

// BatchCreateDealsInput represents input for batch create
type BatchCreateDealsInput struct {
	Inputs []CreateDealInput `json:"inputs"`
}

// BatchUpdateDealsInput represents input for batch update
type BatchUpdateDealsInput struct {
	Inputs []struct {
		ID         string            `json:"id"`
		Properties map[string]string `json:"properties"`
	} `json:"inputs"`
}

// BatchArchiveDealsInput represents input for batch archive
type BatchArchiveDealsInput struct {
	Inputs []struct {
		ID string `json:"id"`
	} `json:"inputs"`
}

// BatchDealsResponse represents response from batch operations
type BatchDealsResponse struct {
	Status      string `json:"status"`
	Results     []Deal `json:"results"`
	StartedAt   string `json:"startedAt"`
	CompletedAt string `json:"completedAt"`
}

// SearchDealsInput represents input for searching deals
type SearchDealsInput struct {
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

// SearchDealsResponse represents response from search
type SearchDealsResponse struct {
	Total   int     `json:"total"`
	Results []Deal  `json:"results"`
	Paging  *Paging `json:"paging"`
}
