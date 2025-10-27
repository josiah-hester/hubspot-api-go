package contacts

import (
	"strings"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

// GetContactOption is a functional option for GetContact
type GetContactOption func(*client.Request)

// WithProperties specifies which properties to retrieve
func WithProperties(properties []string) GetContactOption {
	return func(req *client.Request) {
		req.AddQueryParam("properties", strings.Join(properties, ","))
	}
}

// WithAssociations specifies which associations to retrieve
func WithAssociations(associations []string) GetContactOption {
	return func(req *client.Request) {
		req.AddQueryParam("associations", strings.Join(associations, ","))
	}
}

// WithIdProperty specifies the property to use as the contact identifier
func WithIdProperty(property string) GetContactOption {
	return func(req *client.Request) {
		req.AddQueryParam("idProperty", property)
	}
}

// WithArchived includes archived contacts in results
func WithArchived(archived bool) GetContactOption {
	return func(req *client.Request) {
		if archived {
			req.AddQueryParam("archived", "true")
		}
	}
}

// ListContactsOption is a functional option for ListContacts
type ListContactsOption func(*client.Request)

// WithLimit sets the maximum number of contacts to return
func WithLimit(limit int) ListContactsOption {
	return func(req *client.Request) {
		req.AddQueryParam("limit", string(rune(limit)))
	}
}

// WithAfter sets the pagination cursor
func WithAfter(after string) ListContactsOption {
	return func(req *client.Request) {
		req.AddQueryParam("after", after)
	}
}
