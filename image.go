package imgdiet

import (
	"fmt"
	"io"
	"math"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/davidbyttow/govips/v2/vips"
)

const (
	ErrOpenImage               xerrors.Error = "failed to open image"
	ErrNilImage                xerrors.Error = "image is nil"
	ErrInvalidResizeDimensions xerrors.Error = "dimensions must be greater than 0"
)

// Options represents the parameters used to optimize an image.
type Options struct {
	// Quality defines the quality of the output image. It is a number between 0
	// and 100.
	Quality uint

	// Compression defines the compression level of the output image. It is a
	// number between 0 and 9.
	//
	// Only valid for PNG images.
	Compression uint

	// QuantTable defines the quantization table to be used for the output
	// image. It is a number between 0 and 8.
	//
	// Only valid for JPEG images.
	QuantTable uint

	// OptimizeCoding defines whether the output image should have its coding
	// optimized.
	//
	// Only valid for JPEG images.
	OptimizeCoding bool

	// Interlaced defines whether the output image should be interlaced.
	Interlaced bool

	// StripMetadata defines whether the output image should have its metadata
	// stripped.
	StripMetadata bool

	// OptimizeICCProfile defines whether the output image should have its ICC
	// profile optimized.
	OptimizeICCProfile bool

	// TrellisQuant defines whether the output image should have its
	// quantization tables optimized using trellis quantization.
	//
	// Only valid for JPEG images.
	TrellisQuant bool

	// OvershootDeringing defines whether the output image should have its
	// quantization tables optimized using overshoot deringing.
	//
	// Only valid for JPEG images.
	OvershootDeringing bool

	// OptimizeScans defines whether the output image should have its scans
	// optimized.
	//
	// Only valid for JPEG images.
	OptimizeScans bool
}

// DefaultOptions returns a set of opinionated defaults for optimizing images.
func DefaultOptions() *Options {
	return &Options{
		Quality:            60,
		Compression:        9,
		QuantTable:         3,
		OptimizeCoding:     true,
		Interlaced:         false,
		StripMetadata:      true,
		OptimizeICCProfile: true,
		TrellisQuant:       true,
		OvershootDeringing: true,
		OptimizeScans:      true,
	}
}

// Image defines an image to be optimized and manages its lifecycle.
type Image struct {
	// reference is a govips.ImageRef that contains the image data.
	reference *vips.ImageRef

	// format is a string representation of the image type.
	format string

	// size is the size of the image in bytes.
	size int64

	// saved is the size of the image after optimization in bytes.
	saved int64
}

// Open takes an io.Reader as input for reading and returns an Image instance.
func Open(r io.Reader) (*Image, error) {
	if r == nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, ErrNilImage)
	}

	image, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, err)
	}

	data, err := vips.NewImageFromBuffer(image)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, err)
	}

	imageType, err := DetectImageType(image)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, err)
	}

	return &Image{
		reference: data,
		format:    imageType,
		size:      DetectImageSize(image),
	}, nil
}

// Close releases the resources associated with the Image.
func (i *Image) Close() {
	if i != nil && i.reference != nil {
		i.reference.Close()
	}
}

// Optimize takes the given Options and optimizes the image accordingly. It
// returns the optimized image as a byte slice or an error if the optimization
// fails.
func (i *Image) Optimize(opts *Options) ([]byte, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	if opts.OptimizeICCProfile {
		if err := i.reference.OptimizeICCProfile(); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	var (
		image []byte
		err   error
	)

	switch i.format {
	case ImageTypeJPEG:
		image, err = i.optimizeJPEG(opts)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	case ImageTypePNG:
		image, err = i.optimizePNG(opts)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedImageFormat, i.format)
	}

	i.saved = DetectImageSize(image)

	return image, nil
}

// Resize takes a set of dimensions and resizes the image to those dimensions.
// If opts is not nil, the resulting image is optimized according to the given
// Options.
func (i *Image) Resize(width, height uint, opts *Options) ([]byte, error) {
	if width == 0 && height == 0 {
		return nil, fmt.Errorf("%w", ErrInvalidResizeDimensions)
	}

	var (
		originalWidth  = i.reference.Width()
		originalHeight = i.reference.Height()
		aspectRatio    = float64(originalWidth) / float64(originalHeight)
	)

	if width == 0 {
		width = uint(math.Round(float64(height) * aspectRatio))
	} else if height == 0 {
		height = uint(math.Round(float64(width) / aspectRatio))
	}

	if int(width) > originalWidth {
		width = uint(originalWidth)
	}

	if int(height) > originalHeight {
		height = uint(originalHeight)
	}

	if err := i.reference.Thumbnail(int(width), int(height), vips.InterestingCentre); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if opts != nil {
		return i.Optimize(opts)
	}

	image, _, err := i.reference.ExportNative()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return image, nil
}

// Size returns the size of the image in bytes.
func (i *Image) Size() int64 {
	return i.size
}

// Saved returns the size of the image after optimization in bytes.
func (i *Image) Saved() int64 {
	return i.saved
}

// Width returns the width of the image in pixels.
func (i *Image) Width() int {
	return i.reference.Width()
}

// Height returns the height of the image in pixels.
func (i *Image) Height() int {
	return i.reference.Height()
}

// optimizeJPEG takes the given Options and optimizes the image accordingly. It
// returns the optimized image as a byte slice or an error if the optimization
// fails.
func (i *Image) optimizeJPEG(opts *Options) ([]byte, error) {
	options := &vips.JpegExportParams{
		StripMetadata:      opts.StripMetadata,
		Quality:            int(opts.Quality),
		Interlace:          opts.Interlaced,
		OptimizeCoding:     opts.OptimizeCoding,
		TrellisQuant:       opts.TrellisQuant,
		OvershootDeringing: opts.OvershootDeringing,
		OptimizeScans:      opts.OptimizeScans,
		QuantTable:         int(opts.QuantTable),
	}

	image, _, err := i.reference.ExportJpeg(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return image, nil
}

// optimizePNG takes the given Options and optimizes the image accordingly. It
// returns the optimized image as a byte slice or an error if the optimization
// fails.
func (i *Image) optimizePNG(opts *Options) ([]byte, error) {
	options := &vips.PngExportParams{
		StripMetadata: opts.StripMetadata,
		Compression:   int(opts.Compression),
		Interlace:     opts.Interlaced,
		Quality:       int(opts.Quality),
	}

	image, _, err := i.reference.ExportPng(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return image, nil
}
