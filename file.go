package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadFileInfo gets the abs path and follows the link
func ReadFileInfo(path string) (string, os.FileInfo, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	origin, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", nil, err
	}
	fi, err := os.Stat(origin)
	if err != nil {
		return "", nil, err
	}

	return origin, fi, nil
}

// ByteCountSI converts a given size to human-readable format.
// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
