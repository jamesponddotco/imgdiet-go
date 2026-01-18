package imgdiet_test

import (
	"os"
	"testing"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
)

const (
	_TestDataPath         string = "testdata"
	_TestValidImageJPG    string = "james-pond-hotel-chair.jpg"
	_TestInvalidImageJPG  string = "invalid-image.jpg"
	_TestValidImagePNG    string = "cipherhost-avatar.png"
	_TestValidImageGIF    string = "whoops.gif"
	_TestValidImageWebP   string = "webp-animated.webp"
	_TestValidImageAVIF   string = "avif-8bit.avif"
	_TestValidImageHEIF   string = "heic-24bit.heic"
	_TestValidImageTIFF   string = "tif-16bit.tif"
	_TestNonExistentImage string = "impossible-girl.jpg"
)

func TestMain(m *testing.M) {
	imgdiet.Start(nil)

	defer imgdiet.Stop()

	m.Run()
}

func TestDetectFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    string
		want    imgdiet.Format
		wantErr bool
	}{
		{
			name:    "jpeg",
			give:    _TestDataPath + "/" + _TestValidImageJPG,
			want:    imgdiet.FormatJPEG,
			wantErr: false,
		},
		{
			name:    "png",
			give:    _TestDataPath + "/" + _TestValidImagePNG,
			want:    imgdiet.FormatPNG,
			wantErr: false,
		},
		{
			name:    "gif",
			give:    _TestDataPath + "/" + _TestValidImageGIF,
			want:    imgdiet.FormatGIF,
			wantErr: false,
		},
		{
			name:    "webp",
			give:    _TestDataPath + "/" + _TestValidImageWebP,
			want:    imgdiet.FormatWebP,
			wantErr: false,
		},
		{
			name:    "avif",
			give:    _TestDataPath + "/" + _TestValidImageAVIF,
			want:    imgdiet.FormatAVIF,
			wantErr: false,
		},
		{
			name:    "heif",
			give:    _TestDataPath + "/" + _TestValidImageHEIF,
			want:    imgdiet.FormatHEIF,
			wantErr: false,
		},
		{
			name:    "tiff",
			give:    _TestDataPath + "/" + _TestValidImageTIFF,
			want:    imgdiet.FormatTIFF,
			wantErr: false,
		},
		{
			name:    "invalid",
			give:    _TestDataPath + "/" + _TestInvalidImageJPG,
			want:    imgdiet.FormatUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.ReadFile(tt.give)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			got, err := imgdiet.DetectFormat(file)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error = %v, got error = %v", tt.wantErr, err)
			}

			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want.String(), got.String())
			}
		})
	}
}

func TestFormat_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want string
		give imgdiet.Format
	}{
		{give: imgdiet.FormatJPEG, want: "JPEG"},
		{give: imgdiet.FormatPNG, want: "PNG"},
		{give: imgdiet.FormatGIF, want: "GIF"},
		{give: imgdiet.FormatWebP, want: "WebP"},
		{give: imgdiet.FormatAVIF, want: "AVIF"},
		{give: imgdiet.FormatHEIF, want: "HEIF"},
		{give: imgdiet.FormatTIFF, want: "TIFF"},
		{give: imgdiet.FormatUnknown, want: "Unknown"},
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
