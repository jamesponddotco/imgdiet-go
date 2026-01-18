package imgdiet

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/cshum/vipsgen/vips"
)

const (
	// ErrOpenImage is returned when an image cannot be opened.
	ErrOpenImage xerrors.Error = "failed to open image"

	// ErrExportImage is returned when an image cannot be exported.
	ErrExportImage xerrors.Error = "failed to export image"

	// ErrNilReader is returned when a nil reader is provided.
	ErrNilReader xerrors.Error = "reader is nil"

	// ErrNilContext is returned when a nil context is provided.
	ErrNilContext xerrors.Error = "context is nil"
)

// Image represents a loaded image and manages its lifecycle.
type Image struct {
	// data holds the original image bytes for re-processing.
	data []byte

	// format is the detected image format.
	format Format

	// size is the size of the original image in bytes.
	size int64
}

// Open reads image data from r and returns an [Image] ready for processing.
//
// The entire contents of r are read into memory to enable multiple exports with
// different settings. For very large images, ensure sufficient memory is
// available.
func Open(ctx context.Context, r io.Reader) (*Image, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, ErrNilContext)
	}

	if r == nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, ErrNilReader)
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, ctx.Err())
	default:
	}

	buffer := bytes.NewBuffer(make([]byte, 0))

	_, err := io.Copy(buffer, r)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, err)
	}

	format, err := DetectFormat(buffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenImage, err)
	}

	var img *Image

	// vipsgen may panic on invalid images, so we recover and convert to error.
	defer func() {
		if r := recover(); r != nil {
			img = nil
			err = fmt.Errorf("%w: %v", ErrOpenImage, r)
		}
	}()

	img = &Image{
		data:   buffer.Bytes(),
		format: format,
		size:   int64(buffer.Len()),
	}

	return img, nil
}

// Size returns the size of the original image in bytes.
func (i *Image) Size() int64 {
	return i.size
}

// Format returns the detected format of the original image.
func (i *Image) Format() Format {
	return i.format
}

// Export encodes the image to the specified format and returns the result.
//
// Export returns three values: the encoded image data, the number of bytes
// saved compared to the original, and any error. The bytes saved value may be
// negative if the output is larger than the original, which can occur when
// converting to less efficient formats or using high quality settings.
//
// Each call to Export decodes from the original image data, preserving full
// quality regardless of how many times the image is exported. This allows
// exporting to multiple formats without compounding compression artifacts:
//
//	webp, savedWebP, _ := img.Export(ctx, imgdiet.FormatWebP, opts)
//	avif, savedAVIF, _ := img.Export(ctx, imgdiet.FormatAVIF, opts)
//
// If opts is nil, [DefaultOptions] is used.
func (i *Image) Export(format Format, opts *Options) ([]byte, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	ref, err := vips.NewImageFromBuffer(i.data, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrExportImage, err)
	}
	defer ref.Close()

	if opts.OptimizeICCProfile {
		if err = ref.IccTransform("srgb", nil); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrExportImage, err)
		}
	}

	if opts.Width > 0 || opts.Height > 0 {
		if err = i.applyResize(ref, opts); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrExportImage, err)
		}
	}

	var output []byte

	switch format { //nolint:exhaustive // default would be FormatUnknown
	case FormatJPEG:
		output, err = i.exportJPEG(ref, opts)
	case FormatPNG:
		output, err = i.exportPNG(ref, opts)
	case FormatGIF:
		output, err = i.exportGIF(ref, opts)
	case FormatWebP:
		output, err = i.exportWebP(ref, opts)
	case FormatAVIF:
		output, err = i.exportAVIF(ref, opts)
	case FormatHEIF:
		output, err = i.exportHEIF(ref, opts)
	case FormatTIFF:
		output, err = i.exportTIFF(ref, opts)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, format.String())
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrExportImage, err)
	}

	return output, nil
}

