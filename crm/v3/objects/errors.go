package objects

import (
	"fmt"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

// ObjectNotFoundError is returned when an object is not found
type ObjectNotFoundError struct {
	ObjectType string
	Original   *client.HubSpotError
}

func (e *ObjectNotFoundError) Error() string {
	return fmt.Sprintf("object %s not found", e.ObjectType)
}

type ObjectValidationError struct {
	Field    string
	Message  string
	Original *client.HubSpotError
}

func (e *ObjectValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

type ObjectAlreadyExistsError struct {
	ObjectID string
	Original *client.HubSpotError
}

func (e *ObjectAlreadyExistsError) Error() string {
	return fmt.Sprintf("object with id %s already exists", e.ObjectID)
}

func ParseObjectError(err error, objectType string) error {
	if hubspotErr, ok := err.(*client.HubSpotError); ok {
		switch hubspotErr.Status {
		case 404:
			return &ObjectNotFoundError{
				ObjectType: objectType,
				Original:   hubspotErr,
			}
		case 400:
			if hubspotErr.Category == "VALIDATION_ERROR" {
				return &ObjectValidationError{
					Field:    hubspotErr.Message,
					Original: hubspotErr,
				}
			}
		case 409:
			return &ObjectAlreadyExistsError{
				Original: hubspotErr,
			}
		}
	}
	return err
}
