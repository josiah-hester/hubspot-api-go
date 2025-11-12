package lists

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

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	apiClient, err := client.NewClient(
		client.WithTimeout(5 * time.Second),
	)
	require.NoError(t, err)

	listsClient := NewClient(apiClient)
	assert.NotNil(t, listsClient)
	assert.NotNil(t, listsClient.apiClient)
}

// TestGetListByID_Success tests successful list retrieval by ID
func TestGetListByID_Success(t *testing.T) {
	listJSON := `{
		"list": {
			"listId": "123",
			"name": "My List",
			"objectTypeId": "0-1",
			"processingType": "MANUAL",
			"processingStatus": "COMPLETE",
			"listVersion": 1,
			"size": 100,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-02T00:00:00Z",
			"createdById": "user-1",
			"updatedById": "user-1"
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/lists/123", r.URL.Path)
		respondJSON(w, http.StatusOK, listJSON)
	})
	defer server.Close()

	list, err := listClient.GetListByID(context.Background(), "123")

	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, "123", list.ListID)
	assert.Equal(t, "My List", list.Name)
	assert.Equal(t, "0-1", list.ObjectTypeID)
	assert.Equal(t, Manual, list.ProcessingType)
	assert.Equal(t, Complete, list.ProcessingStatus)
}

// TestGetListByID_WithIncludeFilters tests GetListByID with includeFilters option
func TestGetListByID_WithIncludeFilters(t *testing.T) {
	listJSON := `{
		"list": {
			"listId": "123",
			"name": "Filtered List",
			"objectTypeId": "0-1",
			"processingType": "DYNAMIC",
			"processingStatus": "COMPLETE",
			"listVersion": 1,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-02T00:00:00Z",
			"createdById": "user-1",
			"updatedById": "user-1",
			"filterBranch": {
				"filterBranchType": "AND",
				"filterBranchOperator": "string",
				"filters": [
					{
						"filterType": "PROPERTY",
						"property": "email",
						"operation": {}
					}
				]
			}
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("includeFilters"))
		respondJSON(w, http.StatusOK, listJSON)
	})
	defer server.Close()

	list, err := listClient.GetListByID(context.Background(), "123", WithIncludeFilters(true))

	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.NotNil(t, list.FilterBranch)
	assert.Equal(t, And, list.FilterBranch.FilterBranchType)
}

// TestGetListByID_NotFound tests 404 error handling
func TestGetListByID_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "List not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	list, err := listClient.GetListByID(context.Background(), "999")

	require.Error(t, err)
	assert.Nil(t, list)

	var notFoundErr *ListNotFoundError
	require.ErrorAs(t, err, &notFoundErr)
	assert.Equal(t, "999", notFoundErr.ListID)
}

// TestGetListByID_InvalidJSON tests JSON unmarshal error
func TestGetListByID_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	list, err := listClient.GetListByID(context.Background(), "123")

	require.Error(t, err)
	assert.Nil(t, list)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestGetListByName_Success tests successful list retrieval by name
