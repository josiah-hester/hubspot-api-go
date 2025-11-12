package lists

import (
	"errors"
	"testing"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListNotFoundError_Error tests the Error() method
func TestListNotFoundError_Error(t *testing.T) {
	err := &ListNotFoundError{
		ListID: "123",
		Original: &client.HubSpotError{
			Status:  404,
			Message: "Not found",
		},
	}

	expectedMsg := "list 123 not found"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestListValidationError_Error tests the Error() method
func TestListValidationError_Error(t *testing.T) {
	err := &ListValidationError{
		Field:   "name",
		Message: "Name is required",
		Original: &client.HubSpotError{
			Status:   400,
			Category: "VALIDATION_ERROR",
		},
	}

	expectedMsg := "validation error on field name: Name is required"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestListAlreadyExistsError_Error tests the Error() method
func TestListAlreadyExistsError_Error(t *testing.T) {
	err := &ListAlreadyExistsError{
		ListName: "My List",
		Original: &client.HubSpotError{
			Status:  409,
			Message: "Conflict",
		},
	}

	expectedMsg := "list with name My List already exists"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestRecordNotFoundError_Error tests the Error() method
func TestRecordNotFoundError_Error(t *testing.T) {
	err := &RecordNotFoundError{
		RecordID: "record-123",
		ListID:   "list-456",
		Original: &client.HubSpotError{
			Status:  404,
			Message: "Not found",
		},
	}

	expectedMsg := "record record-123 not found in list list-456"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestParseListError_NotFound tests parsing 404 errors
func TestParseListError_NotFound(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  404,
		Message: "List not found",
	}

	result := ParseListError(hubspotErr, "123")

	var notFoundErr *ListNotFoundError
	require.ErrorAs(t, result, &notFoundErr)
	assert.Equal(t, "123", notFoundErr.ListID)
	assert.Equal(t, hubspotErr, notFoundErr.Original)
}

// TestParseListError_ValidationError tests parsing 400 validation errors
func TestParseListError_ValidationError(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:   400,
		Category: "VALIDATION_ERROR",
		Message:  "Invalid list name",
	}

	result := ParseListError(hubspotErr, "123")

	var validationErr *ListValidationError
	require.ErrorAs(t, result, &validationErr)
	assert.Equal(t, "Invalid list name", validationErr.Field)
	assert.Equal(t, hubspotErr, validationErr.Original)
}

// TestParseListError_BadRequestNonValidation tests parsing 400 non-validation errors
func TestParseListError_BadRequestNonValidation(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:   400,
		Category: "OTHER_ERROR",
		Message:  "Bad request",
	}

	result := ParseListError(hubspotErr, "123")

	// Should return the original error unchanged since it's not a validation error
	assert.Equal(t, hubspotErr, result)
}

// TestParseListError_Conflict tests parsing 409 conflict errors
func TestParseListError_Conflict(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  409,
		Message: "List already exists",
	}

	result := ParseListError(hubspotErr, "My List")

	var existsErr *ListAlreadyExistsError
	require.ErrorAs(t, result, &existsErr)
	assert.Equal(t, hubspotErr, existsErr.Original)
}

// TestParseListError_OtherHubSpotError tests parsing other HubSpot errors
func TestParseListError_OtherHubSpotError(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  500,
		Message: "Internal server error",
	}

	result := ParseListError(hubspotErr, "123")

	// Should return the original error unchanged for non-mapped status codes
	assert.Equal(t, hubspotErr, result)
}

// TestParseListError_NonHubSpotError tests parsing non-HubSpot errors
func TestParseListError_NonHubSpotError(t *testing.T) {
	regularErr := errors.New("some other error")

	result := ParseListError(regularErr, "123")

	// Should return the original error unchanged
	assert.Equal(t, regularErr, result)
}

// TestParseRecordError_NotFound tests parsing 404 record errors
func TestParseRecordError_NotFound(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:  404,
		Message: "Record not found",
	}

	result := ParseRecordError(hubspotErr, "record-123", "list-456")

	var notFoundErr *RecordNotFoundError
	require.ErrorAs(t, result, &notFoundErr)
	assert.Equal(t, "record-123", notFoundErr.RecordID)
	assert.Equal(t, "list-456", notFoundErr.ListID)
	assert.Equal(t, hubspotErr, notFoundErr.Original)
}

// TestParseRecordError_OtherError tests that non-404 errors fall through to ParseListError
func TestParseRecordError_OtherError(t *testing.T) {
	hubspotErr := &client.HubSpotError{
		Status:   400,
		Category: "VALIDATION_ERROR",
		Message:  "Invalid record ID",
	}

	result := ParseRecordError(hubspotErr, "record-123", "list-456")

	// Should fall through to ParseListError
	var validationErr *ListValidationError
	require.ErrorAs(t, result, &validationErr)
	assert.Equal(t, hubspotErr, validationErr.Original)
}

// TestParseRecordError_NonHubSpotError tests non-HubSpot errors
func TestParseRecordError_NonHubSpotError(t *testing.T) {
	regularErr := errors.New("network error")

	result := ParseRecordError(regularErr, "record-123", "list-456")

	// Should fall through to ParseListError which returns the original error
	assert.Equal(t, regularErr, result)
}
