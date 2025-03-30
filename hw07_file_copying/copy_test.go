package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	type Test struct {
		Name          string
		Offset, Limit int64
	}
	tests := []Test{
		{"out_offset0_limit0", 0, 0},
		{"out_offset0_limit10", 0, 10},
		{"out_offset0_limit1000", 0, 1000},
		{"out_offset0_limit10000", 0, 10000},
		{"out_offset100_limit1000", 100, 1000},
		{"out_offset6000_limit1000", 6000, 1000},
		{"out_offset6000_limit0", 6000, 0},
		// {"out_offset6717_limit0", 6717, 0},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			inputFileName := filepath.Join("testdata", "input.txt")
			tmpFileName := filepath.Join("testdata", "out.txt")
			goldenFileName := filepath.Join("testdata", test.Name+".txt")

			// - `os.OpenFile`, `os.Create`, `os.FileMode`
			// - `io.CopyN`
			// - `os.CreateTemp`
			// io.CopyN()

			err := Copy(inputFileName, tmpFileName, test.Offset, test.Limit)
			require.NoError(t, err)

			s, err := os.ReadFile(tmpFileName)
			require.NoError(t, err)

			g, err := os.ReadFile(goldenFileName)
			require.NoError(t, err)

			require.Equal(t, s, g)

			os.Remove(tmpFileName)
		})
	}
}

func TestCopyNegative(t *testing.T) {
	type Test struct {
		Name, srcPath, dstPath string
		Offset, Limit          int64
		err                    error
	}
	tests := []Test{
		{"unsupported_srcFile", "/dev/urandom", "out.txt", 0, 0, ErrUnsupportedFile},
		{"src_eq_dst", "testdata/input.txt", "input.txt", 0, 0, ErrFile},
		{"not_existing_src", "error_path", "out.txt", 0, 0, ErrFile},
		{"dst_problem", "testdata/input.txt", "1/1.txt", 0, 0, ErrFile},
		{"invalid_offset", "testdata/input.txt", "out.txt", 10000, 0, ErrOffsetExceedsFileSize},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			inputFileName := test.srcPath
			tmpFileName := filepath.Join("testdata", test.dstPath)
			err := Copy(inputFileName, tmpFileName, test.Offset, test.Limit)
			require.Truef(t, errors.Is(err, test.err), "actual error %q", err)
		})
	}
}
