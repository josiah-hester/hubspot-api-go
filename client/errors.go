package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type HubSpotError struct {
	Status        int
	Message       string
	ErrorType     string // "RATE_LIMIT", "VALIDATION_ERROR", etc.
	Category      string
	PolicyName    string // "DAILY" or "TEN_SECONDLY_ROLLING"
	CorrelationId string

	// Used to deteremine if request should retry
	IsRetryable bool
	RetryAfter  time.Duration
	RawBody     string
}

// Functions to:
// - Parse error from response
// - Determine if retryable (429s with TEN_SECONDLY are retryable, DAILYs probably aren't)
// - Extract Retry-After header

// Error inplements the error interface
func (e *HubSpotError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("HubSpot API error: %s (type: %s, status: %d)", e.Message, e.ErrorType, e.Status)
	}
	return fmt.Sprintf("HubSpot API error: status %d", e.Status)
}

// ParseHubSpotError parses a response into a HubSpotError
func ParseHubSpotError(statusCode int, body []byte, headers http.Header) *HubSpotError {
	err := &HubSpotError{
		Status:  statusCode,
		RawBody: string(body),
	}

	// Extract Retry-After header if present
	if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
		if seconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
			err.RetryAfter = time.Duration(seconds) * time.Second
		} else {
			if t, parseErr := time.Parse(time.RFC1123, retryAfter); parseErr == nil {
				err.RetryAfter = time.Until(t)
			}
		}
	}

	// Try to unmarshal HubSpot error format
	var hubspotResp struct {
		Status        string `json:"status"`
		Message       string `json:"message"`
		ErrorType     string `json:"errorType"`
		Category      string `json:"category"`
		PolicyName    string `json:"policyName"`
		CorrelationId string `json:"correlationId"`
	}

	if unmarshalErr := json.Unmarshal(body, &hubspotResp); unmarshalErr == nil {
		err.Message = hubspotResp.Message
		err.ErrorType = hubspotResp.ErrorType
		err.Category = hubspotResp.Category
		err.PolicyName = hubspotResp.PolicyName
		err.CorrelationId = hubspotResp.CorrelationId
	}

	// Determine if retryable
	err.IsRetryable = isRetryableStatus(statusCode, err.PolicyName)

	return err
}

// isRetryableStatus determines if an HTTP status should be retried
func isRetryableStatus(statusCode int, policyName string) bool {
	switch statusCode {
	case 429:
		return policyName != "DAILY"
	case 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

// ExtractRateLimitInfo extracts rate limit information from response headers
func ExtractRateLimitInfo(headers http.Header) RateLimitInfo {
	info := RateLimitInfo{
		IntervalMs: 10000,
	}

	if maxStr := headers.Get("X-HubSpot-RateLimit-Max"); maxStr != "" {
		if max, err := strconv.Atoi(maxStr); err == nil {
			info.Max = max
		}
	}

	if remainingStr := headers.Get("X-HubSpot-RateLimit-Remaining"); remainingStr != "" {
		if remaining, err := strconv.Atoi(remainingStr); err == nil {
			info.Remaining = remaining
		}
	}

	if intervalStr := headers.Get("X-HubSpot-RateLimit-Interval-Milliseconds"); intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil {
			info.IntervalMs = interval
		}
	}

	if dailyStr := headers.Get("X-HubSpot-RateLimit-Daily"); dailyStr != "" {
		if daily, err := strconv.Atoi(dailyStr); err == nil {
			info.DailyLimit = daily
		}
	}

	if dailyRemStr := headers.Get("X-HubSpot-RateLimit-Daily-Remaining"); dailyRemStr != "" {
		if dailyRem, err := strconv.Atoi(dailyRemStr); err == nil {
			info.DailyRemaining = dailyRem
		}
	}

	return info
}
