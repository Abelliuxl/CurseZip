package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// ZipArchiver implements the Archiver interface for ZIP format.
type ZipArchiver struct{}

// NewZipArchiver creates a new ZipArchiver.
func NewZipArchiver() *ZipArchiver {
	return &ZipArchiver{}
}

// Archive creates a ZIP archive from the given files.
func (za *ZipArchiver) Archive(files map[string]string, output io.Writer) error {
	zipWriter := zip.NewWriter(output)
	defer zipWriter.Close()

	for archivePath, diskPath := range files {
		info, err := os.Stat(diskPath)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", diskPath, err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create file info header for %s: %w", diskPath, err)
		}

		// Set the name in the archive to the provided archivePath
		header.Name = archivePath

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip entry for %s: %w", archivePath, err)
		}

		if !info.IsDir() {
			srcFile, err := os.Open(diskPath)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", diskPath, err)
			}
			defer srcFile.Close()

			_, err = io.Copy(writer, srcFile)
			if err != nil {
				return fmt.Errorf("failed to copy file %s to zip: %w", diskPath, err)
			}
		}
	}
	return nil
}
