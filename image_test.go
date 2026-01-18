package imgdiet_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
)

type errorReader struct{}

func (*errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("mock error")
}

func TestDefaultOptions(t *testing.T) {
	t.Parallel()

	params := imgdiet.DefaultOptions()
	if params == nil {
		t.Fatal("expected non-nil parameters")
	}
}

func TestOpen_NilContext(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join(_TestDataPath, _TestValidImageJPG))
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}
	defer file.Close()

	_, err = imgdiet.Open(nil, file) //nolint:staticcheck // testing nil context behavior
	if !errors.Is(err, imgdiet.ErrNilContext) {
		t.Fatalf("expected ErrNilContext, got: %v", err)
	}
}

func TestOpen_NilReader(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	_, err := imgdiet.Open(ctx, nil)
	if !errors.Is(err, imgdiet.ErrNilReader) {
		t.Fatalf("expected ErrNilReader, got: %v", err)
	}
}

func TestOpen_ErrorReader(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	_, err := imgdiet.Open(ctx, &errorReader{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOpen_CanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	file, err := os.Open(filepath.Join(_TestDataPath, _TestValidImageJPG))
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}
	defer file.Close()

	_, err = imgdiet.Open(ctx, file)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
}

func TestOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want error
	}{
		{
			name: "valid_JPEG_image",
			give: _TestDataPath + "/" + _TestValidImageJPG,
			want: nil,
		},
		{
			name: "valid_PNG_image",
			give: _TestDataPath + "/" + _TestValidImagePNG,
			want: nil,
		},
		{
			name: "valid_GIF_image",
			give: _TestDataPath + "/" + _TestValidImageGIF,
			want: nil,
		},
		{
			name: "valid_WebP_image",
			give: _TestDataPath + "/" + _TestValidImageWebP,
			want: nil,
		},
		{
			name: "valid_AVIF_image",
			give: _TestDataPath + "/" + _TestValidImageAVIF,
			want: nil,
		},
		{
			name: "valid_HEIF_image",
			give: _TestDataPath + "/" + _TestValidImageHEIF,
			want: nil,
		},
		{
			name: "valid_TIFF_image",
			give: _TestDataPath + "/" + _TestValidImageTIFF,
			want: nil,
		},
		{
			name: "non-existent_image",
			give: _TestDataPath + "/" + _TestNonExistentImage,
			want: imgdiet.ErrNilReader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				file io.Reader
				f    *os.File
				err  error
				ctx  = t.Context()
			)

			if tt.name != "non-existent_image" {
				f, err = os.Open(tt.give)
				if err != nil {
					t.Fatalf("unable to open file: %v", err)
				}
				defer f.Close()

				file = f
			}

			_, err = imgdiet.Open(ctx, file)
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected error %v, got %v", tt.want, err)
			}

			if err != nil {
				return
			}
		})
	}
}

func TestImage_Format(t *testing.T) {
	t.Parallel()

	var (
		ctx   = t.Context()
		tests = []struct {
			name string
			give string
			want imgdiet.Format
		}{
			{
				name: "jpeg",
				give: filepath.Join(_TestDataPath, _TestValidImageJPG),
				want: imgdiet.FormatJPEG,
			},
			{
				name: "png",
				give: filepath.Join(_TestDataPath, _TestValidImagePNG),
				want: imgdiet.FormatPNG,
			},
			{
				name: "gif",
				give: filepath.Join(_TestDataPath, _TestValidImageGIF),
				want: imgdiet.FormatGIF,
			},
			{
				name: "webp",
				give: filepath.Join(_TestDataPath, _TestValidImageWebP),
				want: imgdiet.FormatWebP,
			},
			{
				name: "avif",
				give: filepath.Join(_TestDataPath, _TestValidImageAVIF),
				want: imgdiet.FormatAVIF,
			},
			{
				name: "heif",
				give: filepath.Join(_TestDataPath, _TestValidImageHEIF),
				want: imgdiet.FormatHEIF,
			},
			{
				name: "tiff",
				give: filepath.Join(_TestDataPath, _TestValidImageTIFF),
				want: imgdiet.FormatTIFF,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tt.give)
			if err != nil {
				t.Fatalf("unable to open file: %v", err)
			}
			defer file.Close()

			img, err := imgdiet.Open(ctx, file)
			if err != nil {
				t.Fatalf("Open() failed: %v", err)
			}

			if img.Format() != tt.want {
				t.Fatalf("expected format %s, got %s", tt.want.String(), img.Format().String())
			}
		})
	}
}

