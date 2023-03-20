package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andreyvit/aidev/internal/clipboard"
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

func mustSkippingOSNotExists[T any](v T, err error) T {
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	return v
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

type stringList []string

func (v stringList) String() string {
	return strings.Join(v, " | ")
}

func (v *stringList) Set(str string) error {
	*v = append(*v, str)
	return nil
}

func saveText(fn string, content string) error {
	switch fn {
	case "":
		return nil
	case "-":
		os.Stdout.WriteString(content)
		if !strings.HasSuffix(content, "\n") {
			os.Stdout.WriteString(" ")
		}
		return nil
	case "copy":
		return clipboard.CopyText(content)
	default:
		err := os.MkdirAll(filepath.Dir(fn), 0755)
		if err != nil {
			return err
		}
		return os.WriteFile(fn, []byte(content), 0644)
	}
}

type choiceFlag[T comparable] struct {
	ptr   *T
	value T
}

func (f *choiceFlag[T]) String() string {
	if *f.ptr == f.value {
		return "true"
	} else {
		return ""
	}
}

func (_ *choiceFlag[T]) IsBoolFlag() bool {
	return true
}

func (f *choiceFlag[T]) Set(str string) error {
	*f.ptr = f.value
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
