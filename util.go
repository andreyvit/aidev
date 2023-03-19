package main

import (
	"fmt"
	"os"
	"strings"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func ensure(err error) {
	if err != nil {
		panic(err)
	}
}

func loadEnv(fn string) {
	if fn == "" || fn == "none" {
		return
	}
	raw := string(must(os.ReadFile(fn)))
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, ok := strings.Cut(line, "=")
		if !ok || strings.ContainsAny(k, " ") {
			continue
		}
		os.Setenv(strings.TrimSpace(k), strings.TrimSpace(v))
	}
}

func needEnv(name string) string {
	s := os.Getenv(name)
	if s == "" {
		fmt.Fprintf(os.Stderr, "** missing %s env variable\n", name)
		os.Exit(2)
	}
	return s
}

func subst(s string, values map[string]string) string {
	for k, v := range values {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

type action func()

func (_ action) String() string {
	return ""
}

func (_ action) IsBoolFlag() bool {
	return true
}

func (f action) Set(string) error {
	f()
	os.Exit(0)
	return nil
}