func TestGetListByName_Success(t *testing.T) {
	listJSON := `{
		"list": {
			"listId": "456",
			"name": "Contact List",
			"objectTypeId": "0-1",
			"processingType": "MANUAL",
			"processingStatus": "COMPLETE",
			"listVersion": 1,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-02T00:00:00Z",
			"createdById": "user-1",
			"updatedById": "user-1"
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/lists/object-type-id/0-1/name/Contact List", r.URL.Path)
		respondJSON(w, http.StatusOK, listJSON)
	})
	defer server.Close()

	list, err := listClient.GetListByName(context.Background(), "0-1", "Contact List")

	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, "456", list.ListID)
	assert.Equal(t, "Contact List", list.Name)
}

// TestCreateList_Success tests successful list creation
func TestCreateList_Success(t *testing.T) {
	responseJSON := `{
		"list": {
			"listId": "789",
			"name": "New List",
			"objectTypeId": "0-1",
			"processingType": "MANUAL",
			"processingStatus": "PROCESSING",
			"listVersion": 1,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
			"createdById": "user-1",
			"updatedById": "user-1"
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/lists", r.URL.Path)

		var input ListCreateRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		assert.NoError(t, err)
		assert.Equal(t, "New List", input.Name)
		assert.Equal(t, "0-1", input.ObjectTypeID)
		assert.Equal(t, Manual, input.ProcessingType)

		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &ListCreateRequest{
		Name:           "New List",
		ObjectTypeID:   "0-1",
		ProcessingType: Manual,
	}

	list, err := listClient.CreateList(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, "789", list.ListID)
	assert.Equal(t, "New List", list.Name)
	assert.Equal(t, Processing, list.ProcessingStatus)
}

// TestCreateList_ValidationError tests validation error
func TestCreateList_ValidationError(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "Invalid name",
		"category": "VALIDATION_ERROR"
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, errorJSON)
	})
	defer server.Close()

	input := &ListCreateRequest{
		Name:           "",
		ObjectTypeID:   "0-1",
		ProcessingType: Manual,
	}

	list, err := listClient.CreateList(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, list)

	var validationErr *ListValidationError
	require.ErrorAs(t, err, &validationErr)
}

