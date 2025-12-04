package tickets

import (
	"context"
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

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	apiClient, err := client.NewClient(
		client.WithTimeout(5 * time.Second),
	)
	require.NoError(t, err)

	ticketsClient := NewClient(apiClient)
	assert.NotNil(t, ticketsClient)
	assert.NotNil(t, ticketsClient.apiClient)
}

// TestListTickets_Success tests successful tickets retrieval
func TestListTickets_Success(t *testing.T) {
	ticketsJSON := `{
		"results": [
			{
				"id": "123456",
				"properties": {
					"subject": "Test Ticket",
					"hs_pipeline_stage": "1",
					"hs_ticket_priority": "HIGH"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets", r.URL.Path)
		respondJSON(w, http.StatusOK, ticketsJSON)
	})
	defer server.Close()

	result, err := ticketsClient.ListTickets(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Results, 1)
	assert.Equal(t, "123456", result.Results[0].ID)
}

// TestListTickets_WithOptions tests list with various options
func TestListTickets_WithOptions(t *testing.T) {
	ticketsJSON := `{
		"results": [
			{
				"id": "123456",
				"properties": {
					"subject": "Test Ticket"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "subject,hs_ticket_priority", r.URL.Query().Get("properties"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "abc123", r.URL.Query().Get("after"))
		assert.Equal(t, "subject", r.URL.Query().Get("propertiesWithHistory"))
		assert.Equal(t, "contacts", r.URL.Query().Get("associations"))
		assert.Equal(t, "true", r.URL.Query().Get("archived"))
		respondJSON(w, http.StatusOK, ticketsJSON)
	})
	defer server.Close()

	result, err := ticketsClient.ListTickets(
		context.Background(),
		WithProperties([]string{"subject", "hs_ticket_priority"}),
		WithLimit(10),
		WithAfter("abc123"),
		WithPropertiesWithHistory([]string{"subject"}),
		WithAssociations([]string{"contacts"}),
		WithArchived(),
	)

	require.NoError(t, err)
	assert.Len(t, result.Results, 1)
}

// TestListTickets_NoResults tests when no tickets found
func TestListTickets_NoResults(t *testing.T) {
	ticketsJSON := `{
		"results": []
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, ticketsJSON)
	})
	defer server.Close()

	result, err := ticketsClient.ListTickets(context.Background())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no tickets found")
}

// TestListTickets_InvalidJSON tests invalid JSON response
func TestListTickets_InvalidJSON(t *testing.T) {
	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	result, err := ticketsClient.ListTickets(context.Background())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestCreateTicket_Success tests successful ticket creation
func TestCreateTicket_Success(t *testing.T) {
	responseJSON := `{
		"createdResourceId": "resource-123",
		"entity": {
			"id": "123456",
			"properties": {
				"subject": "New Ticket"
			},
			"createdAt": "2024-01-01T00:00:00.000Z",
			"updatedAt": "2024-01-01T00:00:00.000Z",
			"archived": false
		}
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets", r.URL.Path)
		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateTicketInput{
		Properties: map[string]string{
			"subject": "New Ticket",
		},
	}

	result, err := ticketsClient.CreateTicket(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123456", result.Entity.ID)
}

// TestCreateTicket_ValidationError tests validation error
func TestCreateTicket_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid property",
		"category": "VALIDATION_ERROR"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &CreateTicketInput{
		Properties: map[string]string{
			"invalid_field": "value",
		},
	}

	result, err := ticketsClient.CreateTicket(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
}

// TestCreateTicket_InvalidJSON tests invalid JSON response
func TestCreateTicket_InvalidJSON(t *testing.T) {
	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &CreateTicketInput{
		Properties: map[string]string{
			"subject": "Test",
		},
	}

	result, err := ticketsClient.CreateTicket(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestReadTicket_Success tests successful ticket retrieval
func TestReadTicket_Success(t *testing.T) {
	ticketJSON := `{
		"id": "123456",
		"properties": {
			"subject": "Test Ticket"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/123456", r.URL.Path)
		assert.Equal(t, "email", r.URL.Query().Get("idProperty"))
		respondJSON(w, http.StatusOK, ticketJSON)
	})
	defer server.Close()

	ticket, err := ticketsClient.ReadTicket(context.Background(), "123456", WithIDProperty("email"))

	require.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "123456", ticket.ID)
}

// TestReadTicket_NotFound tests 404 error handling
func TestReadTicket_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Ticket not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	ticket, err := ticketsClient.ReadTicket(context.Background(), "99999")

	require.Error(t, err)
	assert.Nil(t, ticket)
}

// TestReadTicket_InvalidJSON tests invalid JSON response
func TestReadTicket_InvalidJSON(t *testing.T) {
	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	ticket, err := ticketsClient.ReadTicket(context.Background(), "123456")

	require.Error(t, err)
	assert.Nil(t, ticket)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestUpdateTicket_Success tests successful ticket update
func TestUpdateTicket_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"subject": "Updated Ticket"
		},
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"createdAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/123456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateTicketInput{
		Properties: map[string]string{
			"subject": "Updated Ticket",
		},
	}

	ticket, err := ticketsClient.UpdateTicket(context.Background(), "123456", input)

	require.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "123456", ticket.ID)
}

// TestUpdateTicket_NotFound tests 404 error on update
func TestUpdateTicket_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Ticket not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	input := &UpdateTicketInput{
		Properties: map[string]string{
			"subject": "Updated",
		},
	}

	ticket, err := ticketsClient.UpdateTicket(context.Background(), "99999", input)

	require.Error(t, err)
	assert.Nil(t, ticket)
}

// TestUpdateTicket_InvalidJSON tests invalid JSON response
func TestUpdateTicket_InvalidJSON(t *testing.T) {
	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &UpdateTicketInput{
		Properties: map[string]string{
			"subject": "Updated",
		},
	}

	ticket, err := ticketsClient.UpdateTicket(context.Background(), "123456", input)

	require.Error(t, err)
	assert.Nil(t, ticket)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestArchiveTicket_Success tests successful ticket archival
func TestArchiveTicket_Success(t *testing.T) {
	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/123456", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := ticketsClient.ArchiveTicket(context.Background(), "123456")

	assert.NoError(t, err)
}

// TestArchiveTicket_NotFound tests 404 error on archive
func TestArchiveTicket_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Ticket not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	err := ticketsClient.ArchiveTicket(context.Background(), "99999")

	require.Error(t, err)
}

// TestMergeTwoTickets_Success tests successful ticket merge
func TestMergeTwoTickets_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"subject": "Merged Ticket"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"archived": false
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/merge", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &MergeTwoTicketsInput{
		PrimaryObjectID: "123456",
		ObjectIDToMerge: "789012",
	}

	err := ticketsClient.MergeTwoTickets(context.Background(), input)

	assert.NoError(t, err)
}

// TestMergeTwoTickets_Error tests error on merge
func TestMergeTwoTickets_Error(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Cannot merge tickets",
		"category": "VALIDATION_ERROR"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &MergeTwoTicketsInput{
		PrimaryObjectID: "123456",
		ObjectIDToMerge: "123456",
	}

	err := ticketsClient.MergeTwoTickets(context.Background(), input)

	require.Error(t, err)
}

// TestBatchReadTickets_Success tests successful batch read
func TestBatchReadTickets_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"subject": "Ticket 1"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/batch/read", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchReadTicketsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{
			{ID: "1"},
		},
		Properties:            []string{"subject"},
		PropertiesWithHistory: []string{},
	}

	result, err := ticketsClient.BatchReadTickets(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Results, 1)
}

// TestBatchReadTickets_InvalidJSON tests invalid JSON response
func TestBatchReadTickets_InvalidJSON(t *testing.T) {
	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &BatchReadTicketsInput{
		Inputs: []struct {
			ID string `json:"id" required:"yes"`
		}{{ID: "1"}},
		Properties:            []string{"subject"},
		PropertiesWithHistory: []string{},
	}

	result, err := ticketsClient.BatchReadTickets(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestSearchTickets_Success tests successful ticket search
func TestSearchTickets_Success(t *testing.T) {
	responseJSON := `{
		"total": 1,
		"results": [
			{
				"id": "1",
				"properties": {
					"subject": "Urgent Ticket"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		]
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/search", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &SearchTicketsInput{
		FilterGroups: []struct {
			Filters []struct {
				PropertyName string         `json:"propertyName" required:"yes"`
				Operator     FilterOperator `json:"operator" required:"yes"`
				HighValue    string         `json:"highValue"`
				Values       []string       `json:"values"`
				Value        string         `json:"value"`
			} `json:"filters" required:"yes"`
		}{},
		Properties: []string{"subject"},
		Limit:      10,
		After:      "",
		Sorts:      []string{},
	}

	result, err := ticketsClient.SearchTickets(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
}

// TestSearchTickets_InvalidJSON tests invalid JSON response
func TestSearchTickets_InvalidJSON(t *testing.T) {
	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &SearchTicketsInput{
		FilterGroups: []struct {
			Filters []struct {
				PropertyName string         `json:"propertyName" required:"yes"`
				Operator     FilterOperator `json:"operator" required:"yes"`
				HighValue    string         `json:"highValue"`
				Values       []string       `json:"values"`
				Value        string         `json:"value"`
			} `json:"filters" required:"yes"`
		}{},
		Properties: []string{"subject"},
		Limit:      10,
		After:      "",
		Sorts:      []string{},
	}

	result, err := ticketsClient.SearchTickets(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestBatchCreateTickets_Success tests successful batch create
func TestBatchCreateTickets_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"subject": "New Ticket"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/batch/create", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateTicketsInput{}

	result, err := ticketsClient.BatchCreateTickets(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, BatchStatus("COMPLETE"), result.Status)
	assert.Len(t, result.Results, 1)
}

// TestBatchUpdateTickets_Success tests successful batch update
func TestBatchUpdateTickets_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"subject": "Updated Ticket"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-02T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-02T00:00:00.000Z",
		"completedAt": "2024-01-02T00:00:05.000Z"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/batch/update", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchUpdateTicketsInput{}

	result, err := ticketsClient.BatchUpdateTickets(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, BatchStatus("COMPLETE"), result.Status)
	assert.Len(t, result.Results, 1)
}

// TestBatchCreateOrUpdateTickets_Success tests successful batch create or update
func TestBatchCreateOrUpdateTickets_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {
					"subject": "Created or Updated"
				},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/batch/createOrUpdate", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateOrUpdateTicketsInput{}

	result, err := ticketsClient.BatchCreateOrUpdateTickets(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, BatchStatus("COMPLETE"), result.Status)
	assert.Len(t, result.Results, 1)
}

// TestBatchArchiveTickets_Success tests successful batch archive
func TestBatchArchiveTickets_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, ticketsClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/tickets/batch/archive", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchArchiveTicketsInput{}

	result, err := ticketsClient.BatchArchiveTickets(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, BatchStatus("COMPLETE"), result.Status)
}