func TestImage_Export(t *testing.T) {
	t.Parallel()

	var (
		ctx   = t.Context()
		tests = []struct {
			name        string
			give        string
			giveOptions *imgdiet.Options
			giveFormat  imgdiet.Format
			want        bool
		}{
			{
				name:        "export_JPEG_as_JPEG",
				give:        filepath.Join(_TestDataPath, _TestValidImageJPG),
				giveFormat:  imgdiet.FormatJPEG,
				giveOptions: nil,
				want:        false,
			},
			{
				name:        "export_JPEG_as_PNG",
				give:        filepath.Join(_TestDataPath, _TestValidImageJPG),
				giveFormat:  imgdiet.FormatPNG,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_PNG_as_JPEG",
				give:        filepath.Join(_TestDataPath, _TestValidImagePNG),
				giveFormat:  imgdiet.FormatJPEG,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_JPEG_as_WebP",
				give:        filepath.Join(_TestDataPath, _TestValidImageJPG),
				giveFormat:  imgdiet.FormatWebP,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_GIF_as_GIF",
				give:        filepath.Join(_TestDataPath, _TestValidImageGIF),
				giveFormat:  imgdiet.FormatGIF,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_AVIF_as_AVIF",
				give:        filepath.Join(_TestDataPath, _TestValidImageAVIF),
				giveFormat:  imgdiet.FormatAVIF,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_HEIF_as_HEIF",
				give:        filepath.Join(_TestDataPath, _TestValidImageHEIF),
				giveFormat:  imgdiet.FormatHEIF,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_TIFF_as_TIFF",
				give:        filepath.Join(_TestDataPath, _TestValidImageTIFF),
				giveFormat:  imgdiet.FormatTIFF,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_JPEG_as_AVIF",
				give:        filepath.Join(_TestDataPath, _TestValidImageJPG),
				giveFormat:  imgdiet.FormatAVIF,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "export_AVIF_as_JPEG",
				give:        filepath.Join(_TestDataPath, _TestValidImageAVIF),
				giveFormat:  imgdiet.FormatJPEG,
				giveOptions: imgdiet.DefaultOptions(),
				want:        false,
			},
			{
				name:        "invalid_image",
				give:        filepath.Join(_TestDataPath, _TestInvalidImageJPG),
				giveFormat:  imgdiet.FormatJPEG,
				giveOptions: imgdiet.DefaultOptions(),
				want:        true,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tt.give)
			if err != nil {
				t.Fatalf("unable to open file: %v", err)
			}
			defer file.Close()

			img, err := imgdiet.Open(ctx, file)
			if err != nil {
				if tt.want {
					return
				}

				t.Fatalf("Open() failed: %v", err)
			}

			_, err = img.Export(tt.giveFormat, tt.giveOptions)
			if (err != nil) != tt.want {
				t.Errorf("Image.Export() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestImage_Export_MultipleFromSameSource(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	file, err := os.Open(filepath.Join(_TestDataPath, _TestValidImageJPG))
	if err != nil {
		t.Fatalf("unable to open file: %v", err)
	}
	defer file.Close()

	img, err := imgdiet.Open(ctx, file)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}

	opts1 := imgdiet.DefaultOptions()
	opts1.Width = 200
	opts1.Height = 200

	jpeg, err := img.Export(imgdiet.FormatJPEG, opts1)
	if err != nil {
		t.Fatalf("Export JPEG failed: %v", err)
	}

	webp, err := img.Export(imgdiet.FormatWebP, imgdiet.DefaultOptions())
	if err != nil {
		t.Fatalf("Export WebP failed: %v", err)
	}

	opts2 := imgdiet.DefaultOptions()
	opts2.Width = 100

	png, err := img.Export(imgdiet.FormatPNG, opts2)
	if err != nil {
		t.Fatalf("Export PNG failed: %v", err)
	}

	if len(jpeg) == 0 || len(webp) == 0 || len(png) == 0 {
		t.Error("one or more exports produced empty output")
	}
}
