package main

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		desc     string
		patterns []string
		name     string
		relPath  string
		isDir    bool
		expected int
	}{
		{
			desc:     "simple include pattern",
			patterns: []string{"*.go"},
			name:     "main.go",
			relPath:  "main.go",
			isDir:    false,
			expected: 4,
		},
		{
			desc:     "path include pattern",
			patterns: []string{"frontend/*.js"},
			name:     "app.js",
			relPath:  "frontend/app.js",
			isDir:    false,
			expected: 13,
		},
		// Add more test cases
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result := match(test.name, test.relPath, test.isDir, test.patterns)
			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		name     string
		ignorer  *Ignorer
		path     string
		isDir    bool
		expected bool
	}{
		{
			name: "simple include pattern",
			ignorer: newIgnorer(&TreeConfig{
				Includes: []string{"*.go"},
			}, nil),
			path:     "main.go",
			isDir:    false,
			expected: false,
		},
		{
			name: "path include pattern",
			ignorer: newIgnorer(&TreeConfig{
				Dir:      "/tmp",
				Includes: []string{"frontend/*.js"},
			}, nil),
			path:     "/tmp/frontend/app.js",
			isDir:    false,
			expected: false,
		},
		{
			name:     "ignore .draft files",
			ignorer:  newIgnorer(&TreeConfig{}, nil),
			path:     "main.draft.go",
			isDir:    false,
			expected: true,
		},
		{
			name: "ignore directory",
			ignorer: newIgnorer(&TreeConfig{
				Excludes: []string{"ignored_directory/"},
			}, nil),
			path:     "ignored_directory",
			isDir:    true,
			expected: true,
		},
		// Add more test cases
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.ignorer.ShouldIgnore(test.path, test.isDir)
			if result != test.expected {
				t.Errorf("Expected %t, got %t", test.expected, result)
			}
		})
	}
}
