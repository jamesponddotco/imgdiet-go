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
	_TestNonExistentImage string = "impossible-girl.jpg"
)

func TestMain(m *testing.M) {
	imgdiet.Start(nil)
	defer imgdiet.Stop()

	m.Run()
}

func TestDetectImageType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want string
		err  bool
	}{
		{
			name: "jpeg",
			give: _TestDataPath + "/" + _TestValidImageJPG,
			want: "JPEG",
			err:  false,
		},
		{
			name: "png",
			give: _TestDataPath + "/" + _TestValidImagePNG,
			want: "PNG",
			err:  false,
		},
		{
			name: "gif",
			give: _TestDataPath + "/" + _TestValidImageGIF,
			want: "",
			err:  true,
		},
		{
			name: "invalid",
			give: _TestDataPath + "/" + _TestInvalidImageJPG,
			want: "",
			err:  true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.ReadFile(tt.give)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			got, err := imgdiet.DetectImageType(file)
			if err == nil && tt.err {
				t.Fatalf("expected error, got none")
			}

			if err != nil && !tt.err {
				t.Fatalf("expected no error, got: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestDetectImageSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want int64
	}{
		{
			name: "JPG",
			give: _TestDataPath + "/" + _TestValidImageJPG,
			want: 1331858,
		},
		{
			name: "PNG",
			give: _TestDataPath + "/" + _TestValidImagePNG,
			want: 20076,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.ReadFile(tt.give)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			got := imgdiet.DetectImageSize(file)
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
