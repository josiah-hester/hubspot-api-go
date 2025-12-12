// Package client provides a core client for the HubSpot API
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client represents a HubSpot API client
type Client struct {
	config      *Config
	httpClient  *http.Client
	rateLimiter *RateLimiter
	logger      *Logger
}

// Handler represents a function that processes a Request and returns a Response
type Handler func(req *Request) (*Response, error)

func NewClient(opts ...Option) (*Client, error) {
	cfg := NewConfig()

	// Apply all options
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	// Validate required config
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	// Create rate limiter
	rateLimiter := NewRateLimiter(cfg.RateLimit.MaxBurst)

	// Initialize the logger
	logger := NewLogger(cfg.LoggingEnabled, "go_hubspot_sdk_client:", log.LstdFlags|log.Lshortfile, cfg.LogOutputs...)

	return &Client{
		config:      cfg,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		logger:      logger,
	}, nil
}

func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	req.Context = ctx

	// Build and execute middleware chain
	chain := c.buildChain()
	return chain(req)
}

// buildChain constructs the complete middleware chain
func (c *Client) buildChain() Handler {
	// Start with the HTTP handler (innermost)
	handler := c.httpMiddleware()

	// Wrap with retry middleware
	handler = c.wrapRetryMiddleware(handler)

	// Wrap with rate limit middleware
	handler = c.wrapRateLimitMiddleware(handler)

	// Wrap with auth middleware
	handler = c.wrapAuthMiddleware(handler)

	return handler
}

