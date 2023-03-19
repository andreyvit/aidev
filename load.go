package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type item struct {
	relPath string
	content []byte
}

func loadFiles(rootDirs []string, ignore func(path string, isDir bool) bool) (matched []*item, ignored []string) {
	for _, rootDir := range rootDirs {
		rootDir = strings.TrimSuffix(rootDir, "/")
		rootDir = must(filepath.Abs(rootDir))
		ensure(filepath.Walk(rootDir, func(fn string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			rel := computeRelPath(fn, rootDir)

			if info.IsDir() {
				if ignore(fn, true) {
					ignored = append(ignored, rel+"/")
					return filepath.SkipDir
				}
				return nil
			}
			if ignore(fn, false) {
				ignored = append(ignored, rel)
				return nil
			}

			content := must(os.ReadFile(fn))
			if len(content) == 0 {
				return nil
			}
			if !utf8.Valid(content) {
				return nil
			}

			matched = append(matched, &item{
				relPath: rel,
				content: content,
			})
			return nil
		}))
	}
	return
}

func formatItems(items []*item) string {
	var buf strings.Builder
	for _, item := range items {
		buf.WriteString("=#=#= ")
		buf.WriteString(item.relPath)
		buf.WriteString("\n")
		buf.Write(bytes.TrimSpace(item.content))
		buf.WriteString("\n\n")
	}
	buf.WriteString("=#=#= END\n\n")
	return buf.String()
}

func computeRelPath(path, root string) string {
	n := len(root)
	if len(path) > n+1 && path[:n] == root && path[n] == '/' {
		return path[n+1:]
	}
	return path
}
