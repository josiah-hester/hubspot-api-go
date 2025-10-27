package client

import (
	"net/http"
	"time"
)

type Response struct {
	// HTTP details
	StatusCode int
	Body       []byte
	Headers    http.Header

	// Rate limit data extracted from headers
	RateLimit RateLimitInfo

	// HubSpot error from HubSpot (if applicable)
	HubSpotError *HubSpotError
}

type RateLimitInfo struct {
	Max             int // Request allowd in window
	Remaining       int // Request remaining in window
	DailyLimit      int
	DailyRemaining  int
	IntervalMs      int // Window in milliseconds (usually 100000)
	WindowResetTime time.Time
	DailyResetTime  time.Time
}

// IsRateLimited returns true if rate limit information indicates we're limited
func (r *Response) IsRateLimited() bool {
	return r.RateLimit.Remaining <= 0 || r.RateLimit.DailyRemaining <= 0
}

// NewResponse creates a new Response wrapper
func NewResponse(statusCode int, body []byte, headers http.Header) *Response {
	return &Response{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
		RateLimit: RateLimitInfo{
			IntervalMs: 10000,
		},
	}
}
