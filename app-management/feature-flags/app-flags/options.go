package appflags

import (
	"fmt"

	"github.com/josiah-hester/go-hubspot-sdk/client"
)

// AppFlagOption is a functional option for RetrieveAppFeatureFlags
type AppFlagOption func(*client.Request)

func WithLimit(limit int) AppFlagOption {
	return func(req *client.Request) {
		req.AddQueryParam("limit", fmt.Sprintf("%d", limit))
	}
}

func WithStartPortalID(startPortalID int) AppFlagOption {
	return func(req *client.Request) {
		req.AddQueryParam("startPortalId", fmt.Sprintf("%d", startPortalID))
	}
}
