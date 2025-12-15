package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helpers
func setupMockServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(handler))

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithAccessToken("test-token"),
		WithRateLimitEnabled(false),
		WithRetryEnabled(false),
	)
	require.NoError(t, err)

	return server, client
}

func respondJSON(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
}

// TestNewClient tests client creation with various options
func TestNewClient(t *testing.T) {
	t.Run("Default config", func(t *testing.T) {
		client, err := NewClient()
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "https://api.hubapi.com", client.config.BaseURL)
		assert.Equal(t, 30*time.Second, client.config.Timeout)
		assert.True(t, client.config.RateLimit.Enabled)
		assert.True(t, client.config.Retry.Enabled)
	})

	t.Run("With access token", func(t *testing.T) {
		client, err := NewClient(WithAccessToken("test-token-123"))
		require.NoError(t, err)
		assert.Equal(t, "test-token-123", client.config.AccessToken)
	})

	t.Run("With custom base URL", func(t *testing.T) {
		client, err := NewClient(WithBaseURL("https://custom.api.com"))
		require.NoError(t, err)
		assert.Equal(t, "https://custom.api.com", client.config.BaseURL)
	})

	t.Run("With timeout", func(t *testing.T) {
		client, err := NewClient(WithTimeout(60 * time.Second))
		require.NoError(t, err)
		assert.Equal(t, 60*time.Second, client.config.Timeout)
	})

	t.Run("With rate limit config", func(t *testing.T) {
		client, err := NewClient(
			WithRateLimitMaxBurst(200),
			WithRateLimitDailyLimit(500000),
			WithRateLimitEnabled(false),
		)
		require.NoError(t, err)
		assert.Equal(t, 200, client.config.RateLimit.MaxBurst)
		assert.Equal(t, 500000, client.config.RateLimit.DailyLimit)
		assert.False(t, client.config.RateLimit.Enabled)
	})

	t.Run("With retry config", func(t *testing.T) {
		client, err := NewClient(
			WithRetryMaxAttempts(5),
			WithRetryBackoff(2*time.Second, 60*time.Second),
			WithRetryEnabled(false),
		)
		require.NoError(t, err)
		assert.Equal(t, 5, client.config.Retry.MaxAttempts)
		assert.Equal(t, 2*time.Second, client.config.Retry.InitialBackoff)
		assert.Equal(t, 60*time.Second, client.config.Retry.MaxBackoff)
		assert.False(t, client.config.Retry.Enabled)
	})
}

// TestClientDo_Success tests successful HTTP requests
func TestClientDo_Success(t *testing.T) {
	responseJSON := `{"id": "123", "name": "Test"}`

	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/test/path", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		respondJSON(w, http.StatusOK, responseJSON)
	})
	defer server.Close()

	req := NewRequest("GET", "/test/path")
	resp, err := client.Do(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.JSONEq(t, responseJSON, string(resp.Body))
}

