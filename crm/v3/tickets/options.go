package tickets

import (
	"strings"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

type TicketOption func(*client.Request)

// WithLimit sets the maximum number of tickets to return
func WithLimit(limit int) TicketOption {
	return func(req *client.Request) {
		req.AddQueryParam("limit", string(rune(limit)))
	}
}

// WithAfter sets the pagination cursor
func WithAfter(after string) TicketOption {
	return func(req *client.Request) {
		req.AddQueryParam("after", after)
	}
}

// WithProperties specifies which properties to retrieve
func WithProperties(properties []string) TicketOption {
	return func(req *client.Request) {
		req.AddQueryParam("properties", strings.Join(properties, ","))
	}
}

// WithPropertiesWithHistory specifies which properties' history to retrieve
func WithPropertiesWithHistory(properties []string) TicketOption {
	return func(req *client.Request) {
		req.AddQueryParam("propertiesWithHistory", strings.Join(properties, ","))
	}
}

// WithAssociations specifies which associations to retrieve
func WithAssociations(associations []string) TicketOption {
	return func(req *client.Request) {
		req.AddQueryParam("associations", strings.Join(associations, ","))
	}
}

// WithArchived specifies whether to include archived tickets in the response
func WithArchived() TicketOption {
	return func(req *client.Request) {
		req.AddQueryParam("archived", "true")
	}
}

// WithIDProperty specifies the property to use as the ticket identifier
func WithIDProperty(property string) TicketOption {
	return func(req *client.Request) {
		req.AddQueryParam("idProperty", property)
	}
}
