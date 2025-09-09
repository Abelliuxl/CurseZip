package packer

import (
	"fmt"
	"os"
	"path/filepath"
)

// PackResult holds the list of files to be packed.
type PackResult struct {
	Files map[string]string // Map of archive_path -> absolute_disk_path
}

// Packer handles the process of traversing directories and filtering files.
type Packer struct {
	sourceDirs []string // Changed to slice
	filter     *Filter
}

// NewPacker creates a new Packer instance.
func NewPacker(sourceDirs []string, filter *Filter) *Packer { // Changed to slice
	return &Packer{
		sourceDirs: sourceDirs,
		filter:     filter,
	}
}

// Pack traverses the source directories, applies filters, and returns the list of files to pack.
func (p *Packer) Pack() (*PackResult, error) {
	result := &PackResult{
		Files: make(map[string]string),
	}

	for _, sourceDir := range p.sourceDirs { // Iterate over multiple source directories
		absSourceDir, err := filepath.Abs(sourceDir)
		if err != nil {
			return nil, fmt.Errorf("invalid source directory path %s: %w", sourceDir, err)
		}

		err = filepath.Walk(absSourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(absSourceDir, path) // Relative to current sourceDir
			if err != nil {
				return err
			}

			if relPath == "." {
				return nil
			}

			normalizedRelPath := filepath.ToSlash(relPath)

			if p.filter.ShouldExclude(normalizedRelPath, info.IsDir()) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// Determine the path inside the archive.
			// If there's only one source directory, keep the path flat.
			// If multiple, prefix with the source directory name to avoid clashes.
			archivePath := normalizedRelPath
			if len(p.sourceDirs) > 1 {
				// Use the base name of the source directory as a prefix in the archive
				baseSourceDir := filepath.Base(absSourceDir)
				archivePath = filepath.Join(baseSourceDir, normalizedRelPath)
				archivePath = filepath.ToSlash(archivePath) // Ensure forward slashes
			}

			if info.IsDir() {
				result.Files[archivePath+"/"] = path // Store absolute path for directories
			} else {
				result.Files[archivePath] = path // Store absolute path for files
			}
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error walking source directory %s: %w", sourceDir, err)
		}
	}

	return result, nil
}
