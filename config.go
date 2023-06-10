package imgdiet

import (
	"runtime"

	"github.com/davidbyttow/govips/v2/vips"
)

// Config defines the configuration for the libvips library.
type Config struct {
	// Logger defines the logger to be used by libvips.
	Logger func(domain string, verbosity vips.LogLevel, message string)

	// LogLevel is the level of logging to be used by libvips.
	LogLevel vips.LogLevel

	// Cache defines the size of the libvips cache in bytes.
	Cache uint64

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
		Logger:         DefaultLogger,
		LogLevel:       vips.LogLevelError,
		Cache:          1024 * 1024 * 1024,
		MaxConcurrency: runtime.NumCPU(),
		ReportLeaks:    false,
	}
}
