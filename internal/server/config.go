package server

import "net"

const defaultBodyLimit = 4 * 1024 * 1024 // 4 MB

// Config holds the server configuration settings.
type Config struct {
	Addr        string `yaml:"addr" env:"SERVER_ADDR"`
	Port        string `yaml:"port" env:"SERVER_PORT"`
	RunAddress  string `yaml:"run_address" env:"RUN_ADDRESS"`
	AccrualAddr string `yaml:"accrual_addr" env:"ACCRUAL_SYSTEM_ADDRESS"`
	BodyLimitMb *int   `yaml:"body_limit_mb" env:"BODY_LIMIT_MB"`
}

// GetAddr returns the server address to run on.
func (c *Config) GetAddr() string {
	if c.RunAddress != "" {
		return c.RunAddress
	}
	return net.JoinHostPort(c.Addr, c.Port)
}

// GetBodyLimit returns the body size limit in bytes.
func (c *Config) GetBodyLimit() int {
	if c.BodyLimitMb != nil {
		return *c.BodyLimitMb * 1024 * 1024
	}
	return defaultBodyLimit
}