// TestClientDo_WithBody tests POST request with body
func TestClientDo_WithBody(t *testing.T) {
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body map[string]string
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "test", body["key"])

		respondJSON(w, http.StatusCreated, `{"success": true}`)
	})
	defer server.Close()

	req := NewRequest("POST", "/test/create")
	req.WithBody(map[string]string{"key": "test"})

	resp, err := client.Do(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

// TestClientDo_WithQueryParams tests request with query parameters
func TestClientDo_WithQueryParams(t *testing.T) {
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value1", r.URL.Query().Get("param1"))
		assert.Equal(t, "value2", r.URL.Query().Get("param2"))
		respondJSON(w, http.StatusOK, `{"success": true}`)
	})
	defer server.Close()

	req := NewRequest("GET", "/test")
	req.AddQueryParam("param1", "value1")
	req.AddQueryParam("param2", "value2")

	resp, err := client.Do(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// TestClientDo_WithHeaders tests custom headers
func TestClientDo_WithHeaders(t *testing.T) {
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))
		respondJSON(w, http.StatusOK, `{"success": true}`)
	})
	defer server.Close()

	req := NewRequest("GET", "/test")
	req.AddHeader("X-Custom-Header", "custom-value")

	resp, err := client.Do(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// TestClientDo_ErrorResponses tests various error status codes
func TestClientDo_ErrorResponses(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		errorJSON  string
	}{
		{
			name:       "400 Bad Request",
			statusCode: 400,
			errorJSON:  `{"status": "error", "message": "Bad request", "category": "VALIDATION_ERROR"}`,
		},
		{
			name:       "404 Not Found",
			statusCode: 404,
			errorJSON:  `{"status": "error", "message": "Not found", "category": "OBJECT_NOT_FOUND"}`,
		},
		{
			name:       "500 Server Error",
			statusCode: 500,
			errorJSON:  `{"status": "error", "message": "Server error", "category": "INTERNAL_ERROR"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				respondJSON(w, tc.statusCode, tc.errorJSON)
			})
			defer server.Close()

			req := NewRequest("GET", "/test")
			resp, err := client.Do(context.Background(), req)

			require.Error(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			hubspotErr, ok := err.(*HubSpotError)
			require.True(t, ok)
			assert.Equal(t, tc.statusCode, hubspotErr.Status)
		})
	}
}

// TestClientDo_RateLimitError tests rate limit error handling
func TestClientDo_RateLimitError(t *testing.T) {
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		respondJSON(w, 429, `{"status": "error", "message": "Rate limit exceeded", "policyName": "TEN_SECONDLY_ROLLING"}`)
	})
	defer server.Close()

	req := NewRequest("GET", "/test")
	_, err := client.Do(context.Background(), req)

	require.Error(t, err)
	hubspotErr, ok := err.(*HubSpotError)
	require.True(t, ok)
	assert.Equal(t, 429, hubspotErr.Status)
	assert.Equal(t, 60*time.Second, hubspotErr.RetryAfter)
	assert.True(t, hubspotErr.IsRetryable)
}

// TestClientDo_ContextCancellation tests context cancellation
func TestClientDo_ContextCancellation(t *testing.T) {
	server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		respondJSON(w, http.StatusOK, `{"success": true}`)
	})
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := NewRequest("GET", "/test")
	_, err := client.Do(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

// TestAuthMiddleware tests authentication header injection
func TestAuthMiddleware(t *testing.T) {
	t.Run("With access token", func(t *testing.T) {
		server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			respondJSON(w, http.StatusOK, `{"success": true}`)
		})
		defer server.Close()

		req := NewRequest("GET", "/test")
		_, err := client.Do(context.Background(), req)
		require.NoError(t, err)
	})

	t.Run("Without access token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.Header.Get("Authorization"))
			respondJSON(w, http.StatusOK, `{"success": true}`)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithRateLimitEnabled(false),
			WithRetryEnabled(false),
		)
		require.NoError(t, err)

		req := NewRequest("GET", "/test")
		resp, err := client.Do(context.Background(), req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

// TestRetryMiddleware tests retry logic
func TestRetryMiddleware(t *testing.T) {
	t.Run("Retry on retryable error", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				respondJSON(w, 500, `{"status": "error", "message": "Server error"}`)
			} else {
				respondJSON(w, 200, `{"success": true}`)
			}
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithRateLimitEnabled(false),
			WithRetryEnabled(true),
			WithRetryMaxAttempts(3),
			WithRetryBackoff(10*time.Millisecond, 100*time.Millisecond),
		)
		require.NoError(t, err)

		req := NewRequest("GET", "/test")
		resp, err := client.Do(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, 3, attempts)
	})

	t.Run("No retry on non-retryable error", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			respondJSON(w, 400, `{"status": "error", "message": "Bad request"}`)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithRateLimitEnabled(false),
			WithRetryEnabled(true),
			WithRetryMaxAttempts(3),
		)
		require.NoError(t, err)

		req := NewRequest("GET", "/test")
		_, err = client.Do(context.Background(), req)

		require.Error(t, err)
		assert.Equal(t, 1, attempts) // Should not retry
	})

	t.Run("Retry disabled", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			respondJSON(w, 500, `{"status": "error", "message": "Server error"}`)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithRateLimitEnabled(false),
			WithRetryEnabled(false),
		)
		require.NoError(t, err)

		req := NewRequest("GET", "/test")
		_, err = client.Do(context.Background(), req)

		require.Error(t, err)
		assert.Equal(t, 1, attempts) // No retry
	})
}

// TestRateLimitMiddleware tests rate limiting
func TestRateLimitMiddleware(t *testing.T) {
	t.Run("Rate limit enabled", func(t *testing.T) {
		server, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-HubSpot-RateLimit-Max", "100")
			w.Header().Set("X-HubSpot-RateLimit-Remaining", "99")
			w.Header().Set("X-HubSpot-RateLimit-Daily", "250000")
			w.Header().Set("X-HubSpot-RateLimit-Daily-Remaining", "249999")
			respondJSON(w, http.StatusOK, `{"success": true}`)
		})
		defer server.Close()

		// Re-enable rate limiting for this test
		client.config.RateLimit.Enabled = true

		req := NewRequest("GET", "/test")
		resp, err := client.Do(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, 99, resp.RateLimit.Remaining)
		assert.Equal(t, 249999, resp.RateLimit.DailyRemaining)
	})

	t.Run("Daily limit exceeded", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusOK, `{"success": true}`)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithRateLimitEnabled(true),
			WithRetryEnabled(false),
		)
		require.NoError(t, err)

		// Simulate daily limit exceeded
		client.rateLimiter.dailyRemaining = 0

		req := NewRequest("GET", "/test")
		_, err = client.Do(context.Background(), req)

		require.Error(t, err)
		hubspotErr, ok := err.(*HubSpotError)
		require.True(t, ok)
		assert.Equal(t, 429, hubspotErr.Status)
		assert.Contains(t, hubspotErr.Message, "Daily API limit exceeded")
	})
}

// TestMarshalRequestBody tests body marshaling
func TestMarshalRequestBody(t *testing.T) {
	t.Run("JSON object", func(t *testing.T) {
		body := map[string]string{"key": "value"}
		bytes, err := marshalRequestBody(body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"key":"value"}`, string(bytes))
	})

	t.Run("Byte slice", func(t *testing.T) {
		body := []byte("raw bytes")
		bytes, err := marshalRequestBody(body)
		require.NoError(t, err)
		assert.Equal(t, "raw bytes", string(bytes))
	})

	t.Run("String", func(t *testing.T) {
		body := "plain string"
		bytes, err := marshalRequestBody(body)
		require.NoError(t, err)
		assert.Equal(t, "plain string", string(bytes))
	})
}

