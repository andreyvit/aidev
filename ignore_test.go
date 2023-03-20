package main

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		dir      string
		isDir    bool
		expected int
	}{
		{
			name:     "simple include pattern",
			patterns: []string{"*.go"},
			path:     "main.go",
			dir:      ".",
			isDir:    false,
			expected: 4,
		},
		{
			name:     "path include pattern",
			patterns: []string{"frontend/*.js"},
			path:     "frontend/app.js",
			dir:      "frontend",
			isDir:    false,
			expected: 12,
		},
		// Add more test cases
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := match(test.path, test.dir, test.isDir, test.patterns)
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
			}, ""),
			path:     "main.go",
			isDir:    false,
			expected: false,
		},
		{
			name: "path include pattern",
			ignorer: newIgnorer(&TreeConfig{
				Includes: []string{"frontend/*.js"},
			}, ""),
			path:     "frontend/app.js",
			isDir:    false,
			expected: false,
		},
		{
			name: "ignore .draft files",
			ignorer: newIgnorer(&TreeConfig{}, ""),
			path:     "main.draft.go",
			isDir:    false,
			expected: true,
		},
		{
			name: "ignore directory",
			ignorer: newIgnorer(&TreeConfig{
				Excludes: []string{"ignored_directory/"},
			}, ""),
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
