package orders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josiah-hester/go-hubspot-sdk/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions
func setupMockServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(handler))

	apiClient, err := client.NewClient(
		client.WithBaseURL(server.URL),
		client.WithAccessToken("test-token"),
		client.WithRateLimitEnabled(false),
		client.WithRetryEnabled(false),
	)
	require.NoError(t, err)

	ordersClient := NewClient(apiClient)
	return server, ordersClient
}

func respondJSON(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	apiClient, err := client.NewClient()
	require.NoError(t, err)

	ordersClient := NewClient(apiClient)
	assert.NotNil(t, ordersClient)
	assert.NotNil(t, ordersClient.apiClient)
}

// TestCreateOrder tests order creation
func TestCreateOrder_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"hs_order_name": "Order #12345",
			"amount": "1500"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders", r.URL.Path)
		respondJSON(w, http.StatusCreated, responseJSON)
	})
	defer server.Close()

	input := &CreateOrderInput{
		Properties: map[string]string{
			"hs_order_name": "Order #12345",
			"amount":        "1500",
		},
	}

	order, err := ordersClient.CreateOrder(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "123456", order.ID)
	assert.Equal(t, "Order #12345", order.Properties["hs_order_name"])
}

func TestCreateOrder_Error(t *testing.T) {
	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, `{"status": "error", "message": "Invalid input"}`)
	})
	defer server.Close()

	input := &CreateOrderInput{Properties: map[string]string{}}
	_, err := ordersClient.CreateOrder(context.Background(), input)

	require.Error(t, err)
}

// TestGetOrder tests retrieving an order
func TestGetOrder_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"hs_order_name": "Order #12345",
			"amount": "1500"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/123456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	order, err := ordersClient.GetOrder(context.Background(), "123456")

	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "123456", order.ID)
}

func TestGetOrder_WithOptions(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"hs_order_name": "Test Order"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-01T00:00:00.000Z",
		"archived": false
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "hs_order_name,amount", r.URL.Query().Get("properties"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	order, err := ordersClient.GetOrder(context.Background(), "123456",
		WithProperties([]string{"hs_order_name", "amount"}))

	require.NoError(t, err)
	assert.NotNil(t, order)
}

func TestGetOrder_NotFound(t *testing.T) {
	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, `{"status": "error", "message": "Not found"}`)
	})
	defer server.Close()

	_, err := ordersClient.GetOrder(context.Background(), "999999")
	require.Error(t, err)
}

