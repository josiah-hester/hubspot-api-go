package tickets

type AssociationCategory string

const (
	HubspotDefined    AssociationCategory = "HUBSPOT_DEFINED"
	UserDefined       AssociationCategory = "USER_DEFINED"
	IntegratorDefined AssociationCategory = "INTEGRATOR_DEFINED"
)

type BatchStatus string

const (
	Pending    BatchStatus = "PENDING"
	Processing BatchStatus = "PROCESSING"
	Cancelled  BatchStatus = "CANCELLED"
	Complete   BatchStatus = "COMPLETE"
)

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

type Ticket struct {
	CreatedAt             string                           `json:"createdAt" required:"yes"`
	Archived              bool                             `json:"archived" required:"yes"`
	ID                    string                           `json:"id" required:"yes"`
	Properties            map[string]string                `json:"properties" required:"yes"`
	UpdatedAt             string                           `json:"updatedAt" required:"yes"`
	Associations          map[string]Association           `json:"associations"`
	ArchivedAt            string                           `json:"archivedAt"`
	PropertiesWithHistory map[string][]PropertyWithHistory `json:"propertiesWithHistory"`
	ObjectWriteTraceID    string                           `json:"objectWriteTraceId"`
}

type Association struct {
	Results []struct {
		ID   string `json:"id" required:"yes"`
		Type string `json:"type" required:"yes"`
	} `json:"results" required:"yes"`
	Paging struct {
		Next struct {
			After string `json:"after" required:"yes"`
			Link  string `json:"link"`
		} `json:"next"`
		Prev struct {
			Before string `json:"before" required:"yes"`
			Link   string `json:"link"`
		} `json:"prev"`
	} `json:"paging"`
}

type PropertyWithHistory struct {
	SourceType      string `json:"sourceType" required:"yes"`
	Value           string `json:"value" required:"yes"`
	Timestamp       string `json:"timestamp" required:"yes"`
	SourceID        string `json:"sourceId"`
	SourceLabel     string `json:"sourceLabel"`
	UpdatedByUserID int    `json:"updatedByUserId"`
}

type ListTicketsResponse struct {
	Results []Ticket `json:"results" required:"yes"`
	Paging  struct {
		Next struct {
			After string `json:"after" required:"yes"`
			Link  string `json:"link"`
		} `json:"next"`
		Prev struct {
			Before string `json:"before" required:"yes"`
			Link   string `json:"link"`
		} `json:"prev"`
	} `json:"paging"`
}

type CreateTicketInput struct {
	Associations []struct {
		Types []struct {
			AssociationCategory AssociationCategory `json:"associationCategory" required:"yes"`
			AssocaitionTypeID   int                 `json:"associaitonTypeId" required:"yes"`
		} `json:"types" required:"yes"`
		To struct {
			ID string `json:"id" required:"yes"`
		} `json:"to" required:"yes"`
	}
	Properties map[string]string `json:"properties" required:"yes"`
}

type CreateTicketResponse struct {
	CreatedResourceID string `json:"createdResourceId" required:"yes"`
	Entity            Ticket `json:"entity" required:"yes"`
	Location          string `json:"location"`
}

type UpdateTicketInput struct {
	Properties map[string]string `json:"properties" required:"yes"`
}

type MergeTwoTicketsInput struct {
	PrimaryObjectID string `json:"primaryObjectId" required:"yes"`
	ObjectIDToMerge string `json:"objectIdToMerge" required:"yes"`
}

