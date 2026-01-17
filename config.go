package imgdiet

import (
	"runtime"
)

// Config defines the configuration for the libvips library.
type Config struct {
	// Cache defines the size of the libvips cache in bytes.
	Cache int

	// MaxConcurrency defines the maximum number of concurrent operations that
	// libvips can perform.
	MaxConcurrency int

	// ReportLeaks defines whether libvips should report memory leaks.
	ReportLeaks bool
}

// DefaultConfig returns a set of opinionated and sane defaults for the libvips
// library.
func DefaultConfig() *Config {
	return &Config{
		Cache:          1024 * 1024 * 1024,
		MaxConcurrency: runtime.NumCPU(),
		ReportLeaks:    false,
	}
}
