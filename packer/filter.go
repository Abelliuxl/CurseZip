package packer

import (
	"path/filepath"
	"strings"
)

// Filter manages exclusion patterns.
type Filter struct {
	excludePatterns []string
}

// NewFilter creates a new Filter with the given exclusion patterns.
func NewFilter(patterns []string) *Filter {
	return &Filter{
		excludePatterns: patterns,
	}
}

// ShouldExclude checks if a given path should be excluded based on the filter's patterns.
// It handles both file and directory exclusion.
// Paths are expected to be relative to the source directory.
func (f *Filter) ShouldExclude(path string, isDir bool) bool {
	// Normalize path to use forward slashes for consistent matching
	path = filepath.ToSlash(path)

	for _, pattern := range f.excludePatterns {
		// Normalize pattern
		pattern = filepath.ToSlash(pattern)

		// Handle directory patterns (e.g., "dir/" or "dir")
		if strings.HasSuffix(pattern, "/") {
			// If pattern ends with '/', it matches directories and their contents
			dirPattern := strings.TrimSuffix(pattern, "/")
			if isDir {
				// Match directory itself
				matched, _ := filepath.Match(dirPattern, path)
				if matched {
					return true
				}
			}
			// Match contents within the directory
			if strings.HasPrefix(path, dirPattern+"/") {
				return true
			}
		} else {
			// Handle file patterns or exact directory names without trailing slash
			matched, _ := filepath.Match(pattern, path)
			if matched {
				return true
			}
		}
	}
	return false
}
