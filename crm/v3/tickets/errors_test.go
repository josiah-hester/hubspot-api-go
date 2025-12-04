package tickets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBatchError_Error_WithMessage tests BatchError.Error() with message only
func TestBatchError_Error_WithMessage(t *testing.T) {
	err := &BatchError{
		Message:  "Batch operation failed",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors:   []ObjectError{},
	}

	expectedMsg := "Batch operation failed"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestBatchError_Error_WithErrors tests BatchError.Error() with errors
func TestBatchError_Error_WithErrors(t *testing.T) {
	err := &BatchError{
		Message:  "Batch operation failed",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors: []ObjectError{
			{
				Message: "Ticket 1 not found",
			},
			{
				Message: "Ticket 2 validation error",
			},
		},
	}

	expectedMsg := "Batch operation failed: Ticket 1 not found; Ticket 2 validation error"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestBatchError_Error_WithoutMessage tests BatchError.Error() without message
func TestBatchError_Error_WithoutMessage(t *testing.T) {
	err := &BatchError{
		Message:  "",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors: []ObjectError{
			{
				Message: "Ticket 1 not found",
			},
		},
	}

	expectedMsg := "batch error: Ticket 1 not found"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestBatchError_Error_EmptyMessageNoErrors tests BatchError.Error() with no message and no errors
func TestBatchError_Error_EmptyMessageNoErrors(t *testing.T) {
	err := &BatchError{
		Message:  "",
		Category: "BATCH_ERROR",
		Status:   "error",
		Context:  make(map[string][]string),
		Links:    make(map[string]string),
		Errors:   []ObjectError{},
	}

	expectedMsg := "batch error"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestBatchTicketsResponse_HasErrors tests HasErrors method
func TestBatchTicketsResponse_HasErrors(t *testing.T) {
	// Test with errors
	respWithErrors := &BatchTicketsResponse{
		Errors: []BatchError{
			{Message: "Error 1"},
		},
	}
	assert.True(t, respWithErrors.HasErrors())

	// Test without errors
	respWithoutErrors := &BatchTicketsResponse{
		Errors: []BatchError{},
	}
	assert.False(t, respWithoutErrors.HasErrors())
}

// TestBatchTicketsResponse_GetErrors tests GetErrors method
func TestBatchTicketsResponse_GetErrors(t *testing.T) {
	resp := &BatchTicketsResponse{
		Errors: []BatchError{
			{Message: "Error 1"},
			{Message: "Error 2"},
		},
	}

	errors := resp.GetErrors()
	assert.Len(t, errors, 2)
}

// TestBatchTicketsResponse_GetErrorMessages tests GetErrorMessages method
func TestBatchTicketsResponse_GetErrorMessages(t *testing.T) {
	resp := &BatchTicketsResponse{
		Errors: []BatchError{
			{Message: "Ticket 1 not found"},
			{Message: "Ticket 2 validation error"},
			{Message: "Ticket 3 permission denied"},
		},
	}

	expectedMessages := []string{
		"Ticket 1 not found",
		"Ticket 2 validation error",
		"Ticket 3 permission denied",
	}

	assert.Equal(t, expectedMessages, resp.GetErrorMessages())
}

// TestBatchTicketsResponse_GetErrorMessages_Empty tests GetErrorMessages with no errors
func TestBatchTicketsResponse_GetErrorMessages_Empty(t *testing.T) {
	resp := &BatchTicketsResponse{
		Errors: []BatchError{},
	}

	assert.Empty(t, resp.GetErrorMessages())
}
