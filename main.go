package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// Digest hashes a file for a given path
func Digest(h hash.Hash, loc string) (string, error) {
	file, err := os.Open(loc)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	if _, err := io.CopyBuffer(h, file, buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// ReadFileInfo gets the abs path and follows the link
func ReadFileInfo(path string) (string, os.FileInfo, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	linkAbs, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", nil, err
	}
	fi, err := os.Stat(linkAbs)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil, err
	} else if err != nil {
		return "", nil, err
	}

	return linkAbs, fi, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: fp [file or directory]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func exit(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, a...)
	os.Exit(1)
}

var (
	verbose = flag.Bool("v", false, `print the file hash in (md5, sha1, sha256)`)
)

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	fd := os.Stdout // default writer
	path := flag.Arg(0)

	abs, fi, err := ReadFileInfo(path)
	if err != nil {
		exit(fd, "Cannot resolve the path: ", err)
	}

	fmt.Fprintf(fd, "PATH: %s\n", abs)
	if fi.IsDir() {
		os.Exit(0)
	}

	if *verbose {
		hashes := map[string]hash.Hash{
			"MD5":    md5.New(),
			"SHA1":   sha1.New(),
			"SHA256": sha256.New(),
		}
		for name, hf := range hashes {
			if v, err := Digest(hf, abs); err != nil {
				fmt.Fprintf(fd, "Cannot digest %s: %v", name, err)
				os.Exit(1)
			} else {
				fmt.Fprintf(fd, "%s: %s\n", name, v)
			}
		}
	}

}
