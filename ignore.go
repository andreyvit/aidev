package main

import (
	"log"
	"path/filepath"
	"strings"
)

type Ignorer struct {
	overrides     *TreeConfig
	dirConfigs    map[string]*TreeConfig
	selectedSlice string
}

func newIgnorer(overrideConfig *TreeConfig, selectedSlice string) *Ignorer {
	return &Ignorer{
		overrides:     overrideConfig,
		dirConfigs:    make(map[string]*TreeConfig),
		selectedSlice: selectedSlice,
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

	confs := []*TreeConfig{builtinConfig}
	for c := ign.loadConfig(dir); c != nil; c = c.parent {
		confs = append(confs, c)
		if ign.selectedSlice != "" {
			if sliceConf, ok := c.Slices[ign.selectedSlice]; ok {
				confs = append(confs, sliceConf)
			}
		}
	}
	confs = append(confs, ign.overrides)
	// log.Printf("confs: %s", must(json.MarshalIndent(confs, "", "  ")))

	var incl, excl, unex int
	for _, conf := range confs {
		var relPath string
		if conf.Dir != "" {
			relPath = must(filepath.Rel(conf.Dir, path))
		}
		incl = max(incl, match(name, relPath, isDir, conf.Includes))
		excl = max(excl, match(name, relPath, isDir, conf.Excludes))
		unex = max(unex, match(name, relPath, isDir, conf.Unexcludes))
	}

	var hasIncludes bool
	for _, conf := range confs {
		if len(conf.Includes) > 0 {
			hasIncludes = true
		}
	}

	log.Printf("checking: %s (incl = %d, excl = %d, unex = %d, hasinc = %v)", path, incl, excl, unex, hasIncludes)
	if !isDir && hasIncludes && incl == 0 {
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
	conf, err := loadTreeConfig(dir)
	if err != nil {
		log.Fatalf("** %v", err)
	}

	parent := filepath.Dir(dir)
	if dir == "/" || parent == "" || parent == dir {
	} else {
		parentConf := ign.loadConfig(parent)
		if parentConf != nil {
			if conf == nil {
				conf = &TreeConfig{}
			}
			conf.parent = ign.loadConfig(parent)
		}
	}

	ign.dirConfigs[dir] = conf
	// log.Printf("loadConfig loaded for %s", dir)
	return conf
}

// match checks if a given path matches any of the patterns in the list and
// returns the length of the longest matching pattern.
func match(name string, relPath string, isDir bool, list []string) int {
	var score int
	for _, item := range list {
		var wantsDir bool
		item, wantsDir = strings.CutSuffix(item, "/")
		if wantsDir && !isDir {
			continue
		}

		if strings.Contains(item, "/") {
			if relPath != "" && must(filepath.Match(item, relPath)) {
				if len(item) > score {
					score = len(item)
				}
			}
		} else {
			if must(filepath.Match(item, name)) {
				if len(item) > score {
					score = len(item)
				}
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
