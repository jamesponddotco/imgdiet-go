package imgdiet_test

import (
	"reflect"
	"runtime"
	"testing"

	"git.sr.ht/~jamesponddotco/imgdiet-go"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want *imgdiet.Config
		name string
	}{
		{
			name: "Test Default Config",
			want: &imgdiet.Config{
				Cache:          1024 * 1024 * 1024,
				MaxConcurrency: runtime.NumCPU(),
				ReportLeaks:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := imgdiet.DefaultConfig()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
