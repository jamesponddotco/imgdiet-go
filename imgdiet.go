package imgdiet

import (
	"git.sr.ht/~jamesponddotco/imgdiet-go/internal/filetype"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/cshum/vipsgen/vips"
)

// ErrUnsupportedFormat is returned when the image format is not supported by
// this package.
const ErrUnsupportedFormat xerrors.Error = "unsupported image format"

// Supported image formats.
const (
	// FormatUnknown represents an unknown or unsupported image format.
	FormatUnknown Format = iota
	FormatJPEG
	FormatPNG
	FormatGIF
	FormatWebP
	FormatAVIF
	FormatHEIF
	FormatTIFF
)

// Format represents an image format supported by the package.
type Format int

// String returns a human-readable string representation of the format.
func (f Format) String() string {
	switch f { //nolint:exhaustive // default would be FormatUnknown
	case FormatJPEG:
		return "JPEG"
	case FormatPNG:
		return "PNG"
	case FormatGIF:
		return "GIF"
	case FormatWebP:
		return "WebP"
	case FormatAVIF:
		return "AVIF"
	case FormatHEIF:
		return "HEIF"
	case FormatTIFF:
		return "TIFF"
	default:
		return "Unknown"
	}
}

// Start initializes the libvips library with the given configuration.
//
// Start must be called before any image processing operations and should be
// called exactly once during application startup. Calling Start multiple times
// may result in undefined behavior.
//
// If cfg is nil, [DefaultConfig] is used.
func Start(cfg *Config) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	vips.Startup(&vips.Config{
		ConcurrencyLevel: cfg.MaxConcurrency,
		MaxCacheSize:     cfg.Cache,
		ReportLeaks:      cfg.ReportLeaks,
	})
}

// Stop shuts down the libvips library and releases all associated resources.
//
// Stop should be called once when the application is finished processing
// images, typically via defer immediately after [Start]. After Stop is called,
// no further image processing operations should be performed.
func Stop() {
	vips.Shutdown()
}

// DetectFormat identifies the image format by examining magic bytes. It returns
// [FormatUnknown] and [ErrUnsupportedFormat] for unrecognized formats.
func DetectFormat(image []byte) (Format, error) {
	format := filetype.Detect(image)

	switch format { //nolint:exhaustive // default would be FormatUnknown
	case filetype.JPEG:
		return FormatJPEG, nil
	case filetype.PNG:
		return FormatPNG, nil
	case filetype.GIF:
		return FormatGIF, nil
	case filetype.WebP:
		return FormatWebP, nil
	case filetype.AVIF:
		return FormatAVIF, nil
	case filetype.HEIF:
		return FormatHEIF, nil
	case filetype.TIFF:
		return FormatTIFF, nil
	default:
		return FormatUnknown, ErrUnsupportedFormat
	}
}
