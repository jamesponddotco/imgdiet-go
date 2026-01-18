// Package filetype provides image format detection based on magic bytes.
//
// This package is built by Claude Opus based on [h2non/filetype], which is MIT
// licensed. Copyright goes to Tomas Aparicio.
//
// [h2non/filetype]: https://github.com/h2non/filetype
package filetype

import "encoding/binary"

// Format represents an image format.
type Format int

// Supported image formats.
const (
	Unknown Format = iota
	JPEG
	PNG
	GIF
	WebP
	AVIF
	HEIF
	TIFF
)

// String returns a human-readable string representation of the format.
func (f Format) String() string {
	switch f { //nolint:exhaustive // default would be Unknown
	case JPEG:
		return "JPEG"
	case PNG:
		return "PNG"
	case GIF:
		return "GIF"
	case WebP:
		return "WebP"
	case AVIF:
		return "AVIF"
	case HEIF:
		return "HEIF"
	case TIFF:
		return "TIFF"
	default:
		return "Unknown"
	}
}

// Detect identifies the image format from the given byte slice by examining
// magic bytes. Returns Unknown if the format cannot be determined.
func Detect(buf []byte) Format {
	switch {
	case isJPEG(buf):
		return JPEG
	case isPNG(buf):
		return PNG
	case isGIF(buf):
		return GIF
	case isWebP(buf):
		return WebP
	case isTIFF(buf):
		return TIFF
	case isAVIF(buf):
		return AVIF
	case isHEIF(buf):
		return HEIF
	default:
		return Unknown
	}
}

// isJPEG checks for JPEG magic bytes (FF D8 FF).
func isJPEG(buf []byte) bool {
	return len(buf) > 2 &&
		buf[0] == 0xFF &&
		buf[1] == 0xD8 &&
		buf[2] == 0xFF
}

// isPNG checks for PNG magic bytes (89 50 4E 47).
func isPNG(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x89 &&
		buf[1] == 0x50 &&
		buf[2] == 0x4E &&
		buf[3] == 0x47
}

// isGIF checks for GIF magic bytes (47 49 46).
func isGIF(buf []byte) bool {
	return len(buf) > 2 &&
		buf[0] == 0x47 &&
		buf[1] == 0x49 &&
		buf[2] == 0x46
}

// isWebP checks for WebP format by verifying RIFF container and WEBP signature.
// WebP files have RIFF header followed by "WEBP" at bytes 8-11.
func isWebP(buf []byte) bool {
	return len(buf) > 11 &&
		buf[8] == 0x57 &&
		buf[9] == 0x45 &&
		buf[10] == 0x42 &&
		buf[11] == 0x50
}

// isTIFF checks for TIFF magic bytes in both little-endian (49 49 2A 00) and
// big-endian (4D 4D 00 2A) formats.
func isTIFF(buf []byte) bool {
	if len(buf) < 4 {
		return false
	}

	var (
		bigEndian    = buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x00 && buf[3] == 0x2A
		littleEndian = buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2A && buf[3] == 0x00
	)

	return littleEndian || bigEndian
}

// isAVIF checks for AVIF format by examining the ISOBMFF ftyp box.
func isAVIF(buf []byte) bool {
	if !isISOBMFF(buf) {
		return false
	}

	majorBrand, compatibleBrands := getFtyp(buf)

	if majorBrand == "avif" || majorBrand == "avis" {
		return true
	}

	if majorBrand == "mif1" || majorBrand == "msf1" {
		for _, brand := range compatibleBrands {
			if brand == "avif" || brand == "avis" {
				return true
			}
		}
	}

	return false
}

// isHEIF checks for HEIF format by examining the ISOBMFF ftyp box.
func isHEIF(buf []byte) bool {
	if !isISOBMFF(buf) {
		return false
	}

	majorBrand, compatibleBrands := getFtyp(buf)

	if majorBrand == "heic" || majorBrand == "heix" || majorBrand == "hevc" || majorBrand == "hevx" {
		return true
	}

	if majorBrand == "mif1" || majorBrand == "msf1" {
		for _, brand := range compatibleBrands {
			switch brand {
			case "heic", "heix", "hevc", "hevx":
				return true
			}
		}
	}

	return false
}

// isISOBMFF checks whether the buffer represents ISO Base Media File Format
// data. It validates the ftyp box signature and ensures the buffer contains the
// complete ftyp box.
func isISOBMFF(buf []byte) bool {
	if len(buf) < 16 {
		return false
	}

	if buf[4] != 'f' || buf[5] != 't' || buf[6] != 'y' || buf[7] != 'p' {
		return false
	}

	ftypLength := binary.BigEndian.Uint32(buf[0:4])

	return len(buf) >= int(ftypLength)
}

// getFtyp extracts the major brand and compatible brands from an ISOBMFF ftyp box.
// The caller must ensure isISOBMFF returns true before calling this function.
func getFtyp(buf []byte) (majorBrand string, compatibleBrands []string) { //nolint:nonamedreturns // the code is clearer with these
	if len(buf) < 12 {
		return "", nil
	}

	ftypLength := binary.BigEndian.Uint32(buf[0:4])
	majorBrand = string(buf[8:12])

	for i := 16; i < int(ftypLength); i += 4 {
		if len(buf) >= i+4 {
			compatibleBrands = append(compatibleBrands, string(buf[i:i+4]))
		}
	}

	return majorBrand, compatibleBrands
}