// applyResize applies resize transformation to the image reference.
func (*Image) applyResize(ref *vips.Image, opts *Options) error {
	var (
		originalWidth  = ref.Width()
		originalHeight = ref.Height()
		aspectRatio    = float64(originalWidth) / float64(originalHeight)
		width          = opts.Width
		height         = opts.Height
	)

	if width == 0 {
		width = int(math.Round(float64(height) * aspectRatio))
	} else if height == 0 {
		height = int(math.Round(float64(width) / aspectRatio))
	}

	// Cap dimensions to original size to prevent upscaling.
	if width > originalWidth {
		width = originalWidth
	}

	if height > originalHeight {
		height = originalHeight
	}

	if err := ref.ThumbnailImage(width, &vips.ThumbnailImageOptions{
		Height: height,
		Crop:   vips.InterestingCentre,
	}); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// exportJPEG exports the image as JPEG format.
func (*Image) exportJPEG(ref *vips.Image, opts *Options) ([]byte, error) {
	options := &vips.JpegsaveBufferOptions{
		Q:                  opts.Quality,
		Interlace:          opts.Interlaced,
		OptimizeCoding:     opts.OptimizeCoding,
		TrellisQuant:       opts.TrellisQuant,
		OvershootDeringing: opts.OvershootDeringing,
		OptimizeScans:      opts.OptimizeScans,
		QuantTable:         opts.QuantTable,
	}

	if opts.StripMetadata {
		options.Keep = vips.KeepNone
	}

	data, err := ref.JpegsaveBuffer(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

// exportPNG exports the image as PNG format.
func (*Image) exportPNG(ref *vips.Image, opts *Options) ([]byte, error) {
	options := &vips.PngsaveBufferOptions{
		Compression: opts.Compression,
		Interlace:   opts.Interlaced,
		Q:           opts.Quality,
		Dither:      opts.Dither,
		Bitdepth:    opts.Bitdepth,
	}

	if opts.StripMetadata {
		options.Keep = vips.KeepNone
	}

	data, err := ref.PngsaveBuffer(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

// exportGIF exports the image as GIF format.
func (*Image) exportGIF(ref *vips.Image, opts *Options) ([]byte, error) {
	options := &vips.GifsaveBufferOptions{
		Dither:   opts.Dither,
		Effort:   opts.Effort,
		Bitdepth: opts.Bitdepth,
	}

	if opts.StripMetadata {
		options.Keep = vips.KeepNone
	}

	data, err := ref.GifsaveBuffer(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

// exportWebP exports the image as WebP format.
func (*Image) exportWebP(ref *vips.Image, opts *Options) ([]byte, error) {
	options := &vips.WebpsaveBufferOptions{
		Q:              opts.Quality,
		Lossless:       opts.Lossless,
		NearLossless:   opts.NearLossless,
		SmartSubsample: opts.SmartSubsample,
		Effort:         opts.Effort,
	}

	if opts.StripMetadata {
		options.Keep = vips.KeepNone
	}

	data, err := ref.WebpsaveBuffer(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

// exportAVIF exports the image as AVIF format.
func (*Image) exportAVIF(ref *vips.Image, opts *Options) ([]byte, error) {
	options := &vips.HeifsaveBufferOptions{
		Q:           opts.Quality,
		Lossless:    opts.Lossless,
		Compression: vips.HeifCompressionAv1,
		Effort:      opts.Effort,
	}

	if opts.StripMetadata {
		options.Keep = vips.KeepNone
	}

	data, err := ref.HeifsaveBuffer(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

// exportHEIF exports the image as HEIF format.
func (*Image) exportHEIF(ref *vips.Image, opts *Options) ([]byte, error) {
	options := &vips.HeifsaveBufferOptions{
		Q:           opts.Quality,
		Lossless:    opts.Lossless,
		Compression: vips.HeifCompressionHevc,
		Effort:      opts.Effort,
	}

	if opts.StripMetadata {
		options.Keep = vips.KeepNone
	}

	data, err := ref.HeifsaveBuffer(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

// exportTIFF exports the image as TIFF format.
func (*Image) exportTIFF(ref *vips.Image, opts *Options) ([]byte, error) {
	options := &vips.TiffsaveBufferOptions{
		Q:          opts.Quality,
		Tile:       opts.TileWidth > 0 || opts.TileHeight > 0,
		TileWidth:  opts.TileWidth,
		TileHeight: opts.TileHeight,
		Pyramid:    opts.Pyramid,
		Bitdepth:   opts.Bitdepth,
	}

	if opts.StripMetadata {
		options.Keep = vips.KeepNone
	}

	data, err := ref.TiffsaveBuffer(options)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}
