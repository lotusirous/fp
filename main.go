package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// DigestGroup hashes a file for a given path
func DigestGroup(hg map[string]hash.Hash, loc string) (map[string]string, error) {
	var fns []io.Writer
	for _, fn := range hg {
		fns = append(fns, fn)
	}
	w := io.MultiWriter(fns...)

	file, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	if _, err := io.CopyBuffer(w, file, buf); err != nil {
		return nil, err
	}

	out := make(map[string]string)
	for name, v := range hg {
		out[name] = hex.EncodeToString(v.Sum(nil))
	}

	return out, nil
}

func writeTo(values map[string]string, w io.Writer) {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		fmt.Fprintf(w, "%s: %s\n", name, values[name])
	}
}

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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: fp [file or directory]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
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
		fmt.Fprintln(fd, "Cannot resolve the path: ", err)
		os.Exit(1)
	}

	fmt.Fprintf(fd, "PATH: %s\n", abs)
	if *verbose && !fi.IsDir() {
		hg := make(map[string]hash.Hash)
		hg["MD5"] = md5.New()
		hg["SHA1"] = sha1.New()
		hg["SHA256"] = sha256.New()

		values, err := DigestGroup(hg, abs)
		if err != nil {
			fmt.Fprintf(fd, "Cannot digest group: %v", err)
			os.Exit(1)
		}
		writeTo(values, fd)
	}

}
