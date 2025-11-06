package objects

import "strings"

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

type Object struct {
	CreatedAt             string                           `json:"createdAt" required:"yes"`
	ID                    string                           `json:"id" required:"yes"`
	Properties            map[string]string                `json:"properties" required:"yes"`
	UpdatedAt             string                           `json:"updatedAt" required:"yes"`
	Archived              bool                             `json:"archived" required:"yes"`
	Associations          map[string]AssociationResponse   `json:"associations"`
	ArchivedAt            string                           `json:"archivedAt"`
	PropertiesWithHistory map[string][]PropertyWithHistory `json:"propertiesWithHistory"`
	ObjectWriteTraceID    string                           `json:"objectWriteTraceId"`
}

type Paging struct {
	Next struct {
		After string `json:"after"`
		Link  string `json:"link"`
	} `json:"next"`
	Prev struct {
		Before string `json:"before"`
		Link   string `json:"link"`
	}
}

type Association struct {
	Types []struct {
		AssociationCategory AssociationCategory `json:"associationCategory"`
		AssociationTypeID   int                 `json:"associationTypeId"`
	} `json:"types"`
	To struct {
		ID string `json:"id"`
	} `json:"id"`
}

type AssociationResponse struct {
	Results struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"results"`
	Paging Paging `json:"paging"`
}

type PropertyWithHistory struct {
	SourceType      string `json:"sourceType"`
	Value           string `json:"value"`
	Timestamp       string `json:"timestamp"`
	SourceID        string `json:"sourceId"`
	SourceLabel     string `json:"sourceLabel"`
	UpdatedByUserID int    `json:"updatedByUserId"`
}

type ListObjectsResponse struct {
	Results []Object `json:"results"`
	Paging  Paging   `json:"paging"`
}

type CreateObjectInput struct {
	Associations []Association     `json:"associations" required:"yes"`
	Properties   map[string]string `json:"properties" required:"yes"`
}

type CreateObjectResponse struct {
	CreateResourceID string `json:"createResourceId" required:"yes"`
	Entity           Object `json:"entity" required:"yes"`
}

type UpdateObjectInput struct {
	Properties map[string]string `json:"properties" required:"yes"`
}

type MergeObjectsInput struct {
	ObjectIDToMerge string `json:"objectIdToMerge" required:"yes"`
	PrimaryObjectID string `json:"primaryObjectID" required:"yes"`
}

type ObjectError struct {
	Message     string              `json:"message" required:"yes"`
	SubCategory string              `json:"subCategory"`
	Code        string              `json:"code"`
	In          string              `json:"in"`
	Context     map[string][]string `json:"context"`
}

type BatchError struct {
	Context     map[string][]string `json:"context" required:"yes"`
	Links       map[string]string   `json:"links" required:"yes"`
	Category    string              `json:"category" required:"yes"`
	Message     string              `json:"message" required:"yes"`
	Errors      []ObjectError       `json:"errors" required:"yes"`
	Status      string              `json:"status" required:"yes"`
	SubCategory any                 `json:"subCategory"`
	ID          string              `json:"id"`
}

func (e *BatchError) Error() string {
	if e.Message != "" && len(e.Errors) == 0 {
		return e.Message
	}

	var sb strings.Builder
	if e.Message != "" {
		sb.WriteString(e.Message)
	} else {
		sb.WriteString("batch error")
	}

	if len(e.Errors) > 0 {
		sb.WriteString(": ")
		for i, err := range e.Errors {
			if i > 0 {
				sb.WriteString("; ")
			}
			sb.WriteString(err.Message)
		}
	}

	return sb.String()
}

type BatchReadObjectsInput struct {
	PropertiesWithHistory []string `json:"propertiesWithHistory" required:"yes"`
	Inputs                []struct {
		ID string `json:"id" required:"yes"`
	} `json:"inputs" required:"yes"`
	Properties []string `json:"properties" required:"yes"`
	IDProperty string   `json:"idProperty"`
}

type BatchCreateObjectsInput struct {
	Inputs []struct {
		Associations       []Association     `json:"associations" required:"yes"`
		Properties         map[string]string `json:"properties" required:"yes"`
		ObjectWriteTraceID string            `json:"objectWriteTraceId"`
	} `json:"inputs" required:"yes"`
}

type BatchUpdateObjectsInput struct {
	Inputs []struct {
		ID                 string            `json:"id" required:"yes"`
		Properties         map[string]string `json:"properties" required:"yes"`
		IDProperty         string            `json:"idProperty"`
		ObjectWriteTraceID string            `json:"objectWriteTraceId"`
	} `json:"inputs" required:"yes"`
}

type BatchCreateOrUpdateObjectsInput struct {
	Inputs []struct {
		ID                 string            `json:"id" required:"yes"`
		Properties         map[string]string `json:"properties" required:"yes"`
		IDProperty         string            `json:"idProperty"`
		ObjectWriteTraceID string            `json:"objectWriteTraceId"`
	} `json:"inputs" required:"yes"`
}

type BatchArchiveObjectsInput struct {
	Inputs []struct {
		ID string `json:"id" required:"yes"`
	} `json:"inputs" required:"yes"`
}

type BatchResponse struct {
	CompletedAt string            `json:"completedAt" required:"yes"`
	StartedAt   string            `json:"startedAt" required:"yes"`
	Results     []Object          `json:"results" required:"yes"`
	Status      BatchStatus       `json:"status" required:"yes"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt string            `json:"requestedAt"`
	Links       map[string]string `json:"links"`
	Errors      []BatchError      `json:"errors"`
}

type SearchObjectsInput struct {
	Limit        int      `json:"limit" required:"yes"`
	After        string   `json:"after" required:"yes"`
	Sorts        []string `json:"sorts" required:"yes"`
	Properties   []string `json:"properties" required:"yes"`
	FilterGroups []struct {
		Filters []struct {
			PropertyName string         `json:"propertyName" required:"yes"`
			Operator     FilterOperator `json:"operator" required:"yes"`
			HighValue    string         `json:"highValue"`
			Values       []string       `json:"values"`
			Value        string         `json:"value"`
		} `json:"filters" required:"yes"`
	} `json:"filterGroups" required:"yes"`
	Query string `json:"query"`
}

type SearchObjectsResponse struct {
	Total   int      `json:"total" required:"yes"`
	Results []Object `json:"results" required:"yes"`
	Paging  Paging   `json:"paging"`
}
