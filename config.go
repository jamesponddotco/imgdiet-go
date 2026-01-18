package imgdiet

import (
	"runtime"
)

// Config defines the runtime configuration for the libvips image processing
// library, the underlying library that powers imgdiet.
type Config struct {
	// Cache specifies the maximum size of the libvips operation cache in bytes.
	// The operation cache stores intermediate results to accelerate repeated
	// processing of similar images.
	//
	// Larger values improve performance for repetitive workloads but increase
	// memory consumption. If zero, caching is disabled entirely.
	Cache int

	// MaxConcurrency specifies the maximum number of threads libvips uses for
	// image processing operations. Each operation may use multiple threads for
	// parallel processing of image regions.
	//
	// Higher values improve throughput on multi-core systems but increase
	// memory usage, as each thread maintains its own working buffers. If zero,
	// libvips uses a single thread.
	MaxConcurrency int

	// ReportLeaks enables memory leak detection in libvips. When enabled,
	// libvips logs any unreleased allocations when the library shuts down.
	//
	// This is intended for debugging and development use only, as leak
	// detection adds overhead to memory operations.
	ReportLeaks bool
}

// DefaultConfig returns a [Config] instance suitable for server environments
// with moderate memory availability.
//
// The defaults balance throughput and memory usage:
//   - Cache: 1 GB operation cache.
//   - MaxConcurrency: one thread per CPU core.
//   - ReportLeaks: disabled.
func DefaultConfig() *Config {
	return &Config{
		Cache:          1024 * 1024 * 1024,
		MaxConcurrency: runtime.NumCPU(),
		ReportLeaks:    false,
	}
}
