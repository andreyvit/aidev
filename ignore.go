package main

import (
	"path/filepath"
	"regexp"
)

type Ignorer struct {
}

func (ign *Ignorer) ShouldIgnore(path string, isDir bool) bool {
	if path == "." {
		return false
	}
	_, name := filepath.Split(path)
	if ignoredNames[name] {
		return true
	}
	if ignoredRe.MatchString(name) {
		return true
	}
	if isDir && ignoredDirRe.MatchString(name) {
		return true
	}
	return false
}

func newIgnorer() *Ignorer {
	return &Ignorer{}
}

var ignoredRe = regexp.MustCompile(`~`)
var ignoredDirRe = regexp.MustCompile(`^\.`)

var ignoredNames = map[string]bool{
	"node_modules": true,
}
