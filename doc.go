// Package imgdiet provides high-performance image resizing, format conversion,
// and compression optimized for web delivery. The package wraps the [libvips]
// image processing library using [vipsgen].
//
// # Supported Formats
//
// The package supports the following image formats for both input and output:
//
//   - JPEG
//   - PNG
//   - GIF
//   - WebP
//   - AVIF
//   - HEIF
//   - TIFF
//
// Input format is automatically detected from file contents using magic bytes.
//
// # Usage
//
// Initialize libvips before processing any images, and shut it down when done:
//
//	func main() {
//		imgdiet.Start(imgdiet.DefaultConfig())
//		defer imgdiet.Stop()
//
//		// Process images...
//	}
//
// Open an image from any [io.Reader], process it, and export to the desired format:
//
//	img, err := imgdiet.Open(ctx, file)
//	if err != nil {
//		return err
//	}
//
//	opts := imgdiet.DefaultOptions()
//	opts.Width = 800
//	opts.Quality = 75
//
//	output, err := img.Export(ctx, imgdiet.FormatWebP, opts)
//	if err != nil {
//		return err
//	}
//
// Each call to [Image.Export] works from the original image data, allowing
// multiple exports with different formats or settings without quality loss from
// repeated processing:
//
//	webp, _ := img.Export(ctx, imgdiet.FormatWebP, opts)
//	avif, _ := img.Export(ctx, imgdiet.FormatAVIF, opts)
//	jpeg, _ := img.Export(ctx, imgdiet.FormatJPEG, opts)
//
// # Configuration
//
// Use [Config] to control libvips memory usage and concurrency. The
// [DefaultConfig] function returns settings suitable for most server
// environments.
//
// Use [Options] to control image output dimensions, quality, and
// format-specific encoding parameters. The [DefaultOptions] function returns
// settings optimized for web delivery.
//
// # Requirements
//
// This package requires libvips 8.18 or later to be installed on the system.
// See the libvips documentation for installation instructions.
//
// [libvips]: https://github.com/libvips/libvips
// [vipsgen]: https://github.com/cshum/vipsgen
package imgdiet
