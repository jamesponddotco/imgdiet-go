package imgdiet

// Options controls the export behavior for images.
//
// Zero values for numeric fields indicate that the encoder should use
// its default value or preserve the original image property.
type Options struct {
	// Width specifies the target width in pixels.
	//
	// If zero, the original width is preserved. If only Width is specified,
	// Height is calculated to maintain aspect ratio.
	Width int

	// Height specifies the target height in pixels.
	//
	// If zero, the original height is preserved. If only Height is specified,
	// Width is calculated to maintain aspect ratio.
	Height int

	// Effort controls the CPU effort spent on encoding, ranging from 0
	// (fastest) to 9 (slowest, best compression). Higher values produce smaller
	// files at the cost of encoding speed.
	//
	// This setting applies to GIF, WebP, and AVIF.
	Effort int

	// Quality specifies the lossy compression quality, ranging from 1 (lowest)
	// to 100 (highest). Higher values produce larger files with better visual
	// fidelity.
	//
	// This setting applies to JPEG, WebP, AVIF, and HEIF.
	Quality int

	// Compression specifies the PNG compression level, ranging from 0 (no
	// compression, fastest) to 9 (maximum compression, slowest). This affects
	// file size and encoding time but not image quality, as PNG is lossless.
	Compression int

	// Speed controls the AVIF encoder speed preset, ranging from 0 (slowest,
	// best compression) to 9 (fastest, largest files). This is similar to
	// Effort, but specific to AVIF encoding.
	Speed int

	// TileWidth specifies the TIFF tile width in pixels.
	//
	// If zero, the encoder uses strip-based encoding instead of tiled encoding.
	// When set, TileHeight should also be specified.
	TileWidth int

	// TileHeight specifies the TIFF tile height in pixels.
	//
	// If zero, the encoder uses strip-based encoding instead of tiled encoding.
	// When set, TileWidth should also be specified.
	TileHeight int

	// Bitdepth specifies the number of bits per channel in the output image.
	// Valid values are 1, 2, 4, 8, and 16. If zero, the encoder preserves the
	// original bit depth or uses a format-appropriate default.
	Bitdepth int

	// QuantTable selects the JPEG quantization table, ranging from 0 to 8.
	QuantTable int

	// Dither specifies the amount of Floyd-Steinberg dithering to apply when
	// reducing color depth for palette-based formats. Values range from 0.0 (no
	// dithering) to 1.0 (full dithering).
	//
	// Dithering reduces banding artifacts but may introduce noise.
	Dither float64

	// Interlaced enables progressive encoding for JPEG and PNG, or interlaced
	// encoding for GIF. Progressive images render incrementally, displaying a
	// low-quality preview before the full image loads.
	Interlaced bool

	// StripMetadata removes EXIF, XMP, IPTC, and ICC profile data from the
	// output image.
	//
	// This reduces file size but discards information such as camera settings,
	// GPS coordinates, and color profiles.
	StripMetadata bool

	// OptimizeICCProfile converts embedded ICC color profiles to sRGB and
	// removes the profile from the output. This reduces file size while
	// maintaining color accuracy for displays using the sRGB color space.
	OptimizeICCProfile bool

	// OptimizeCoding enables Huffman table optimization for JPEG encoding. This
	// produces smaller files with a modest increase in encoding time.
	OptimizeCoding bool

	// TrellisQuant enables trellis quantization for JPEG encoding. This
	// rate-distortion optimization technique produces smaller files at
	// equivalent quality levels, at the cost of slower encoding.
	TrellisQuant bool

	// OvershootDeringing enables overshoot deringing for JPEG encoding. This
	// reduces ringing artifacts near sharp edges in the image.
	OvershootDeringing bool

	// OptimizeScans enables progressive scan optimization for JPEG encoding.
	// This reorders the spectral data in progressive JPEGs for improved
	// compression.
	//
	// Only effective when Interlaced is true.
	OptimizeScans bool

	// Lossless enables mathematically lossless encoding for WebP and AVIF.
	//
	// The output image will be pixel-identical to the input, but file sizes
	// will be significantly larger than lossy encoding.
	Lossless bool

	// NearLossless enables near-lossless encoding for WebP. This preprocessing
	// step slightly modifies the image to improve compression while minimizing
	// perceptible quality loss.
	//
	// Only effective when Lossless is true.
	NearLossless bool

	// SmartSubsample enables context-dependent chroma subsampling for WebP.
	// This adaptively applies subsampling based on image content, preserving
	// detail in areas with sharp color transitions.
	SmartSubsample bool

	// Pyramid enables multi-resolution pyramid generation for TIFF output.
	//
	// Pyramid TIFFs store the image at multiple resolutions, enabling efficient
	// display at different zoom levels.
	Pyramid bool
}

// DefaultOptions returns an [Options] instance configured for web-optimized
// image export.
//
// The defaults prioritize smaller file sizes over maximum quality, making them
// suitable for web delivery where bandwidth and load times are important.
// Metadata is stripped for privacy and reduced file size, and ICC profiles are
// converted to sRGB for broad display compatibility.
//
// Image dimensions are preserved by default.
func DefaultOptions() *Options {
	return &Options{
		Quality:            60,
		Compression:        9,
		Effort:             7,
		QuantTable:         3,
		Bitdepth:           8,
		OptimizeCoding:     true,
		Interlaced:         false,
		StripMetadata:      true,
		OptimizeICCProfile: true,
		TrellisQuant:       true,
		OvershootDeringing: true,
		OptimizeScans:      true,
		Lossless:           false,
		NearLossless:       false,
		SmartSubsample:     false,
		Speed:              6,
	}
}
