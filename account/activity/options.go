package activity

import (
	"fmt"
	"strings"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// AccountActivityOption is a functional option for RetrieveAuditLogs
type AccountActivityOption func(*client.Request)

// WithActingUserID The ID of a user, for retrieving user-specific logs.
func WithActingUserID(actingUserID int) AccountActivityOption {
	return func(req *client.Request) {
		req.AddQueryParam("actingUserId", fmt.Sprintf("%d", actingUserID))
	}
}

// WithUserID The ID of a user, for retrieving user-specific logs.
func WithUserID(userID int) AccountActivityOption {
	return func(req *client.Request) {
		req.AddQueryParam("userId", fmt.Sprintf("%d", userID))
	}
}

// WithAfter The paging cursor token of the last successfully read resource will be returned as the paging.next.after JSON property of a paged response containing more results.
func WithAfter(after string) AccountActivityOption {
	return func(req *client.Request) {
		req.AddQueryParam("after", after)
	}
}

// WithLimit The maximum number of results to display per page.
func WithLimit(limit int) AccountActivityOption {
	return func(req *client.Request) {
		req.AddQueryParam("limit", fmt.Sprintf("%d", limit))
	}
}

// WithOccurredBeforeAndAfter Retrieve audit logs that occurred between the specified times
func WithOccurredBeforeAndAfter(after, before string) AccountActivityOption {
	return func(req *client.Request) {
		if after != "" {
			req.AddQueryParam("occurredAfter", after)
		}
		if before != "" {
			req.AddQueryParam("occurredBefore", before)
		}
	}
}

// WithSort The fields by which results are sorted
func WithSort(sort []string) AccountActivityOption {
	return func(req *client.Request) {
		req.AddQueryParam("sort", strings.Join(sort, ","))
	}
}
