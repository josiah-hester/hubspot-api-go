package orders

import (
	"fmt"
	"strings"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// OrderOption represents a functional option for order requests
type OrderOption func(*client.Request)

// WithProperties specifies which properties to return
func WithProperties(properties []string) OrderOption {
	return func(req *client.Request) {
		req.AddQueryParam("properties", strings.Join(properties, ","))
	}
}

// WithPropertiesWithHistory specifies which properties to return with history
func WithPropertiesWithHistory(properties []string) OrderOption {
	return func(req *client.Request) {
		req.AddQueryParam("propertiesWithHistory", strings.Join(properties, ","))
	}
}

// WithAssociations specifies which associations to return
func WithAssociations(associations []string) OrderOption {
	return func(req *client.Request) {
		req.AddQueryParam("associations", strings.Join(associations, ","))
	}
}

// WithLimit sets the maximum number of results per page
func WithLimit(limit int) OrderOption {
	return func(req *client.Request) {
		req.AddQueryParam("limit", fmt.Sprintf("%d", limit))
	}
}

// WithAfter sets the paging cursor
func WithAfter(after string) OrderOption {
	return func(req *client.Request) {
		req.AddQueryParam("after", after)
	}
}

// WithArchived includes archived orders
func WithArchived() OrderOption {
	return func(req *client.Request) {
		req.AddQueryParam("archived", "true")
	}
}

// WithIDProperty specifies a unique identifier property to use instead of ID
func WithIDProperty(property string) OrderOption {
	return func(req *client.Request) {
		req.AddQueryParam("idProperty", property)
	}
}
