package filetype_test

import (
	"os"
	"testing"

	"git.sr.ht/~jamesponddotco/imgdiet-go/internal/filetype"
)

const testDataPath = "../../testdata"

func TestDetect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		file string
		want filetype.Format
	}{
		{
			name: "jpeg",
			file: testDataPath + "/james-pond-hotel-chair.jpg",
			want: filetype.JPEG,
		},
		{
			name: "png",
			file: testDataPath + "/cipherhost-avatar.png",
			want: filetype.PNG,
		},
		{
			name: "gif",
			file: testDataPath + "/whoops.gif",
			want: filetype.GIF,
		},
		{
			name: "webp",
			file: testDataPath + "/webp-animated.webp",
			want: filetype.WebP,
		},
		{
			name: "avif",
			file: testDataPath + "/avif-8bit.avif",
			want: filetype.AVIF,
		},
		{
			name: "heif",
			file: testDataPath + "/heic-24bit.heic",
			want: filetype.HEIF,
		},
		{
			name: "tiff",
			file: testDataPath + "/tif-16bit.tif",
			want: filetype.TIFF,
		},
		{
			name: "invalid",
			file: testDataPath + "/invalid-image.jpg",
			want: filetype.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			got := filetype.Detect(data)
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestDetect_MagicBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data []byte
		want filetype.Format
	}{
		{
			name: "jpeg magic bytes",
			data: []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10},
			want: filetype.JPEG,
		},
		{
			name: "png magic bytes",
			data: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			want: filetype.PNG,
		},
		{
			name: "gif87a magic bytes",
			data: []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61},
			want: filetype.GIF,
		},
		{
			name: "gif89a magic bytes",
			data: []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61},
			want: filetype.GIF,
		},
		{
			name: "webp magic bytes",
			data: []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50},
			want: filetype.WebP,
		},
		{
			name: "tiff little-endian magic bytes",
			data: []byte{0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00},
			want: filetype.TIFF,
		},
		{
			name: "tiff big-endian magic bytes",
			data: []byte{0x4D, 0x4D, 0x00, 0x2A, 0x00, 0x00, 0x00, 0x08},
			want: filetype.TIFF,
		},
		{
			name: "empty buffer",
			data: []byte{},
			want: filetype.Unknown,
		},
		{
			name: "too short buffer",
			data: []byte{0xFF, 0xD8},
			want: filetype.Unknown,
		},
		{
			name: "random bytes",
			data: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			want: filetype.Unknown,
		},
		{
			name: "riff without webp",
			data: []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x41, 0x56, 0x49, 0x20},
			want: filetype.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := filetype.Detect(tt.data)
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestDetect_ISOBMFF(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data []byte
		want filetype.Format
	}{
		{
			name: "avif major brand",
			data: buildFtypBox("avif", nil),
			want: filetype.AVIF,
		},
		{
			name: "avis major brand",
			data: buildFtypBox("avis", nil),
			want: filetype.AVIF,
		},
		{
			name: "heic major brand",
			data: buildFtypBox("heic", nil),
			want: filetype.HEIF,
		},
		{
			name: "heix major brand",
			data: buildFtypBox("heix", nil),
			want: filetype.HEIF,
		},
		{
			name: "mif1 with avif compatible brand",
			data: buildFtypBox("mif1", []string{"miaf", "avif"}),
			want: filetype.AVIF,
		},
		{
			name: "mif1 with heic compatible brand",
			data: buildFtypBox("mif1", []string{"miaf", "heic"}),
			want: filetype.HEIF,
		},
		{
			name: "msf1 with avif compatible brand",
			data: buildFtypBox("msf1", []string{"msf1", "avif"}),
			want: filetype.AVIF,
		},
		{
			name: "msf1 with heic compatible brand",
			data: buildFtypBox("msf1", []string{"iso8", "heic"}),
			want: filetype.HEIF,
		},
		{
			name: "mif1 without avif or heic",
			data: buildFtypBox("mif1", []string{"miaf", "MiPr"}),
			want: filetype.Unknown,
		},
		{
			name: "incomplete ftyp box",
			data: []byte{0x00, 0x00, 0x00, 0x14, 'f', 't', 'y', 'p'},
			want: filetype.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := filetype.Detect(tt.data)
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestFormat_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give filetype.Format
		want string
	}{
		{give: filetype.JPEG, want: "JPEG"},
		{give: filetype.PNG, want: "PNG"},
		{give: filetype.GIF, want: "GIF"},
		{give: filetype.WebP, want: "WebP"},
		{give: filetype.AVIF, want: "AVIF"},
		{give: filetype.HEIF, want: "HEIF"},
		{give: filetype.TIFF, want: "TIFF"},
		{give: filetype.Unknown, want: "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			got := tt.give.String()
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

// buildFtypBox creates a minimal ISOBMFF ftyp box for testing.
func buildFtypBox(majorBrand string, compatibleBrands []string) []byte {
	// ftyp box structure:
	// - 4 bytes: box size (big-endian)
	// - 4 bytes: "ftyp"
	// - 4 bytes: major brand
	// - 4 bytes: minor version
	// - n*4 bytes: compatible brands

	boxSize := 16 + len(compatibleBrands)*4
	buf := make([]byte, boxSize)

	// Box size (big-endian)
	buf[0] = byte(boxSize >> 24)
	buf[1] = byte(boxSize >> 16)
	buf[2] = byte(boxSize >> 8)
	buf[3] = byte(boxSize)

	// "ftyp"
	copy(buf[4:8], "ftyp")

	// Major brand
	copy(buf[8:12], majorBrand)

	// Minor version (zeros)
	// buf[12:16] is already zero

	// Compatible brands
	for i, brand := range compatibleBrands {
		copy(buf[16+i*4:20+i*4], brand)
	}

	return buf
}
