package imgdiet_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
)

type errorReader struct{}

func (*errorReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error")
}

func TestDefaultOptions(t *testing.T) {
	t.Parallel()

	params := imgdiet.DefaultOptions()
	if params == nil {
		t.Fatal("expected non-nil parameters")
	}
}

func TestOpen_ErrorReader(t *testing.T) {
	t.Parallel()

	_, err := imgdiet.Open(&errorReader{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		err  error
	}{
		{
			name: "valid_JPEG_image",
			give: _TestDataPath + "/" + _TestValidImageJPG,
			err:  nil,
		},
		{
			name: "valid_PNG_image",
			give: _TestDataPath + "/" + _TestValidImagePNG,
			err:  nil,
		},
		{
			name: "invalid_image",
			give: _TestDataPath + "/" + _TestValidImageGIF,
			err:  imgdiet.ErrUnsupportedImageFormat,
		},
		{
			name: "non-existent_image",
			give: _TestDataPath + "/" + _TestNonExistentImage,
			err:  imgdiet.ErrNilImage,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				file io.Reader
				f    *os.File
				err  error
			)

			if tt.name != "non-existent_image" {
				f, err = os.Open(tt.give)
				if err != nil {
					t.Fatalf("unable to open file: %v", err)
				}
				defer f.Close()

				file = f
			}

			image, err := imgdiet.Open(file)
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, got %v", tt.err, err)
			}
			defer image.Close()

			if err != nil {
				return
			}
		})
	}
}

func TestImage_Optimize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		file    string
		options *imgdiet.Options
		wantErr bool
	}{
		{
			name:    "valid_JPEG_image",
			file:    filepath.Join(_TestDataPath, _TestValidImageJPG),
			options: nil,
			wantErr: false,
		},
		{
			name:    "invalid_JPEG_image",
			file:    filepath.Join(_TestDataPath, _TestInvalidImageJPG),
			options: imgdiet.DefaultOptions(),
			wantErr: true,
		},
		{
			name:    "valid_PNG_image",
			file:    filepath.Join(_TestDataPath, _TestValidImagePNG),
			options: imgdiet.DefaultOptions(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tt.file)
			if err != nil {
				t.Fatalf("unable to open file: %v", err)
			}
			defer file.Close()

			img, err := imgdiet.Open(file)
			if err != nil {
				if tt.wantErr {
					return
				}

				t.Fatalf("Open() failed: %v", err)
			}
			defer img.Close()

			_, err = img.Optimize(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.Optimize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImage_Resize(t *testing.T) { //nolint:gocognit // TODO: simplify
	t.Parallel()

	tests := []struct {
		name           string
		file           string
		width          uint
		height         uint
		options        *imgdiet.Options
		expectedWidth  uint
		expectedHeight uint
		wantErr        bool
	}{
		{
			name:           "valid_JPEG_image",
			file:           filepath.Join(_TestDataPath, _TestValidImageJPG),
			width:          100,
			height:         100,
			options:        imgdiet.DefaultOptions(),
			expectedWidth:  100,
			expectedHeight: 100,
			wantErr:        false,
		},
		{
			name:           "invalid_JPEG_image",
			file:           filepath.Join(_TestDataPath, _TestInvalidImageJPG),
			width:          100,
			height:         100,
			expectedWidth:  100,
			expectedHeight: 100,
			wantErr:        true,
		},
		{
			name:           "valid_PNG_image",
			file:           filepath.Join(_TestDataPath, _TestValidImagePNG),
			width:          100,
			height:         100,
			expectedWidth:  100,
			expectedHeight: 100,
			wantErr:        false,
		},
		{
			name:    "invalid_dimensions",
			file:    filepath.Join(_TestDataPath, _TestValidImageJPG),
			width:   0,
			height:  0,
			wantErr: true,
		},
		{
			name:           "resize_with_only_width",
			file:           filepath.Join(_TestDataPath, _TestValidImageJPG),
			width:          500,
			height:         0,
			expectedWidth:  500,
			expectedHeight: 750,
			wantErr:        false,
		},
		{
			name:           "resize_with_only_height",
			file:           filepath.Join(_TestDataPath, _TestValidImageJPG),
			width:          0,
			height:         500,
			expectedWidth:  333,
			expectedHeight: 500,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tt.file)
			if err != nil {
				t.Fatalf("unable to open file: %v", err)
			}
			defer file.Close()

			img, err := imgdiet.Open(file)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("Open() failed: %v", err)
			}
			defer img.Close()

			originalWidth := img.Width()
			originalHeight := img.Height()

			_, err = img.Resize(tt.width, tt.height, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.Resize() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				if img.Width() != originalWidth || img.Height() != originalHeight {
					t.Errorf("Image.Resize() changed dimensions on error, got width = %d, height = %d, want width = %d, height = %d",
						img.Width(), img.Height(), originalWidth, originalHeight)
				}
			} else {
				if img.Width() != int(tt.expectedWidth) || img.Height() != int(tt.expectedHeight) {
					t.Errorf("Image.Resize() got width = %d, height = %d, want width = %d, height = %d",
						img.Width(), img.Height(), tt.expectedWidth, tt.expectedHeight)
				}
			}
		})
	}
}

func TestImage_SizeAndSaved(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid_JPEG_image",
			file:    filepath.Join(_TestDataPath, _TestValidImageJPG),
			wantErr: false,
		},
		{
			name:    "invalid_JPEG_image",
			file:    filepath.Join(_TestDataPath, _TestInvalidImageJPG),
			wantErr: true,
		},
		{
			name:    "valid_PNG_image",
			file:    filepath.Join(_TestDataPath, _TestValidImagePNG),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tt.file)
			if err != nil {
				t.Fatalf("unable to open file: %v", err)
			}
			defer file.Close()

			img, err := imgdiet.Open(file)
			if err != nil {
				if tt.wantErr {
					return
				}

				t.Fatalf("Open() failed: %v", err)
			}
			defer img.Close()

			_, err = img.Optimize(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.Optimize() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if img.Size() <= img.Saved() {
				t.Errorf("Image.Size() is less than or equal to Image.Saved(), got size = %d, saved = %d",
					img.Size(), img.Saved())
			}
		})
	}
}
