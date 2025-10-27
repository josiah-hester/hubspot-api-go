package client

import "time"

// Config holds all configuration for the HubSpot API client
type Config struct {
	AccessToken string
	BaseURL     string
	Timeout     time.Duration
	RateLimit   RateLimitConfig
	Retry       RetryConfig
}

// RateLimitConfig configures rate limiting behavior
type RateLimitConfig struct {
	MaxBurst   int
	DailyLimit int
	Enabled    bool
}

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Enabled        bool
}

// Option is a functional option for configuring the Client
type Option func(*Config) error

// NewConfig creates a Config with sensible defaults
func NewConfig() *Config {
	return &Config{
		BaseURL: "https://api.hubapi.com",
		Timeout: 30 * time.Second,
		RateLimit: RateLimitConfig{
			MaxBurst:   100,
			DailyLimit: 250000,
			Enabled:    true,
		},
		Retry: RetryConfig{
			MaxAttempts:    3,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     30 * time.Second,
			Enabled:        true,
		},
	}
}

// WithAccessToken sets the API access token
func WithAccessToken(token string) Option {
	return func(cfg *Config) error {
		cfg.AccessToken = token
		return nil
	}
}

// WithBaseURL sets the API base URL (useful for testing)
func WithBaseURL(url string) Option {
	return func(cfg *Config) error {
		cfg.BaseURL = url
		return nil
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(cfg *Config) error {
		cfg.Timeout = timeout
		return nil
	}
}

// WithRateLimitMaxBurst sets the maximum burst for rate limiting
func WithRateLimitMaxBurst(burst int) Option {
	return func(cfg *Config) error {
		cfg.RateLimit.MaxBurst = burst
		return nil
	}
}

// WithRateLimitDailyLimit sets the daily rate limit
func WithRateLimitDailyLimit(limit int) Option {
	return func(cfg *Config) error {
		cfg.RateLimit.DailyLimit = limit
		return nil
	}
}

// WithRateLimitEnabled enables/disables rate limiting
func WithRateLimitEnabled(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.RateLimit.Enabled = enabled
		return nil
	}
}

// WithRetryMaxAttempts sets the maximum number of retry attempts
func WithRetryMaxAttempts(attempts int) Option {
	return func(cfg *Config) error {
		cfg.Retry.MaxAttempts = attempts
		return nil
	}
}

// WithRetryEnabled enables/disables retries
func WithRetryEnabled(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.Retry.Enabled = enabled
		return nil
	}
}

// WithRetryBackoff sets the initial and max backoff for retries
func WithRetryBackoff(initial, max time.Duration) Option {
	return func(cfg *Config) error {
		cfg.Retry.InitialBackoff = initial
		cfg.Retry.MaxBackoff = max
		return nil
	}
}
