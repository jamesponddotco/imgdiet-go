package imgdiet_test

import (
	"reflect"
	"runtime"
	"testing"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
	"github.com/davidbyttow/govips/v2/vips"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *imgdiet.Config
	}{
		{
			name: "Test Default Config",
			want: &imgdiet.Config{
				Logger:         imgdiet.DefaultLogger,
				LogLevel:       vips.LogLevelError,
				Cache:          1024 * 1024 * 1024,
				MaxConcurrency: runtime.NumCPU(),
				ReportLeaks:    false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := imgdiet.DefaultConfig()

			if got.Logger == nil {
				t.Errorf("Expected Logger to not be nil")
			}

			// Ignore Logger for reflect.DeepEqual.
			got.Logger = nil
			tt.want.Logger = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
