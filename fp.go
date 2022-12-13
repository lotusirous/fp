package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: fp [file or directory]\n")
	fmt.Fprintf(os.Stderr, "\nThis command allows you to view information about a specified file or directory.\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	flagVerbose = flag.Bool("d", false, `Print detailed file hash information, including the md5, sha1, and sha256 hashes`)
	flagClip    = flag.Bool("c", false, "Copy the path to the clipboard")
)

func main() {
	log.SetPrefix("fp: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	path := flag.Arg(0)

	abs, fi, err := ReadFileInfo(path)
	if err != nil {
		log.Fatal(err)
	}

	rr := NewRowRender()

	if *flagVerbose {
		if fi.IsDir() {
			rr.AddRow("DIR", abs)
		} else {
			rr.AddRow("FILE", abs)
		}
		rr.AddRow("SIZE", ByteCountSI(fi.Size()))
		hg := make(map[string]hash.Hash)
		hg["MD5"] = md5.New()
		hg["SHA1"] = sha1.New()
		hg["SHA256"] = sha256.New()

		values, err := DigestGroup(hg, abs)
		if err != nil {
			log.Fatal("unable to digest group: ", err)
		}
		rr.AddRowMap(values)
		rr.Write(os.Stdout)
		return
	}
	rr.AddRow("", abs)

	if *flagClip {
		buf := new(bytes.Buffer)
		rr.Write(buf)
		line := buf.String()
		clipboard.WriteAll(strings.TrimSuffix(line, "\n"))
		return
	}
	rr.Write(os.Stdout)
}
