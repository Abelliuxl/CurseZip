package archiver

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// TarGzArchiver implements the Archiver interface for TAR.GZ format.
type TarGzArchiver struct{}

// NewTarGzArchiver creates a new TarGzArchiver.
func NewTarGzArchiver() *TarGzArchiver {
	return &TarGzArchiver{}
}

// Archive creates a TAR.GZ archive from the given files.
func (tga *TarGzArchiver) Archive(files map[string]string, output io.Writer) error {
	gzipWriter := gzip.NewWriter(output)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for archivePath, diskPath := range files {
		info, err := os.Stat(diskPath)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", diskPath, err)
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create file info header for %s: %w", diskPath, err)
		}

		// Set the name in the archive to the provided archivePath
		header.Name = archivePath

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header for %s: %w", archivePath, err)
		}

		if !info.IsDir() {
			srcFile, err := os.Open(diskPath)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", diskPath, err)
			}
			defer srcFile.Close()

			_, err = io.Copy(tarWriter, srcFile)
			if err != nil {
				return fmt.Errorf("failed to copy file %s to tar.gz: %w", diskPath, err)
			}
		}
	}
	return nil
}
