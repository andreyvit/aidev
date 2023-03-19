package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func loadFiles(buf *strings.Builder, root string, ignore func(path string, isDir bool) bool) {
	root = strings.TrimSuffix(root, "/")

	ensure(filepath.Walk(root, func(fn string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if ignore(fn, true) {
				log.Printf("ignore: %s/", fn)
				return filepath.SkipDir
			}
			return nil
		}
		if ignore(fn, false) {
			log.Printf("ignore: %s", fn)
			return nil
		}

		data := must(os.ReadFile(fn))
		if len(data) == 0 {
			return nil
		}
		if !utf8.Valid(data) {
			return nil
		}
		data = bytes.TrimSpace(data) // mostly for trailing whitespace

		buf.WriteString("=#=#= ")
		buf.WriteString(computeRelPath(fn, root))
		buf.WriteString("\n")
		buf.Write(data)
		buf.WriteString("\n\n")
		return nil
	}))
}

func computeRelPath(path, root string) string {
	n := len(root)
	if len(path) > n+1 && path[:n] == root && path[n] == '/' {
		return path[n+1:]
	}
	return path
}