// TestGetListsByIDs_Success tests successful retrieval of multiple lists
func TestGetListsByIDs_Success(t *testing.T) {
	responseJSON := `{
		"lists": [
			{
				"listId": "1",
				"name": "List 1",
				"objectTypeId": "0-1",
				"processingType": "MANUAL",
				"processingStatus": "COMPLETE",
				"listVersion": 1,
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-02T00:00:00Z",
				"createdById": "user-1",
				"updatedById": "user-1"
			},
			{
				"listId": "2",
				"name": "List 2",
				"objectTypeId": "0-1",
				"processingType": "DYNAMIC",
				"processingStatus": "COMPLETE",
				"listVersion": 1,
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-02T00:00:00Z",
				"createdById": "user-1",
				"updatedById": "user-1"
			}
		]
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/lists", r.URL.Path)
		// Query params will have multiple listIds
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	lists, err := listClient.GetListsByIDs(context.Background(), []string{"1", "2"})

	require.NoError(t, err)
	assert.Len(t, lists, 2)
	assert.Equal(t, "1", lists[0].ListID)
	assert.Equal(t, "2", lists[1].ListID)
}

// TestSearchLists_Success tests successful list search
func TestSearchLists_Success(t *testing.T) {
	responseJSON := `{
		"lists": [
			{
				"listId": "1",
				"name": "Search Result",
				"objectTypeId": "0-1",
				"processingType": "MANUAL",
				"processingStatus": "COMPLETE",
				"listVersion": 1,
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-02T00:00:00Z",
				"createdById": "user-1",
				"updatedById": "user-1"
			}
		],
		"total": 1,
		"hasMore": false,
		"offset": 0
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/lists/search", r.URL.Path)

		var input ListSearchRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		assert.NoError(t, err)

		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	query := "Search"
	input := &ListSearchRequest{
		Query: &query,
	}

	result, err := listClient.SearchLists(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Lists, 1)
	assert.Equal(t, 1, result.Total)
	assert.False(t, result.HasMore)
}

// TestUpdateListName_Success tests successful list name update
func TestUpdateListName_Success(t *testing.T) {
	responseJSON := `{
		"list": {
			"listId": "123",
			"name": "Updated Name",
			"objectTypeId": "0-1",
			"processingType": "MANUAL",
			"processingStatus": "COMPLETE",
			"listVersion": 2,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-03T00:00:00Z",
			"createdById": "user-1",
			"updatedById": "user-1"
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/update-list-name", r.URL.Path)
		assert.Equal(t, "Updated Name", r.URL.Query().Get("listName"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	list, err := listClient.UpdateListName(context.Background(), "123", "Updated Name", false)

	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, "Updated Name", list.Name)
	assert.Equal(t, 2, list.ListVersion)
}

// TestUpdateListName_WithIncludeFilters tests update with includeFilters
func TestUpdateListName_WithIncludeFilters(t *testing.T) {
	responseJSON := `{
		"list": {
			"listId": "123",
			"name": "Updated Name",
			"objectTypeId": "0-1",
			"processingType": "DYNAMIC",
			"processingStatus": "COMPLETE",
			"listVersion": 2,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-03T00:00:00Z",
			"createdById": "user-1",
			"updatedById": "user-1",
			"filterBranch": {
				"filterBranchType": "AND",
				"filterBranchOperator": "string",
				"filters": []
			}
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("includeFilters"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	list, err := listClient.UpdateListName(context.Background(), "123", "Updated Name", true)

	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.NotNil(t, list.FilterBranch)
}

// TestUpdateListFilters_Success tests successful list filter update
func TestUpdateListFilters_Success(t *testing.T) {
	responseJSON := `{
		"list": {
			"listId": "123",
			"name": "Filtered List",
			"objectTypeId": "0-1",
			"processingType": "DYNAMIC",
			"processingStatus": "PROCESSING",
			"listVersion": 2,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-03T00:00:00Z",
			"createdById": "user-1",
			"updatedById": "user-1"
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/update-list-filters", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	filterBranch := FilterBranch{
		FilterBranchType: And,
		Filters:          []Filter{},
	}

	list, err := listClient.UpdateListFilters(context.Background(), "123", filterBranch, false)

	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 2, list.ListVersion)
}

// TestDeleteList_Success tests successful list deletion
func TestDeleteList_Success(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/lists/123", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := listClient.DeleteList(context.Background(), "123")

	assert.NoError(t, err)
}

// TestDeleteList_NotFound tests 404 error on delete
func TestDeleteList_NotFound(t *testing.T) {
	errorJSON := `{
		"status": "error",
		"message": "List not found",
		"category": "OBJECT_NOT_FOUND"
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, errorJSON)
	})
	defer server.Close()

	err := listClient.DeleteList(context.Background(), "999")

	require.Error(t, err)

	var notFoundErr *ListNotFoundError
	require.ErrorAs(t, err, &notFoundErr)
	assert.Equal(t, "999", notFoundErr.ListID)
}

// TestRestoreList_Success tests successful list restoration
func TestRestoreList_Success(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/restore", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := listClient.RestoreList(context.Background(), "123")

	assert.NoError(t, err)
}

// TestGetRecordMemberships_Success tests successful retrieval of record memberships
func TestGetRecordMemberships_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"listId": "1",
				"listVersion": 1,
				"firstAddedTimestamp": "2024-01-01T00:00:00Z",
				"lastAddedTimestamp": "2024-01-02T00:00:00Z"
			}
		],
		"total": 1
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/lists/records/0-1/contact-123/memberships", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	memberships, err := listClient.GetRecordMemberships(context.Background(), "0-1", "contact-123")

	require.NoError(t, err)
	assert.NotNil(t, memberships)
	assert.Len(t, memberships.Results, 1)
	assert.Equal(t, "1", memberships.Results[0].ListID)
}

// TestBatchGetRecordMemberships_Success tests batch retrieval of record memberships
func TestBatchGetRecordMemberships_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"results": [
					{
						"listId": "1",
						"listVersion": 1,
						"firstAddedTimestamp": "2024-01-01T00:00:00Z",
						"lastAddedTimestamp": "2024-01-02T00:00:00Z"
					}
				],
				"total": 1
			}
		]
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/lists/records/memberships/batch/read", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	inputs := []MembershipRecordIdentifier{
		{ObjectTypeID: "0-1", RecordID: "contact-123"},
	}

	result, err := listClient.BatchGetRecordMemberships(context.Background(), inputs)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Results, 1)
}