// TestCalculateBackoffDuration tests backoff calculation
func TestCalculateBackoffDuration(t *testing.T) {
	cfg := RetryConfig{
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
	}

	t.Run("First attempt", func(t *testing.T) {
		backoff := calculateBackoffDuration(0, 0, cfg)
		// Should be around 1 second (with jitter)
		assert.Greater(t, backoff, 900*time.Millisecond)
		assert.Less(t, backoff, 1100*time.Millisecond)
	})

	t.Run("With Retry-After", func(t *testing.T) {
		retryAfter := 5 * time.Second
		backoff := calculateBackoffDuration(0, retryAfter, cfg)
		assert.Equal(t, 5*time.Second, backoff)
	})

	t.Run("Capped at max", func(t *testing.T) {
		backoff := calculateBackoffDuration(10, 0, cfg)
		assert.LessOrEqual(t, backoff, cfg.MaxBackoff)
	})
}

// TestRequest tests Request helper methods
func TestRequest(t *testing.T) {
	t.Run("NewRequest", func(t *testing.T) {
		req := NewRequest("GET", "/test")
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/test", req.Path)
		assert.NotNil(t, req.QueryParams)
		assert.NotNil(t, req.Headers)
	})

	t.Run("WithContext", func(t *testing.T) {
		type Key string
		var key Key = "test"
		ctx := context.WithValue(context.Background(), key, "value")
		req := NewRequest("GET", "/test").WithContext(ctx)
		assert.Equal(t, ctx, req.Context)
	})

	t.Run("WithResourceType", func(t *testing.T) {
		req := NewRequest("GET", "/test").WithResourceType("contacts")
		assert.Equal(t, "contacts", req.ResourceType)
	})

	t.Run("WithBody", func(t *testing.T) {
		body := map[string]string{"key": "value"}
		req := NewRequest("POST", "/test").WithBody(body)
		assert.Equal(t, body, req.Body)
	})

	t.Run("AddQueryParam", func(t *testing.T) {
		req := NewRequest("GET", "/test").AddQueryParam("limit", "10")
		assert.Equal(t, "10", req.QueryParams["limit"])
	})

	t.Run("AddHeader", func(t *testing.T) {
		req := NewRequest("GET", "/test").AddHeader("X-Custom", "value")
		assert.Equal(t, "value", req.Headers["X-Custom"])
	})
}

