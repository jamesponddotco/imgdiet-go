package imgdiet

import (
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xlog"
	"github.com/davidbyttow/govips/v2/vips"
)

// List of libvips logging levels for convenience.
const (
	LogLevelError    = vips.LogLevelError
	LogLevelCritical = vips.LogLevelCritical
	LogLevelWarning  = vips.LogLevelWarning
	LogLevelMessage  = vips.LogLevelMessage
	LogLevelInfo     = vips.LogLevelInfo
	LogLevelDebug    = vips.LogLevelDebug
)

// DefaultLogger is a very simple logger the package uses by default. It
// complies with vips.LoggingHandlerFunction.
func DefaultLogger(_ string, verbosity vips.LogLevel, message string) {
	timestamp := time.Now().Format(time.RFC3339)

	levels := map[vips.LogLevel]string{
		LogLevelError:    "error",
		LogLevelCritical: "critical",
		LogLevelWarning:  "warning",
		LogLevelMessage:  "message",
		LogLevelInfo:     "info",
		LogLevelDebug:    "debug",
	}

	level, ok := levels[verbosity]
	if !ok {
		xlog.Printf("Invalid log level: %v. Defaulting to 'info'", verbosity)

		level = "info"
	}

	xlog.Printf("[%s] %s: %s", timestamp, level, message)
}
