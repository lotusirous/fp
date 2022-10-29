package main

import (
	"path/filepath"
	"testing"
)

func TestByteCountSI(t *testing.T) {
	t.Parallel()

	cases := []struct {
		num  int64
		want string
	}{
		{
			num:  int64(1000),
			want: "1.0 kB",
		},
		{
			num:  int64(1),
			want: "1 B",
		},
		{
			num:  int64(0),
			want: "0 B",
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			got := ByteCountSI(tc.num)
			if got != tc.want {
				t.Errorf("Got: %s - want: %s", got, tc.want)
			}
		})
	}
}

func TestReadFileInfo(t *testing.T) {
	testFile := "file.go"
	abs, _, err := ReadFileInfo(testFile)
	if err != nil {
		t.Fatalf("cannot read file info: %v", err)
	}

	want, _ := filepath.Abs(testFile)
	if abs != want {
		t.Fatalf("invalid abs file, want: %v - got: %v", want, abs)
	}
}
