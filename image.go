package imgdiet

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/davidbyttow/govips/v2/vips"
)

const (
	ErrOpenImage               xerrors.Error = "failed to open image"
	ErrInvalidResizeDimensions xerrors.Error = "dimensions must be greater than 0"
)

// Options represents the parameters used to optimize an image.
type Options struct {
	// Quality defines the quality of the output image. It is a number between 0
	// and 100.
	Quality int

	// Compression defines the compression level of the output image. It is a
	// number between 0 and 9.
	//
	// Only valid for PNG images.
	Compression int

	// QuantTable defines the quantization table to be used for the output
	// image. It is a number between 0 and 8.
	//
	// Only valid for JPEG images.
	QuantTable int

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
		Compression:        6,
		QuantTable:         3,
		OptimizeCoding:     true,
		Interlaced:         true,
		StripMetadata:      true,
		OptimizeICCProfile: true,
		TrellisQuant:       true,
		OvershootDeringing: true,
		OptimizeScans:      true,
	}
}

// Image defines an image to be optimized and manages its lifecycle.
type Image struct {
	// Reference is a govips.ImageRef that contains the image data.
	Reference *vips.ImageRef

	// Options defines the parameters for image optimization.
	Options *Options

	// Type is a string representation of the image type.
	Type string

	// Size is the size of the image in bytes.
	Size int64

	// Saved is the size of the image after optimization in bytes.
	Saved int64
}

// Open takes a named image file as input for reading and returns an Image
// instance.
func Open(name string) (*Image, error) {
	image, err := os.ReadFile(name)
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
		Reference: data,
		Type:      imageType,
		Size:      DetectImageSize(image),
	}, nil
}

// Close releases the resources associated with the Image.
func (i *Image) Close() {
	i.Reference.Close()
}

// Optimize takes the given Options and optimizes the image accordingly. It
// returns the optimized image as a byte slice or an error if the optimization
// fails.
func (i *Image) Optimize(opts *Options) ([]byte, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	i.Options = opts

	if i.Options.OptimizeICCProfile {
		if err := i.Reference.OptimizeICCProfile(); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	var (
		image []byte
		err   error
	)

	switch i.Type {
	case ImageTypeJPEG:
		image, err = i.optimizeJPEG()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	case ImageTypePNG:
		image, err = i.optimizePNG()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedImageFormat, i.Type)
	}

	i.Saved = DetectImageSize(image)

	return image, nil
}

// Resize takes a set of dimensions and resizes the image to those dimensions.
// The resulting image is optimized according to the given Options.
func (i *Image) Resize(width, height uint, opts *Options) ([]byte, error) {
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("%w", ErrInvalidResizeDimensions)
	}

	if opts == nil {
		opts = DefaultOptions()
	}

	i.Options = opts

	if err := i.Reference.Thumbnail(int(width), int(height), vips.InterestingCentre); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return i.Optimize(i.Options)
}

func (i *Image) optimizeJPEG() ([]byte, error) {
	options := &vips.JpegExportParams{
		StripMetadata:      i.Options.StripMetadata,
		Quality:            i.Options.Quality,
		Interlace:          i.Options.Interlaced,
		OptimizeCoding:     i.Options.OptimizeCoding,
		TrellisQuant:       i.Options.TrellisQuant,
		OvershootDeringing: i.Options.OvershootDeringing,
		OptimizeScans:      i.Options.OptimizeScans,
		QuantTable:         i.Options.QuantTable,
	}

	jpg, _, err := i.Reference.ExportJpeg(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return jpg, nil
}

func (i *Image) optimizePNG() ([]byte, error) {
	options := &vips.PngExportParams{
		StripMetadata: i.Options.StripMetadata,
		Compression:   i.Options.Compression,
		Interlace:     i.Options.Interlaced,
		Quality:       i.Options.Quality,
	}

	png, _, err := i.Reference.ExportPng(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return png, nil
}
