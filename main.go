package main

import (
	"fmt"
	"os"
	"path/filepath" // Re-added filepath import

	"cursezip/archiver"
	"cursezip/config"
	"cursezip/packer"

	"github.com/spf13/cobra"
)

var (
	format     string
	configFile string
	excludePatterns []string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cursezip <source_directory_1> [source_directory_2...]",
		Short: "Packages one or more CurseForge plugin directories into a single archive.",
		Long: `CurseZip helps you create archives of your CurseForge plugins, allowing you to exclude unnecessary files like .git repositories or macOS system files.

The output archive will be named after the first source directory and placed in its parent directory.`,
		Args:  cobra.MinimumNArgs(1), // Requires at least one source directory
		RunE:  runPack, // runPack function will now be the main command's logic
	}

	// Define flags for the root command
	rootCmd.Flags().StringVarP(&format, "format", "f", "zip", "Output archive format (zip, tar.gz)")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to exclusion configuration file")
	rootCmd.Flags().StringArrayVarP(&excludePatterns, "exclude", "e", []string{}, "Additional exclusion patterns (e.g., *.tmp, **/test/)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runPack(cmd *cobra.Command, args []string) error {
	sourceDirs := args // All positional arguments are source directories

	// Determine output file name and path
	firstSourceDir := sourceDirs[0]
	outputFileName := filepath.Base(firstSourceDir) + "." + format
	if format == "tar.gz" {
		outputFileName = filepath.Base(firstSourceDir) + ".tar.gz"
	}
	outputDir := filepath.Dir(firstSourceDir)
	if outputDir == "" { // If source is in current dir, output to current dir
		outputDir = "."
	}
	outputFilePath := filepath.Join(outputDir, outputFileName)

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Add command-line exclusion patterns to config
	cfg.Exclude = append(cfg.Exclude, excludePatterns...)

	// Create filter
	fileFilter := packer.NewFilter(cfg.Exclude)

	// Create packer with multiple source directories
	filePacker := packer.NewPacker(sourceDirs, fileFilter)

	// Pack files
	packResult, err := filePacker.Pack()
	if err != nil {
		return fmt.Errorf("failed to pack files: %w", err)
	}

	if len(packResult.Files) == 0 {
		fmt.Printf("No files to archive after applying exclusion rules.\n")
		return nil
	}

	// Create output file
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFilePath, err)
	}
	defer outFile.Close()

	var arch archiver.Archiver
	switch format {
	case "zip":
		arch = archiver.NewZipArchiver()
	case "tar.gz":
		arch = archiver.NewTarGzArchiver()
	// TODO: Add support for 7z and rar
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}

	fmt.Printf("Archiving %d files to '%s' using format '%s'...\n", len(packResult.Files), outputFilePath, format)
	if err := arch.Archive(packResult.Files, outFile); err != nil { // Pass map to Archive
		return fmt.Errorf("failed to create archive: %w", err)
	}

	fmt.Printf("Successfully created archive: %s\n", outputFilePath)
	return nil
}