// TestUpdateOrder tests updating an order
func TestUpdateOrder_Success(t *testing.T) {
	responseJSON := `{
		"id": "123456",
		"properties": {
			"hs_order_name": "Updated Order",
			"amount": "2000"
		},
		"createdAt": "2024-01-01T00:00:00.000Z",
		"updatedAt": "2024-01-02T00:00:00.000Z",
		"archived": false
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/123456", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &UpdateOrderInput{
		Properties: map[string]string{
			"hs_order_name": "Updated Order",
		},
	}

	order, err := ordersClient.UpdateOrder(context.Background(), "123456", input)

	require.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "Updated Order", order.Properties["hs_order_name"])
}

// TestArchiveOrder tests archiving an order
func TestArchiveOrder_Success(t *testing.T) {
	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/123456", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := ordersClient.ArchiveOrder(context.Background(), "123456")
	require.NoError(t, err)
}

// TestListOrders tests listing orders
func TestListOrders_Success(t *testing.T) {
	responseJSON := `{
		"results": [
			{
				"id": "1",
				"properties": {"hs_order_name": "Order A", "amount": "1000"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			},
			{
				"id": "2",
				"properties": {"hs_order_name": "Order B", "amount": "2500"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"paging": {
			"next": {
				"after": "abc123",
				"link": "?after=abc123"
			}
		}
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	resp, err := ordersClient.ListOrders(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
	assert.Equal(t, "abc123", resp.Paging.Next.After)
}

func TestListOrders_WithOptions(t *testing.T) {
	responseJSON := `{"results": [], "paging": null}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "50", r.URL.Query().Get("limit"))
		assert.Equal(t, "xyz789", r.URL.Query().Get("after"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	_, err := ordersClient.ListOrders(context.Background(),
		WithLimit(50),
		WithAfter("xyz789"))

	require.NoError(t, err)
}

// TestBatchReadOrders tests batch read
func TestBatchReadOrders_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {"hs_order_name": "Order A"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/batch/read", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchReadOrdersInput{
		Properties: []string{"hs_order_name", "amount"},
		Inputs: []struct {
			ID string `json:"id"`
		}{
			{ID: "1"},
		},
	}

	resp, err := ordersClient.BatchReadOrders(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "COMPLETE", resp.Status)
	assert.Len(t, resp.Results, 1)
}

// TestBatchCreateOrders tests batch create
func TestBatchCreateOrders_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [
			{
				"id": "1",
				"properties": {"hs_order_name": "New Order"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/batch/create", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchCreateOrdersInput{
		Inputs: []CreateOrderInput{
			{Properties: map[string]string{"hs_order_name": "New Order"}},
		},
	}

	resp, err := ordersClient.BatchCreateOrders(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "COMPLETE", resp.Status)
}

// TestBatchUpdateOrders tests batch update
func TestBatchUpdateOrders_Success(t *testing.T) {
	responseJSON := `{
		"status": "COMPLETE",
		"results": [],
		"startedAt": "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:05.000Z"
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/batch/update", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &BatchUpdateOrdersInput{}
	resp, err := ordersClient.BatchUpdateOrders(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestBatchArchiveOrders tests batch archive
func TestBatchArchiveOrders_Success(t *testing.T) {
	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/batch/archive", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	input := &BatchArchiveOrdersInput{}
	err := ordersClient.BatchArchiveOrders(context.Background(), input)

	require.NoError(t, err)
}

// TestSearchOrders tests search functionality
func TestSearchOrders_Success(t *testing.T) {
	responseJSON := `{
		"total": 1,
		"results": [
			{
				"id": "1",
				"properties": {"hs_order_name": "Order #12345", "amount": "5000"},
				"createdAt": "2024-01-01T00:00:00.000Z",
				"updatedAt": "2024-01-01T00:00:00.000Z",
				"archived": false
			}
		],
		"paging": null
	}`

	server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/crm/v3/objects/orders/search", r.URL.Path)
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	input := &SearchOrdersInput{
		FilterGroups: []FilterGroup{
			{
				Filters: []Filter{
					{
						PropertyName: "amount",
						Operator:     "GTE",
						Value:        "1000",
					},
				},
			},
		},
		Properties: []string{"hs_order_name", "amount"},
		Limit:      10,
	}

	resp, err := ordersClient.SearchOrders(context.Background(), input)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Total)
	assert.Len(t, resp.Results, 1)
}

// TestOptions tests all option functions
func TestOptions(t *testing.T) {
	tests := []struct {
		name     string
		option   OrderOption
		expected string
		param    string
	}{
		{"WithPropertiesWithHistory", WithPropertiesWithHistory([]string{"hs_order_name", "amount"}), "hs_order_name,amount", "propertiesWithHistory"},
		{"WithAssociations", WithAssociations([]string{"contacts", "deals"}), "contacts,deals", "associations"},
		{"WithIDProperty", WithIDProperty("orderid"), "orderid", "idProperty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.expected, r.URL.Query().Get(tt.param))
				respondJSON(w, http.StatusOK, `{"id": "123", "properties": {"hs_order_name": "Test"}, "createdAt": "2024-01-01T00:00:00.000Z", "updatedAt": "2024-01-01T00:00:00.000Z", "archived": false}`)
			})
			defer server.Close()

			_, err := ordersClient.GetOrder(context.Background(), "123", tt.option)
			require.NoError(t, err)
		})
	}

	t.Run("WithArchived", func(t *testing.T) {
		server, ordersClient := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "true", r.URL.Query().Get("archived"))
			respondJSON(w, http.StatusOK, `{"results": [], "paging": null}`)
		})
		defer server.Close()

		_, err := ordersClient.ListOrders(context.Background(), WithArchived())
		require.NoError(t, err)
	})
}
