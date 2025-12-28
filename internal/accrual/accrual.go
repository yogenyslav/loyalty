// Package accrual provides methods to interact with the accrual system.
package accrual

import (
	"net/http"
)

// Config holds configuration for the AccrualClient.
type Config struct {
	BaseURL      string `yaml:"accrual_addr"    env:"ACCRUAL_SYSTEM_ADDRESS"`
	PollInterval int    `yaml:"poll_accrual"    env:"POLL_ACCRUAL_INTERVAL"  env-default:"10"` // in seconds
	NumWorkers   int    `yaml:"accrual_workers" env:"ACCRUAL_WORKERS"        env-default:"5"`
}

// HTTPClient defines the interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AccrualClient implements client for the accrual system.
type AccrualClient struct {
	cfg    *Config
	client HTTPClient
}

// NewClient creates a new AccrualClient instance.
func NewClient(cfg *Config, client HTTPClient) *AccrualClient {
	return &AccrualClient{
		cfg:    cfg,
		client: client,
	}
}

// PollInterval returns the polling interval for the accrual system.
func (ac *AccrualClient) PollInterval() int {
	return ac.cfg.PollInterval
}
