package contacts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockServer creates a test server with custom handler
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)

	apiClient, err := client.NewClient(
		client.WithTimeout(5*time.Second),
		client.WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	return server, NewClient(apiClient)
}

// respondJSON writes a JSON string response
func respondJSON(w http.ResponseWriter, statusCode int, jsonString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(jsonString))
}

// TestGetContact_Success tests successful contact retrieval
func TestGetContact_Success(t *testing.T) {
	contactJSON := `{
		"id": "12345",
		"properties": {
			"email": "test@example.com",
			"firstname": "John",
			"lastname": "Doe"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/12345", r.URL.Path)
		respondJSON(w, http.StatusOK, contactJSON)
	})
	defer server.Close()

	contact, err := contactClient.GetContact(context.Background(), "12345")

	require.NoError(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, "12345", contact.ID)
	assert.Equal(t, "test@example.com", contact.Properties["email"])
	assert.Equal(t, "John", contact.Properties["firstname"])
	assert.Equal(t, "Doe", contact.Properties["lastname"])
	assert.False(t, contact.Archived)
}

// TestGetContact_WithProperties tests GetContact with properties option
func TestGetContact_WithProperties(t *testing.T) {
	contactJSON := `{
		"id": "12345",
		"properties": {
			"email": "test@example.com",
			"firstname": "John",
			"lastname": "Doe"
		},
		"archived": false
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/12345", r.URL.Path)
		assert.Equal(t, "email,firstname,lastname", r.URL.Query().Get("properties"))
		respondJSON(w, http.StatusOK, contactJSON)
	})
	defer server.Close()

	contact, err := contactClient.GetContact(
		context.Background(),
		"12345",
		WithProperties([]string{"email", "firstname", "lastname"}),
	)

	require.NoError(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, "12345", contact.ID)
}

// TestGetContact_WithAssociations tests GetContact with associations option
func TestGetContact_WithAssociations(t *testing.T) {
	contactJSON := `{
		"id": "12345",
		"properties": {},
		"archived": false
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "companies,deals", r.URL.Query().Get("associations"))
		respondJSON(w, http.StatusOK, contactJSON)
	})
	defer server.Close()

	contact, err := contactClient.GetContact(
		context.Background(),
		"12345",
		WithAssociations([]string{"companies", "deals"}),
	)

	require.NoError(t, err)
	assert.NotNil(t, contact)
}

// TestGetContact_WithAllOptions tests GetContact with multiple options
func TestGetContact_WithAllOptions(t *testing.T) {
	contactJSON := `{
		"id": "12345",
		"properties": {
			"email": "test@example.com"
		},
		"archived": true
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "email,firstname", r.URL.Query().Get("properties"))
		assert.Equal(t, "companies", r.URL.Query().Get("associations"))
		assert.Equal(t, "email", r.URL.Query().Get("idProperty"))
		assert.Equal(t, "true", r.URL.Query().Get("archived"))
		respondJSON(w, http.StatusOK, contactJSON)
	})
	defer server.Close()

	contact, err := contactClient.GetContact(
		context.Background(),
		"12345",
		WithProperties([]string{"email", "firstname"}),
		WithAssociations([]string{"companies"}),
		WithIDProperty("email"),
		WithArchived(),
	)

	require.NoError(t, err)
	assert.NotNil(t, contact)
	assert.True(t, contact.Archived)
}

// TestGetContact_NotFound tests 404 error handling
func TestGetContact_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Contact not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	contact, err := contactClient.GetContact(context.Background(), "99999")

	require.Error(t, err)
	assert.Nil(t, contact)

	var notFoundErr *ContactNotFoundError
	require.ErrorAs(t, err, &notFoundErr)
	assert.Equal(t, "99999", notFoundErr.ContactID)
}

