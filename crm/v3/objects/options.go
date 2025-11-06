package objects

import (
	"strings"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

// ObjectsOption is a functional option for Object calls Query Parameters
type ObjectsOption func(*client.Request)

// WithLimit sets the maximum number of objects to return
func WithLimit(limit int) ObjectsOption {
	return func(req *client.Request) {
		if limit > 0 {
			req.AddQueryParam("limit", string(rune(limit)))
		}
	}
}

// WithAfter sets the pagination cursor
func WithAfter(after string) ObjectsOption {
	return func(req *client.Request) {
		if after != "" {
			req.AddQueryParam("after", after)
		}
	}
}

// WithProperties specifies which properties to retrieve
func WithProperties(props []string) ObjectsOption {
	return func(req *client.Request) {
		if len(props) > 0 {
			req.AddQueryParam("properties", strings.Join(props, ","))
		}
	}
}

// WithPropertiesWithHistory specifies which properties' history to retrieve
func WithPropertiesWithHistory(props []string) ObjectsOption {
	return func(req *client.Request) {
		if len(props) > 0 {
			req.AddQueryParam("propertiesWithHistory", strings.Join(props, ","))
		}
	}
}

// WithAssociations specifies which associations to retrieve
func WithAssociations(associations []string) ObjectsOption {
	return func(req *client.Request) {
		if len(associations) > 0 {
			req.AddQueryParam("associations", strings.Join(associations, ","))
		}
	}
}

// WithArchived specifies whether to include archived objects in the response
func WithArchived() ObjectsOption {
	return func(req *client.Request) {
		req.AddQueryParam("archived", "true")
	}
}

// WithIDProperty specifies the property to use as the object identifier
func WithIDProperty(property string) ObjectsOption {
	return func(req *client.Request) {
		req.AddQueryParam("idProperty", property)
	}
}