// wrapAuthMiddleware wraps a handler with authentication
func (c *Client) wrapAuthMiddleware(next Handler) Handler {
	return func(req *Request) (*Response, error) {
		if c.config.AccessToken != "" {
			req.AddHeader("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
		}
		return next(req)
	}
}

// wrapRateLimitMiddleware wraps a handler with rate limiting
func (c *Client) wrapRateLimitMiddleware(next Handler) Handler {
	return func(req *Request) (*Response, error) {
		if !c.config.RateLimit.Enabled {
			return next(req)
		}

		if !c.rateLimiter.CheckDailyLimit() {
			return nil, &HubSpotError{
				Status:      429,
				Message:     "Daily API limit exceeded",
				ErrorType:   "RATE_LIMIT",
				PolicyName:  "DAILY",
				IsRetryable: false,
			}
		}

		if err := c.rateLimiter.Wait(req.Context); err != nil {
			return nil, err
		}

		resp, err := next(req)

		if resp != nil {
			c.rateLimiter.UpdateFromResponse(resp)
		}

		return resp, err
	}
}

// wrapRetryMiddleware wraps a handler with retry logic
func (c *Client) wrapRetryMiddleware(next Handler) Handler {
	return func(req *Request) (*Response, error) {
		if !c.config.Retry.Enabled {
			return next(req)
		}

		var lastErr error
		var lastResp *Response

		for attempt := 0; attempt < c.config.Retry.MaxAttempts; attempt++ {
			req.RetryCount = attempt

			resp, err := next(req)

			if err == nil {
				return resp, nil
			}

			lastResp = resp
			lastErr = err

			if hubspotErr, ok := err.(*HubSpotError); ok {
				if !hubspotErr.IsRetryable {
					return resp, err
				}

				if attempt < c.config.Retry.MaxAttempts-1 {
					backoff := calculateBackoffDuration(attempt, hubspotErr.RetryAfter, c.config.Retry)
					select {
					case <-time.After(backoff):
					case <-req.Context.Done():
						return lastResp, req.Context.Err()
					}
				}
			} else {
				return resp, err
			}
		}

		return lastResp, lastErr
	}
}

// calculateBackoffDuration calculates exponential backoff with jitter
func calculateBackoffDuration(attempt int, retryAfter time.Duration, cfg RetryConfig) time.Duration {
	// If Retry-After header was provided, respect it
	if retryAfter > 0 {
		return retryAfter
	}

	// Calculate exponential backoff: initial * 2^attempt
	backoff := time.Duration(math.Pow(2, float64(attempt))) * cfg.InitialBackoff

	// Add jitter: ±10% randomness
	jitter := time.Duration(rand.Int63n(int64(backoff / 10)))
	if rand.Intn(2) == 0 {
		backoff += jitter
	} else {
		backoff -= jitter
	}

	// Cap at max backoff
	if backoff > cfg.MaxBackoff {
		backoff = cfg.MaxBackoff
	}

	return backoff
}

// httpMiddleware performs the actual HTTP request
func (c *Client) httpMiddleware() Handler {
	return func(req *Request) (*Response, error) {
		// Build full URL
		fullURL := c.config.BaseURL + req.Path

		// Add query parameters
		if len(req.QueryParams) > 0 {
			values := url.Values{}
			for k, v := range req.QueryParams {
				values.Add(k, v)
			}
			fullURL += "?" + values.Encode()
		}

		// Prepare request body
		var bodyReader io.Reader
		if req.Body != nil {
			bodyBytes, err := marshalRequestBody(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(bodyBytes)
			req.AddHeader("Content-Type", "application/json")
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(req.Context, req.Method, fullURL, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Copy headers from request wrapper
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		// Set default headers
		httpReq.Header.Set("User-Agent", "go-hubspot-sdk/1.0")

		c.LogPrintf("Request Headers: %v", httpReq.Header)
		c.LogPrintf("Making request: %s %s", req.Method, fullURL)

		// Perform request
		httpResp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}
		defer httpResp.Body.Close()

		// Read response body
		respBodyBytes, err := readResponseBody(httpResp)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Create response wrapper
		resp := NewResponse(httpResp.StatusCode, respBodyBytes, httpResp.Header)
		resp.RateLimit = ExtractRateLimitInfo(httpResp.Header)

		c.LogPrintf("Rate Limit: %v", resp.RateLimit)

		// Handle error responses
		if httpResp.StatusCode >= 400 {
			resp.HubSpotError = ParseHubSpotError(httpResp.StatusCode, respBodyBytes, httpResp.Header)
			return resp, resp.HubSpotError
		}

		return resp, nil
	}
}

func (c *Client) LogPrintf(format string, v ...any) {
	if c.config.LoggingEnabled {
		c.logger.Printf(format, v...)
	}
}

// marshalRequestBody marshals the request body to JSON bytes
func marshalRequestBody(body any) ([]byte, error) {
	switch v := body.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return jsonMarshal(v)
	}
}

// readResponseBody reads the HTTP response body
func readResponseBody(httpResp *http.Response) ([]byte, error) {
	defer httpResp.Body.Close()

	// Read all body content
	// Use 0 capacity if ContentLength is unknown (-1) or negative
	capacity := httpResp.ContentLength
	if capacity < 0 {
		capacity = 0
	}
	bodyBytes := make([]byte, 0, capacity)

	// Use a temporary buffer to read
	tmpBuffer := make([]byte, 4096)
	for {
		n, err := httpResp.Body.Read(tmpBuffer)
		if n > 0 {
			bodyBytes = append(bodyBytes, tmpBuffer[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
	}

	return bodyBytes, nil
}

// jsonMarshal is a wrapper around json.Marshal for consistency
func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// PrintRateLimit is used to test and verify the rate limiter is being properly updated
func (c *Client) PrintRateLimit(writers ...io.Writer) {
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}
	var b []byte
	b = fmt.Appendf(b, "Daily limit: %d\nDaily remaining: %d\nDaily reset time: %s\n", c.rateLimiter.GetDailyLimit(), c.rateLimiter.GetDailyRemaining(), c.rateLimiter.GetDailyResetTime().String())
	for _, writer := range writers {
		n, err := writer.Write(b)
		if n != len(b) {
			fmt.Printf("Failed to write all bytes to writer: %d != %d\n", n, len(b))
		}
		if err != nil {
			fmt.Printf("Failed to write to writer: %s\n", err.Error())
		}
	}
}
