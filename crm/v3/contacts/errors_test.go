package contacts

import (
	"errors"
	"testing"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContactNotFoundError_Error tests the Error() method
func TestContactNotFoundError_Error(t *testing.T) {
	err := &ContactNotFoundError{
		ContactID: "12345",
		Original: &client.HubSpotError{
			Status:  404,
			Message: "Not found",
		},
	}

	expectedMsg := "contact 12345 not found"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestContactValidationError_Error tests the Error() method
func TestContactValidationError_Error(t *testing.T) {
	err := &ContactValidationError{
		Field:   "email",
		Message: "Invalid email format",
		Original: &client.HubSpotError{
			Status:   400,
			Category: "VALIDATION_ERROR",
		},
	}

	expectedMsg := "validation error on field email: Invalid email format"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestContactAlreadyExistsError_Error tests the Error() method
func TestContactAlreadyExistsError_Error(t *testing.T) {
	err := &ContactAlreadyExistsError{
		ContactID: "test@example.com",
		Original: &client.HubSpotError{
			Status:  409,
			Message: "Conflict",
		},
	}

	expectedMsg := "contact with email test@example.com already exists"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestParseContactError_NotFound tests parsing 404 errors
func TestParseContactError_NotFound(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  404,
		Message: "Contact not found",
	}

	result := ParseContactError(hubspotErr, "12345")

	var notFoundErr *ContactNotFoundError
	require.ErrorAs(t, result, &notFoundErr)
	assert.Equal(t, "12345", notFoundErr.ContactID)
	assert.Equal(t, hubspotErr, notFoundErr.Original)
}

// TestParseContactError_ValidationError tests parsing 400 validation errors
func TestParseContactError_ValidationError(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:   400,
		Category: "VALIDATION_ERROR",
		Message:  "Invalid email format",
	}

	result := ParseContactError(hubspotErr, "12345")

	var validationErr *ContactValidationError
	require.ErrorAs(t, result, &validationErr)
	assert.Equal(t, "Invalid email format", validationErr.Field)
	assert.Equal(t, hubspotErr, validationErr.Original)
}

// TestParseContactError_BadRequestNonValidation tests parsing 400 non-validation errors
func TestParseContactError_BadRequestNonValidation(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:   400,
		Category: "OTHER_ERROR",
		Message:  "Bad request",
	}

	result := ParseContactError(hubspotErr, "12345")

	// Should return the original error unchanged since it's not a validation error
	assert.Equal(t, hubspotErr, result)
}

// TestParseContactError_Conflict tests parsing 409 conflict errors
func TestParseContactError_Conflict(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  409,
		Message: "Contact already exists",
	}

	result := ParseContactError(hubspotErr, "test@example.com")

	var existsErr *ContactAlreadyExistsError
	require.ErrorAs(t, result, &existsErr)
	assert.Equal(t, hubspotErr, existsErr.Original)
}

// TestParseContactError_OtherHubSpotError tests parsing other HubSpot errors
func TestParseContactError_OtherHubSpotError(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  500,
		Message: "Internal server error",
	}

	result := ParseContactError(hubspotErr, "12345")

	// Should return the original error unchanged for non-mapped status codes
	assert.Equal(t, hubspotErr, result)
}

// TestParseContactError_NonHubSpotError tests parsing non-HubSpot errors
func TestParseContactError_NonHubSpotError(t *testing.T) {
	regularErr := errors.New("some other error")

	result := ParseContactError(regularErr, "12345")

	// Should return the original error unchanged
	assert.Equal(t, regularErr, result)
}
