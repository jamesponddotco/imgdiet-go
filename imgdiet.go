// Package imgdiet offers a simple and fast image processing and compression
// solution by leveraging C's [libvips] image processing library and its Go
// binding, [govips].
//
// [libvips]: https://github.com/libvips/libvips [govips]:
// https://github.com/davidbyttow/govips
package imgdiet

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/davidbyttow/govips/v2/vips"
)

// List of image types supported by this package.
const (
	ImageTypeJPEG string = "JPEG"
	ImageTypePNG  string = "PNG"
	ImageTypeGIF  string = "GIF"
)

// ErrUnsupportedImageFormat is returned when the image format is not supported by this package.
const ErrUnsupportedImageFormat xerrors.Error = "unsupported image format"

// Start initializes the libvips library with the given configuration.
func Start(cfg *Config) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	vips.LoggingSettings(cfg.Logger, cfg.LogLevel)

	vips.Startup(&vips.Config{
		ConcurrencyLevel: cfg.MaxConcurrency,
		MaxCacheSize:     int(cfg.Cache),
		ReportLeaks:      cfg.ReportLeaks,
		CollectStats:     false,
	})
}

// Stop shuts down the libvips library.
func Stop() {
	vips.Shutdown()
}

// DetectImageType takes an image as a byte array input and detects its type based on
// its magic bytes. It returns a string representation of the image type and an
// error if the image type is not supported.
func DetectImageType(image []byte) (string, error) {
	switch http.DetectContentType(image) {
	case "image/jpeg":
		return ImageTypeJPEG, nil
	case "image/png":
		return ImageTypePNG, nil
	case "image/gif":
		return ImageTypeGIF, nil
	default:
		return "", ErrUnsupportedImageFormat
	}
}

// DetectImageSize takes an image as a byte array input and detects the image
// size in bytes.
func DetectImageSize(image []byte) int64 {
	return int64(cap(image))
}