// TestResponse tests Response methods
func TestResponse(t *testing.T) {
	t.Run("NewResponse", func(t *testing.T) {
		headers := http.Header{}
		resp := NewResponse(200, []byte("test"), headers)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, []byte("test"), resp.Body)
	})

	t.Run("IsRateLimited - remaining zero", func(t *testing.T) {
		resp := &Response{
			RateLimit: RateLimitInfo{
				Remaining:      0,
				DailyRemaining: 100,
			},
		}
		assert.True(t, resp.IsRateLimited())
	})

	t.Run("IsRateLimited - daily remaining zero", func(t *testing.T) {
		resp := &Response{
			RateLimit: RateLimitInfo{
				Remaining:      100,
				DailyRemaining: 0,
			},
		}
		assert.True(t, resp.IsRateLimited())
	})

	t.Run("IsRateLimited - not limited", func(t *testing.T) {
		resp := &Response{
			RateLimit: RateLimitInfo{
				Remaining:      100,
				DailyRemaining: 10000,
			},
		}
		assert.False(t, resp.IsRateLimited())
	})
}

// TestHubSpotError tests HubSpotError methods
func TestHubSpotError(t *testing.T) {
	t.Run("Error with message", func(t *testing.T) {
		err := &HubSpotError{
			Status:    400,
			Message:   "Invalid request",
			ErrorType: "VALIDATION_ERROR",
		}
		assert.Contains(t, err.Error(), "Invalid request")
		assert.Contains(t, err.Error(), "VALIDATION_ERROR")
		assert.Contains(t, err.Error(), "400")
	})

	t.Run("Error without message", func(t *testing.T) {
		err := &HubSpotError{
			Status: 500,
		}
		assert.Contains(t, err.Error(), "status 500")
	})
}

// TestParseHubSpotError tests error parsing
func TestParseHubSpotError(t *testing.T) {
	t.Run("Parse complete error", func(t *testing.T) {
		body := []byte(`{
			"status": "error",
			"message": "Test error",
			"errorType": "VALIDATION_ERROR",
			"category": "VALIDATION",
			"policyName": "DAILY",
			"correlationId": "abc-123"
		}`)
		headers := http.Header{}

		err := ParseHubSpotError(400, body, headers)
		assert.Equal(t, 400, err.Status)
		assert.Equal(t, "Test error", err.Message)
		assert.Equal(t, "VALIDATION_ERROR", err.ErrorType)
		assert.Equal(t, "VALIDATION", err.Category)
		assert.Equal(t, "DAILY", err.PolicyName)
		assert.Equal(t, "abc-123", err.CorrelationID)
		assert.False(t, err.IsRetryable) // 400 is not retryable
	})

	t.Run("Parse with Retry-After header (seconds)", func(t *testing.T) {
		body := []byte(`{"status": "error", "message": "Rate limited"}`)
		headers := http.Header{}
		headers.Set("Retry-After", "60")

		err := ParseHubSpotError(429, body, headers)
		assert.Equal(t, 60*time.Second, err.RetryAfter)
	})

	t.Run("Parse with Retry-After header (RFC1123)", func(t *testing.T) {
		future := time.Now().Add(30 * time.Second)
		body := []byte(`{"status": "error", "message": "Rate limited"}`)
		headers := http.Header{}
		headers.Set("Retry-After", future.Format(time.RFC1123))

		err := ParseHubSpotError(429, body, headers)
		// Should be around 30 seconds
		assert.Greater(t, err.RetryAfter, 25*time.Second)
		assert.Less(t, err.RetryAfter, 35*time.Second)
	})

	t.Run("Parse invalid JSON", func(t *testing.T) {
		body := []byte(`invalid json`)
		headers := http.Header{}

		err := ParseHubSpotError(500, body, headers)
		assert.Equal(t, 500, err.Status)
		assert.Equal(t, "invalid json", err.RawBody)
	})

	t.Run("Retryable status codes", func(t *testing.T) {
		testCases := []struct {
			status      int
			policyName  string
			isRetryable bool
		}{
			{429, "TEN_SECONDLY_ROLLING", true},
			{429, "DAILY", false},
			{500, "", true},
			{502, "", true},
			{503, "", true},
			{504, "", true},
			{400, "", false},
			{404, "", false},
		}

		for _, tc := range testCases {
			body := []byte(`{"policyName": "` + tc.policyName + `"}`)
			err := ParseHubSpotError(tc.status, body, http.Header{})
			assert.Equal(t, tc.isRetryable, err.IsRetryable,
				"Status %d with policy %s should have IsRetryable=%v",
				tc.status, tc.policyName, tc.isRetryable)
		}
	})
}

