package contacts

import (
	"fmt"

	"github.com/aacc-dev/go-hubspot-sdk/client"
)

// ContactNotFoundError is returned when a contact is not found
type ContactNotFoundError struct {
	ContactID string
	Original  *client.HubSpotError
}

func (e *ContactNotFoundError) Error() string {
	return fmt.Sprintf("contact %s not found", e.ContactID)
}

// ContactValidationError is returned on validation failures
type ContactValidationError struct {
	Field    string
	Message  string
	Original *client.HubSpotError
}

func (e *ContactValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// ContactAlreadyExistsError is returned when trying to create a duplicate
type ContactAlreadyExistsError struct {
	ContactID string
	Original  *client.HubSpotError
}

func (e *ContactAlreadyExistsError) Error() string {
	return fmt.Sprintf("contact with email %s already exists", e.ContactID)
}

// ParseContactError converts a generic HubSpot error to a contact-specific error
func ParseContactError(err error, contactID string) error {
	if hubspotErr, ok := err.(*client.HubSpotError); ok {
		switch hubspotErr.Status {
		case 404:
			return &ContactNotFoundError{
				ContactID: contactID,
				Original:  hubspotErr,
			}
		case 400:
			if hubspotErr.Category == "VALIDATION_ERROR" {
				return &ContactValidationError{
					Field:    hubspotErr.Message,
					Original: hubspotErr,
				}
			}
		case 409:
			return &ContactAlreadyExistsError{
				Original: hubspotErr,
			}
		}
	}
	return err
}
