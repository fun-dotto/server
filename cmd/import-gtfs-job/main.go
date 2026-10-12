package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: go run ./cmd/import-gtfs-job <zip-url>")
	}
	if err := downloadAndExtract(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func downloadAndExtract(url string) error {
	tempDir, err := os.MkdirTemp("", "gtfs-import-")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download ZIP: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download ZIP: unexpected HTTP status %s", response.Status)
	}

	archivePath := filepath.Join(tempDir, "feed.zip")
	archive, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("create temporary ZIP file: %w", err)
	}
	if _, err := io.Copy(archive, response.Body); err != nil {
		archive.Close()
		return fmt.Errorf("save downloaded ZIP: %w", err)
	}
	if err := archive.Close(); err != nil {
		return fmt.Errorf("close downloaded ZIP: %w", err)
	}

	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("open downloaded ZIP: %w", err)
	}
	defer reader.Close()

	extractDir := filepath.Join(tempDir, "extracted")
	for _, file := range reader.File {
		if !filepath.IsLocal(file.Name) {
			return fmt.Errorf("ZIP contains an unsafe file path: %q", file.Name)
		}

		path := filepath.Join(extractDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("create directory %q: %w", file.Name, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("create parent directory for %q: %w", file.Name, err)
		}
		source, err := file.Open()
		if err != nil {
			return fmt.Errorf("open ZIP entry %q: %w", file.Name, err)
		}
		destination, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			source.Close()
			return fmt.Errorf("create extracted file %q: %w", file.Name, err)
		}
		_, copyErr := io.Copy(destination, source)
		sourceErr := source.Close()
		destinationErr := destination.Close()
		if copyErr != nil {
			return fmt.Errorf("extract file %q: %w", file.Name, copyErr)
		}
		if sourceErr != nil {
			return fmt.Errorf("close ZIP entry %q: %w", file.Name, sourceErr)
		}
		if destinationErr != nil {
			return fmt.Errorf("close extracted file %q: %w", file.Name, destinationErr)
		}
		fmt.Println(file.Name)
	}

	return nil
}