// TestGetContact_InvalidJSON tests JSON unmarshal error
func TestGetContact_InvalidJSON(t *testing.T) {
	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	contact, err := contactClient.GetContact(context.Background(), "12345")

	require.Error(t, err)
	assert.Nil(t, contact)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestCreateContact_Success tests successful contact creation
func TestCreateContact_Success(t *testing.T) {
	responseJSON := `{
		"id": "12345",
		"properties": {
			"email": "test@example.com",
			"firstname": "John",
			"lastname": "Doe"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts", r.URL.Path)

		var input CreateContactInput
		err := json.NewDecoder(r.Body).Decode(&input)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", input.Properties["email"])
		assert.Equal(t, "John", input.Properties["firstname"])

		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateContactInput{
		Properties: map[string]string{
			"email":     "test@example.com",
			"firstname": "John",
			"lastname":  "Doe",
		},
	}

	contact, err := contactClient.CreateContact(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, "12345", contact.ID)
	assert.Equal(t, "test@example.com", contact.Properties["email"])
}

// TestCreateContact_ValidationError tests 400 validation error
func TestCreateContact_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid email",
		"category": "VALIDATION_ERROR"
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &CreateContactInput{
		Properties: map[string]string{
			"email": "invalid-email",
		},
	}

	contact, err := contactClient.CreateContact(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, contact)

	var validationErr *ContactValidationError
	require.ErrorAs(t, err, &validationErr)
}

// TestCreateContact_AlreadyExists tests 409 conflict error
func TestCreateContact_AlreadyExists(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Contact already exists",
		"category": "CONFLICT"
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusConflict, errorJSON)
	})
	defer server.Close()

	input := &CreateContactInput{
		Properties: map[string]string{
			"email": "existing@example.com",
		},
	}

	contact, err := contactClient.CreateContact(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, contact)

	var existsErr *ContactAlreadyExistsError
	require.ErrorAs(t, err, &existsErr)
}

// TestUpdateContact_Success tests successful contact update
func TestUpdateContact_Success(t *testing.T) {
	responseJSON := `{
		"id": "12345",
		"properties": {
			"email": "test@example.com",
			"firstname": "Jane",
			"lastname": "Doe"
		},
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"archived": false
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/12345", r.URL.Path)

		var input UpdateContactInput
		err := json.NewDecoder(r.Body).Decode(&input)
		assert.NoError(t, err)
		assert.Equal(t, "Jane", input.Properties["firstname"])

		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateContactInput{
		Properties: map[string]string{
			"firstname": "Jane",
		},
	}

	contact, err := contactClient.UpdateContact(context.Background(), "12345", input)

	require.NoError(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, "12345", contact.ID)
	assert.Equal(t, "Jane", contact.Properties["firstname"])
}

// TestUpdateContact_NotFound tests 404 error on update
func TestUpdateContact_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Contact not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	input := &UpdateContactInput{
		Properties: map[string]string{
			"firstname": "Jane",
		},
	}

	contact, err := contactClient.UpdateContact(context.Background(), "99999", input)

	require.Error(t, err)
	assert.Nil(t, contact)

	var notFoundErr *ContactNotFoundError
	require.ErrorAs(t, err, &notFoundErr)
	assert.Equal(t, "99999", notFoundErr.ContactID)
}

// TestUpdateContact_ValidationError tests 400 validation error on update
func TestUpdateContact_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid property",
		"category": "VALIDATION_ERROR"
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &UpdateContactInput{
		Properties: map[string]string{
			"invalid_field": "value",
		},
	}

	contact, err := contactClient.UpdateContact(context.Background(), "12345", input)

	require.Error(t, err)
	assert.Nil(t, contact)

	var validationErr *ContactValidationError
	require.ErrorAs(t, err, &validationErr)
}

// TestDeleteContact_Success tests successful contact deletion
func TestDeleteContact_Success(t *testing.T) {
	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts/12345", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := contactClient.DeleteContact(context.Background(), "12345")

	assert.NoError(t, err)
}

// TestDeleteContact_NotFound tests 404 error on delete
func TestDeleteContact_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Contact not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	err := contactClient.DeleteContact(context.Background(), "99999")

	require.Error(t, err)

	var notFoundErr *ContactNotFoundError
	require.ErrorAs(t, err, &notFoundErr)
	assert.Equal(t, "99999", notFoundErr.ContactID)
}

// TestListContacts_Success tests successful contact listing
func TestListContacts_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test1@example.com",
					"firstname": "John",
					"lastname": "Doe"
				},
				"archived": false
			},
			{
				"id": "2",
				"properties": {
					"email": "test2@example.com",
					"firstname": "Jane",
					"lastname": "Smith"
				},
				"archived": false
			}
		],
		"paging": {
			"next": {
				"after": "next-cursor",
				"link": "https://api.hubapi.com/next"
			}
		}
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/contacts", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	contacts, nextCursor, err := contactClient.ListContacts(context.Background())

	require.NoError(t, err)
	assert.Len(t, contacts, 2)
	assert.Equal(t, "1", contacts[0].ID)
	assert.Equal(t, "2", contacts[1].ID)
	assert.Equal(t, "next-cursor", nextCursor)
}

// TestListContacts_WithLimit tests list with limit option
func TestListContacts_WithLimit(t *testing.T) {
	responseJSON := `{
		"results": [],
		"paging": {}
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	_, _, err := contactClient.ListContacts(context.Background(), WithLimit(10))

	require.NoError(t, err)
}

// TestListContacts_WithAfter tests list with pagination cursor
func TestListContacts_WithAfter(t *testing.T) {
	responseJSON := `{
		"results": [],
		"paging": {}
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "cursor123", r.URL.Query().Get("after"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	_, _, err := contactClient.ListContacts(context.Background(), WithAfter("cursor123"))

	require.NoError(t, err)
}

// TestListContacts_EmptyResults tests list with no results
func TestListContacts_EmptyResults(t *testing.T) {
	responseJSON := `{
		"results": [],
		"paging": {}
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	contacts, nextCursor, err := contactClient.ListContacts(context.Background())

	require.NoError(t, err)
	assert.Empty(t, contacts)
	assert.Empty(t, nextCursor)
}

// TestListContacts_NoPaging tests list without pagination
func TestListContacts_NoPaging(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"id": "1",
				"properties": {
					"email": "test@example.com"
				},
				"archived": false
			}
		]
	}`

	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	contacts, nextCursor, err := contactClient.ListContacts(context.Background())

	require.NoError(t, err)
	assert.Len(t, contacts, 1)
	assert.Empty(t, nextCursor)
}

// TestListContacts_InvalidJSON tests JSON unmarshal error
func TestListContacts_InvalidJSON(t *testing.T) {
	server, contactClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	contacts, nextCursor, err := contactClient.ListContacts(context.Background())

	require.Error(t, err)
	assert.Nil(t, contacts)
	assert.Empty(t, nextCursor)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}
