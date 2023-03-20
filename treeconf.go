package main

import (
	"fmt"
	"os"
	"strings"
)

type TreeConfig struct {
	Includes   []string
	Excludes   []string
	Unexcludes []string
	Slices     map[string]*TreeConfig
}

func (conf *TreeConfig) Append(src *TreeConfig) {
	if src == nil {
		return
	}
	conf.Includes = append(conf.Includes, src.Includes...)
	conf.Excludes = append(conf.Excludes, src.Excludes...)
	conf.Unexcludes = append(conf.Unexcludes, src.Unexcludes...)
}

func loadTreeConfig(fn string) (*TreeConfig, error) {
	raw, err := os.ReadFile(fn)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(string(raw), "\n")
	conf := &TreeConfig{
		Slices: make(map[string]*TreeConfig),
	}

	var sliceName string
	for lno, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		directive, args, _ := strings.Cut(line, " ")
		args = strings.TrimSpace(args)
		switch directive {
		case "only":
			patterns := strings.Fields(args)
			if len(patterns) == 0 {
				return nil, fmt.Errorf("%s:%d: empty %s directive", fn, lno+1, directive)
			}
			if sliceName == "" {
				conf.Includes = append(conf.Includes, patterns...)
			} else {
				conf.Slices[sliceName].Includes = append(conf.Slices[sliceName].Includes, patterns...)
			}
		case "ignore":
			patterns := strings.Fields(args)
			if len(patterns) == 0 {
				return nil, fmt.Errorf("%s:%d: empty %s directive", fn, lno+1, directive)
			}
			if sliceName == "" {
				conf.Excludes = append(conf.Excludes, patterns...)
			} else {
				conf.Slices[sliceName].Excludes = append(conf.Slices[sliceName].Excludes, patterns...)
			}
		case "unignore":
			patterns := strings.Fields(args)
			if len(patterns) == 0 {
				return nil, fmt.Errorf("%s:%d: empty %s directive", fn, lno+1, directive)
			}
			if sliceName == "" {
				conf.Unexcludes = append(conf.Unexcludes, patterns...)
			} else {
				conf.Slices[sliceName].Unexcludes = append(conf.Slices[sliceName].Unexcludes, patterns...)
			}
		case "slice":
			sliceName = args
			conf.Slices[sliceName] = &TreeConfig{}
		default:
			return nil, fmt.Errorf("%s:%d: unknown directive %s", fn, lno+1, directive)
		}
	}
	return conf, nil
}