type BatchTicketsResponse struct {
	CompletedAt string            `json:"completedAt" required:"yes"`
	StartedAt   string            `json:"startedAt" required:"yes"`
	Results     []Ticket          `json:"results" required:"yes"`
	Status      BatchStatus       `json:"status" required:"yes"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt string            `json:"requestedAt"`
	Links       map[string]string `json:"links"`
	Errors      []BatchError      `json:"errors"`
}

func (batch *BatchTicketsResponse) HasErrors() bool {
	return len(batch.Errors) > 0
}

func (batch *BatchTicketsResponse) GetErrors() []error {
	var errs []error
	for _, err := range batch.Errors {
		errs = append(errs, &err)
	}
	return errs
}

func (batch *BatchTicketsResponse) GetErrorMessages() []string {
	var msgs []string
	for _, err := range batch.Errors {
		msgs = append(msgs, err.Message)
	}
	return msgs
}

type BatchReadTicketsInput struct {
	PropertiesWithHistory []string `json:"propertiesWithHistory" required:"yes"`
	Inputs                []struct {
		ID string `json:"id" required:"yes"`
	} `json:"inputs" required:"yes"`
	Properties []string `json:"properties" required:"yes"`
	IDProperty string   `json:"idProperty"`
}

type BatchCreateTicketsInput struct {
	Inputs []struct {
		Associations []struct {
			Types []struct {
				AssociationCategory AssociationCategory `json:"associationCategory" required:"yes"`
				AssociationTypeID   int                 `json:"associationTypeId" required:"yes"`
			} `json:"types" required:"yes"`
			To struct {
				ID string `json:"id" required:"yes"`
			} `json:"to" required:"yes"`
		} `json:"associations" required:"yes"`
		Properties         map[string]string `json:"properties" required:"yes"`
		ObjectWriteTraceID string            `json:"objectWriteTraceId"`
	} `json:"inputs" required:"yes"`
}

type BatchUpdateTicketsInput struct {
	Inputs []struct {
		ID                 string            `json:"id" required:"yes"`
		Properties         map[string]string `json:"properties" required:"yes"`
		IDProperty         string            `json:"idPoperty"`
		ObjectWriteTraceID string            `json:"objectWriteTraceId"`
	} `json:"inputs" required:"yes"`
}

type BatchCreateOrUpdateTicketsInput struct {
	Inputs []struct {
		ID                 string            `json:"id" required:"yes"`
		Properties         map[string]string `json:"properties" required:"yes"`
		IDProperty         string            `json:"idProperty"`
		ObjectWriteTraceID string            `json:"objectWriteTraceId"`
	} `json:"inputs" required:"yes"`
}

type BatchArchiveTicketsInput struct {
	Inputs []struct {
		ID string `json:"id" required:"yes"`
	} `json:"inputs" required:"yes"`
}

type SearchTicketsInput struct {
	// The maximum results to return, up to 200 objects.
	Limit int `json:"limit" required:"yes"`
	// A paging cursor token for retrieving subsequent pages.
	After string `json:"after" required:"yes"`
	// Specifies sorting order based on object properties.
	Sorts []string `json:"sorts" required:"yes"`
	// A list of property names to include in the response.
	Properties []string `json:"properties" required:"yes"`
	// Up to 6 groups of filters defining additional query criteria.
	FilterGroups []struct {
		Filters []struct {
			// The name of the property to apply the filter to.
			PropertyName string         `json:"propertyName" required:"yes"`
			Operator     FilterOperator `json:"operator" required:"yes"`
			// The upper boundary value when using ranged-based filters.
			HighValue string `json:"highValue"`
			// The values to match against the property.
			Values []string `json:"values"`
			// The value to match against the property.
			Value string `json:"value"`
		} `json:"filters" required:"yes"`
	} `json:"filterGroups" required:"yes"`
	// The search query string, up to 3000 characters.
	Query string `json:"query"`
}

type SearchTicketsResponse struct {
	Total   int      `json:"total" required:"yes"`
	Results []Ticket `json:"results" required:"yes"`
	Paging  struct {
		Next struct {
			After string `json:"after" required:"yes"`
			Link  string `json:"link"`
		} `json:"next"`
		Prev struct {
			Before string `json:"before" required:"yes"`
			Link   string `json:"link"`
		} `json:"prev"`
	} `json:"paging"`
}