// TestAddRecordsToList_Success tests successfully adding records to a list
func TestAddRecordsToList_Success(t *testing.T) {
	responseJSON := `{
		"recordIdsAdded": ["contact-1", "contact-2"]
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/memberships/add", r.URL.Path)

		var recordIDs []string
		err := json.NewDecoder(r.Body).Decode(&recordIDs)
		assert.NoError(t, err)
		assert.Len(t, recordIDs, 2)

		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	result, err := listClient.AddRecordsToList(context.Background(), "123", []string{"contact-1", "contact-2"})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.RecordIDsAdded, 2)
}

// TestAddFromSourceList_Success tests successfully adding from source list
func TestAddFromSourceList_Success(t *testing.T) {
	responseJSON := `{
		"recordIdsAdded": ["contact-1", "contact-2", "contact-3"]
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/memberships/add-from/456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	result, err := listClient.AddFromSourceList(context.Background(), "123", "456")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.RecordIDsAdded, 3)
}

// TestGetListMemberships_Success tests successful retrieval of list memberships
func TestGetListMemberships_Success(t *testing.T) {
	responseJSON := `{
		"results": ["contact-1", "contact-2", "contact-3"],
		"hasMore": true,
		"offset": "next-page-token"
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/memberships", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	memberships, err := listClient.GetListMemberships(context.Background(), "123")

	require.NoError(t, err)
	assert.NotNil(t, memberships)
	assert.Len(t, memberships.Results, 3)
	assert.NotNil(t, memberships.HasMore)
	assert.True(t, *memberships.HasMore)
}

// TestGetListMemberships_WithOptions tests GetListMemberships with pagination options
func TestGetListMemberships_WithOptions(t *testing.T) {
	responseJSON := `{
		"results": ["contact-1", "contact-2"],
		"hasMore": false
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "token123", r.URL.Query().Get("offset"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	memberships, err := listClient.GetListMemberships(
		context.Background(),
		"123",
		WithMembershipsLimit(10),
		WithMembershipsOffset("token123"),
	)

	require.NoError(t, err)
	assert.NotNil(t, memberships)
	assert.Len(t, memberships.Results, 2)
}

// TestRemoveAllRecords_Success tests successfully removing all records from a list
func TestRemoveAllRecords_Success(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/memberships", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := listClient.RemoveAllRecords(context.Background(), "123")

	assert.NoError(t, err)
}

// TestRemoveRecordsFromList_Success tests successfully removing specific records
func TestRemoveRecordsFromList_Success(t *testing.T) {
	responseJSON := `{
		"recordIdsRemoved": ["contact-1", "contact-2"]
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/memberships/remove", r.URL.Path)

		var recordIDs []string
		err := json.NewDecoder(r.Body).Decode(&recordIDs)
		assert.NoError(t, err)
		assert.Len(t, recordIDs, 2)

		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	result, err := listClient.RemoveRecordsFromList(context.Background(), "123", []string{"contact-1", "contact-2"})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.RecordIDsRemoved, 2)
}

// TestScheduleConversion_Success tests successfully scheduling a list conversion
func TestScheduleConversion_Success(t *testing.T) {
	responseJSON := `{
		"listId": "123",
		"requestedConversionTime": {
			"conversionType": "CONVERSION_DATE",
			"year": 2024,
			"month": 12,
			"day": 31
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/schedule-conversion", r.URL.Path)

		var input ScheduleConversionRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		assert.NoError(t, err)
		assert.Equal(t, ConversionDate, input.ConversionType)

		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	year := 2024
	month := 12
	day := 31
	input := &ScheduleConversionRequest{
		ConversionType: ConversionDate,
		Year:           &year,
		Month:          &month,
		Day:            &day,
	}

	result, err := listClient.ScheduleConversion(context.Background(), "123", input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123", result.ListID)
	assert.Equal(t, ConversionDate, result.RequestedConversionTime.ConversionType)
}

// TestGetConversionSchedule_Success tests retrieving conversion schedule
func TestGetConversionSchedule_Success(t *testing.T) {
	responseJSON := `{
		"listId": "123",
		"requestedConversionTime": {
			"conversionType": "INACTIVITY",
			"timeUnit": "DAY",
			"offset": 30
		}
	}`

	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/schedule-conversion", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	result, err := listClient.GetConversionSchedule(context.Background(), "123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123", result.ListID)
	assert.Equal(t, Inactivity, result.RequestedConversionTime.ConversionType)
}

// TestDeleteConversionSchedule_Success tests deleting conversion schedule
func TestDeleteConversionSchedule_Success(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/lists/123/schedule-conversion", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := listClient.DeleteConversionSchedule(context.Background(), "123")

	assert.NoError(t, err)
}

// Additional tests for improving coverage to 95%+

// TestGetListByName_InvalidJSON tests JSON unmarshal error
func TestGetListByName_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	list, err := listClient.GetListByName(context.Background(), "0-1", "Test List")

	require.Error(t, err)
	assert.Nil(t, list)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestCreateList_InvalidJSON tests JSON unmarshal error
func TestCreateList_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	input := &ListCreateRequest{
		Name:           "New List",
		ObjectTypeID:   "0-1",
		ProcessingType: Manual,
	}

	list, err := listClient.CreateList(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, list)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestGetListsByIDs_InvalidJSON tests JSON unmarshal error
func TestGetListsByIDs_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	lists, err := listClient.GetListsByIDs(context.Background(), []string{"1", "2"})

	require.Error(t, err)
	assert.Nil(t, lists)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestSearchLists_InvalidJSON tests JSON unmarshal error
func TestSearchLists_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	query := "test"
	input := &ListSearchRequest{
		Query: &query,
	}

	result, err := listClient.SearchLists(context.Background(), input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestUpdateListName_InvalidJSON tests JSON unmarshal error
func TestUpdateListName_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	list, err := listClient.UpdateListName(context.Background(), "123", "New Name", false)

	require.Error(t, err)
	assert.Nil(t, list)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestUpdateListFilters_InvalidJSON tests JSON unmarshal error
func TestUpdateListFilters_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	filterBranch := FilterBranch{
		FilterBranchType: And,
		Filters:          []Filter{},
	}

	list, err := listClient.UpdateListFilters(context.Background(), "123", filterBranch, false)

	require.Error(t, err)
	assert.Nil(t, list)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestGetRecordMemberships_InvalidJSON tests JSON unmarshal error
func TestGetRecordMemberships_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	memberships, err := listClient.GetRecordMemberships(context.Background(), "0-1", "contact-123")

	require.Error(t, err)
	assert.Nil(t, memberships)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestBatchGetRecordMemberships_InvalidJSON tests JSON unmarshal error
func TestBatchGetRecordMemberships_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	inputs := []MembershipRecordIdentifier{
		{ObjectTypeID: "0-1", RecordID: "contact-123"},
	}

	result, err := listClient.BatchGetRecordMemberships(context.Background(), inputs)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestAddRecordsToList_InvalidJSON tests JSON unmarshal error
func TestAddRecordsToList_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	result, err := listClient.AddRecordsToList(context.Background(), "123", []string{"contact-1"})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestAddFromSourceList_InvalidJSON tests JSON unmarshal error
func TestAddFromSourceList_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	result, err := listClient.AddFromSourceList(context.Background(), "123", "456")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestGetListMemberships_InvalidJSON tests JSON unmarshal error
func TestGetListMemberships_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	memberships, err := listClient.GetListMemberships(context.Background(), "123")

	require.Error(t, err)
	assert.Nil(t, memberships)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestRemoveRecordsFromList_InvalidJSON tests JSON unmarshal error
func TestRemoveRecordsFromList_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	result, err := listClient.RemoveRecordsFromList(context.Background(), "123", []string{"contact-1"})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestScheduleConversion_InvalidJSON tests JSON unmarshal error
func TestScheduleConversion_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	year := 2024
	input := &ScheduleConversionRequest{
		ConversionType: ConversionDate,
		Year:           &year,
	}

	result, err := listClient.ScheduleConversion(context.Background(), "123", input)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

// TestGetConversionSchedule_InvalidJSON tests JSON unmarshal error
func TestGetConversionSchedule_InvalidJSON(t *testing.T) {
	server, listClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	})
	defer server.Close()

	result, err := listClient.GetConversionSchedule(context.Background(), "123")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}
