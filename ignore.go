package main

import (
	"path/filepath"
	"strings"
)

type Ignorer struct {
	Include []string
	Exclude []string
}

func (ign *Ignorer) ShouldIgnore(path string, isDir bool) bool {
	if path == "." {
		return false
	}

	_, name := filepath.Split(path)
	if strings.Contains(name, ".draft.") {
		return true
	}

	incl := match(name, ign.Include)
	excl := match(name, ign.Exclude)

	if len(ign.Include) > 0 && incl == 0 {
		return true
	}
	if excl > 0 && excl >= incl {
		// in case of -i and -x conflict, longest pattern wins
		return true
	}

	if incl == 0 {
		if isDir && match(name, standardDirExcludes) > 0 {
			return true
		}
		if match(name, standardExcludes) > 0 {
			return true
		}
	}

	return false
}

func match(path string, list []string) int {
	var score int
	for _, item := range list {
		if must(filepath.Match(item, path)) {
			if len(item) > score {
				score = len(item)
			}
		}
	}
	return score
}

var standardExcludes = []string{
	"go.sum",
	".env*",
	"~*",
	".gitignore",
}

var standardDirExcludes = []string{
	".git",
	".svn",
	".*",
	"node_modules",
}
