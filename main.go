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
	"sort"
)

// HashGroup composes the name to hash hash function
type HashGroup map[string]hash.Hash

// Values extracts the values from group hash.
func (h HashGroup) Values() map[string]string {
	out := make(map[string]string)
	for name, v := range h {
		out[name] = hex.EncodeToString(v.Sum(nil))
	}
	return out
}

// Writer composes the writer in the group to hash
func (h HashGroup) Writer() io.Writer {
	var fns []io.Writer
	for _, fn := range h {
		fns = append(fns, fn)
	}
	return io.MultiWriter(fns...)
}

// MultiDigest hashes a file for a given path
func MultiDigest(loc string, group HashGroup) (map[string]string, error) {
	w := group.Writer()
	file, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	if _, err := io.CopyBuffer(w, file, buf); err != nil {
		return nil, err
	}

	return group.Values(), nil
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

	hashFuncs := HashGroup{
		"MD5":    md5.New(),
		"SHA1":   sha1.New(),
		"SHA256": sha256.New(),
	}

	if *verbose {
		hashes, err := MultiDigest(abs, hashFuncs)
		if err != nil {
			fmt.Fprintf(fd, "Cannot digest: %v", err)
			os.Exit(1)
		}
		// sort by key since the go map is unordered
		keys := make([]string, 0, len(hashes))
		for k := range hashes {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, name := range keys {
			fmt.Fprintf(fd, "%s: %s\n", name, hashes[name])
		}
	}

}
