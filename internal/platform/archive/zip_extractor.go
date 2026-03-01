package archive

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/seternate/go-lanty-client/internal/game"
)

type ZipExtractor struct {
	mu       sync.RWMutex
	progress map[string]game.InstallationProgress
}

func NewZipExtractor() *ZipExtractor {
	return &ZipExtractor{
		progress: make(map[string]game.InstallationProgress),
	}
}

func (extractor *ZipExtractor) Extract(ctx context.Context, slug string, archivePath string, extractionPath string, removeArchive bool) error {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip reader: %w", err)
	}
	defer zipReader.Close()

	err = os.MkdirAll(extractionPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}

	extractor.mu.Lock()
	extractor.progress[slug] = game.InstallationProgress{
		Slug:      slug,
		Total:     uint64(len(zipReader.File)),
		Completed: 0,
		StartedAt: time.Now(),
	}
	extractor.mu.Unlock()

	var extractionErr error
extractionLoop:
	for _, zipFile := range zipReader.File {
		select {
		case <-ctx.Done():
			extractionErr = ctx.Err()
			break extractionLoop
		default:
		}

		extractionErr = extractZipFile(zipFile, extractionPath)
		if extractionErr != nil {
			break extractionLoop
		}

		extractor.mu.Lock()
		progress := extractor.progress[slug]
		progress.Completed++
		extractor.progress[slug] = progress
		extractor.mu.Unlock()
	}

	extractor.mu.Lock()
	defer extractor.mu.Unlock()

	if extractionErr != nil {
		progress := extractor.progress[slug]
		progress.EndedAt = time.Now()
		progress.Failed = true
		extractor.progress[slug] = progress
		return fmt.Errorf("failed to extract zip file: %w", extractionErr)
	}

	progress := extractor.progress[slug]
	progress.EndedAt = time.Now()
	progress.Succeeded = true
	extractor.progress[slug] = progress

	if removeArchive {
		err := os.Remove(archivePath)
		if err != nil {
			return fmt.Errorf("failed to remove archive: %w", err)
		}
	}

	return nil
}

func extractZipFile(zipFile *zip.File, dest string) error {
	cleanDest := filepath.Clean(dest)

	targetPath := filepath.Join(cleanDest, zipFile.Name)
	targetPath = filepath.Clean(targetPath)

	if !strings.HasPrefix(targetPath, cleanDest+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path in zip: %s", zipFile.Name)
	}

	if zipFile.FileInfo().IsDir() {
		err := os.MkdirAll(targetPath, zipFile.Mode())
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		return nil
	}

	err := os.MkdirAll(filepath.Dir(targetPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	src, err := zipFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func (extractor *ZipExtractor) GetAllInstallationProgresses(ctx context.Context) (map[string]game.InstallationProgress, error) {
	extractor.mu.RLock()
	defer extractor.mu.RUnlock()

	return maps.Clone(extractor.progress), nil
}
