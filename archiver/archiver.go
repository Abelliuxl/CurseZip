package archiver

import "io"

// Archiver defines the interface for different archiving formats.
type Archiver interface {
	Archive(files map[string]string, output io.Writer) error // Changed signature
}
