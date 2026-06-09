package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		offset   int64
		limit    int64
		wantFile string
		wantErr  error
	}{
		{
			name:     "full file limit 0",
			offset:   0,
			limit:    0,
			wantFile: "testdata/out_offset0_limit0.txt",
		},
		{
			name:     "limit 10",
			offset:   0,
			limit:    10,
			wantFile: "testdata/out_offset0_limit10.txt",
		},
		{
			name:     "limit 1000",
			offset:   0,
			limit:    1000,
			wantFile: "testdata/out_offset0_limit1000.txt",
		},
		{
			name:     "limit 10000",
			offset:   0,
			limit:    10000,
			wantFile: "testdata/out_offset0_limit10000.txt",
		},
		{
			name:     "offset 100 limit 1000",
			offset:   100,
			limit:    1000,
			wantFile: "testdata/out_offset100_limit1000.txt",
		},
		{
			name:     "offset 6000 limit 1000",
			offset:   6000,
			limit:    1000,
			wantFile: "testdata/out_offset6000_limit1000.txt",
		},
		{
			name:    "offset exceeds file size",
			offset:  10_000_000,
			limit:   0,
			wantErr: ErrOffsetExceedsFileSize,
		},
		{
			name:     "limit greater than file copies until EOF",
			offset:   0,
			limit:    10_000_000,
			wantFile: "testdata/out_offset0_limit0.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			outPath := filepath.Join(t.TempDir(), "out.txt")
			err := Copy("testdata/input.txt", outPath, tt.offset, tt.limit)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			got, err := os.ReadFile(outPath)
			require.NoError(t, err)
			want, err := os.ReadFile(tt.wantFile)
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}
}
