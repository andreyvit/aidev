package main

import (
	"log"
	"path/filepath"
	"strings"
)

type Ignorer struct {
	overrides  *TreeConfig
	dirConfigs map[string]*TreeConfig
}

func newIgnorer(overrideConfig *TreeConfig) *Ignorer {
	return &Ignorer{
		overrides:  overrideConfig,
		dirConfigs: make(map[string]*TreeConfig),
	}
}

func (ign *Ignorer) ShouldIgnore(path string, isDir bool) bool {
	if path == "." {
		return false
	}

	dir, name := filepath.Split(path)
	if strings.Contains(name, ".draft.") || strings.HasSuffix(name, ".draft") {
		return true
	}

	conf := &TreeConfig{}
	conf.Append(builtinConfig)
	conf.Append(ign.loadConfig(dir))
	conf.Append(ign.overrides)

	incl := match(name, isDir, conf.Includes)
	excl := match(name, isDir, conf.Excludes)
	unex := match(name, isDir, conf.Unexcludes)

	if len(conf.Includes) > 0 && incl == 0 {
		return true
	}
	if excl > 0 && excl >= incl && unex == 0 {
		// in case of -i and -x conflict, longest pattern wins
		return true
	}

	return false
}

// loadConfig loads the configuration for a given directory and caches it.
func (ign *Ignorer) loadConfig(dir string) *TreeConfig {
	conf := ign.dirConfigs[dir]
	if conf != nil {
		return conf
	}

	// log.Printf("loadConfig loading for %s", dir)
	conf, err := loadTreeConfig(filepath.Join(dir, ".aidev"))
	if err != nil {
		log.Fatalf("** %v", err)
	}

	parent := filepath.Dir(dir)
	if dir == "/" || parent == "" {
	} else {
		combined := &TreeConfig{}
		combined.Append(ign.loadConfig(parent))
		combined.Append(conf)
		conf = combined
	}

	ign.dirConfigs[dir] = conf
	// log.Printf("loadConfig loaded for %s", dir)
	return conf
}

// match checks if a given path matches any of the patterns in the list and
// returns the length of the longest matching pattern.
func match(path string, isDir bool, list []string) int {
	var score int
	for _, item := range list {
		var wantsDir bool
		item, wantsDir = strings.CutSuffix(item, "/")
		if wantsDir && !isDir {
			continue
		}

		if must(filepath.Match(item, path)) {
			if len(item) > score {
				score = len(item)
			}
		}
	}
	return score
}

var builtinConfig = &TreeConfig{
	Excludes: []string{
		"go.sum",
		".env*",
		"~*",
		".gitignore",
		".git/",
		".svn/",
		".*/",
		".aidev",
		"modd*.conf",
		"node_modules/",
	},
}