// TestExtractRateLimitInfo tests rate limit header parsing
func TestExtractRateLimitInfo(t *testing.T) {
	t.Run("Extract all headers", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-HubSpot-RateLimit-Max", "100")
		headers.Set("X-HubSpot-RateLimit-Remaining", "95")
		headers.Set("X-HubSpot-RateLimit-Interval-Milliseconds", "10000")
		headers.Set("X-HubSpot-RateLimit-Daily", "250000")
		headers.Set("X-HubSpot-RateLimit-Daily-Remaining", "249500")

		info := ExtractRateLimitInfo(headers)
		assert.Equal(t, 100, info.Max)
		assert.Equal(t, 95, info.Remaining)
		assert.Equal(t, 10000, info.IntervalMs)
		assert.Equal(t, 250000, info.DailyLimit)
		assert.Equal(t, 249500, info.DailyRemaining)
	})

	t.Run("Extract with missing headers", func(t *testing.T) {
		headers := http.Header{}
		info := ExtractRateLimitInfo(headers)
		assert.Equal(t, 0, info.Max)
		assert.Equal(t, 0, info.Remaining)
		assert.Equal(t, 10000, info.IntervalMs) // Default value
	})

	t.Run("Extract with invalid values", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-HubSpot-RateLimit-Max", "invalid")
		headers.Set("X-HubSpot-RateLimit-Remaining", "not-a-number")

		info := ExtractRateLimitInfo(headers)
		assert.Equal(t, 0, info.Max)
		assert.Equal(t, 0, info.Remaining)
	})
}

// TestRateLimiter tests RateLimiter methods
func TestRateLimiter(t *testing.T) {
	t.Run("GetDailyRemaining", func(t *testing.T) {
		rl := NewRateLimiter(100)
		rl.dailyRemaining = 12345
		assert.Equal(t, 12345, rl.GetDailyRemaining())
	})

	t.Run("CheckDailyLimit - has quota", func(t *testing.T) {
		rl := NewRateLimiter(100)
		rl.dailyRemaining = 1000
		assert.True(t, rl.CheckDailyLimit())
	})

	t.Run("CheckDailyLimit - no quota", func(t *testing.T) {
		rl := NewRateLimiter(100)
		rl.dailyRemaining = 0
		assert.False(t, rl.CheckDailyLimit())
	})

	t.Run("UpdateFromResponse", func(t *testing.T) {
		rl := NewRateLimiter(100)
		resp := &Response{
			RateLimit: RateLimitInfo{
				DailyLimit:     250000,
				DailyRemaining: 249000,
			},
		}
		rl.UpdateFromResponse(resp)
		assert.Equal(t, 250000, rl.dailyLimit)
		assert.Equal(t, 249000, rl.dailyRemaining)
	})
}
