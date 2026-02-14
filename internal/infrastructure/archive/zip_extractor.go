package archive

import (
	"archive/zip"
	"context"
	"errors"
	"io"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gameapp "github.com/seternate/go-lanty-client/internal/application/game"
	gameappsrv "github.com/seternate/go-lanty-client/internal/application/game/service"
	"github.com/seternate/go-lanty-client/internal/application/setting"
	settingappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
)

var _ gameappsrv.BlobExtractor = (*ZipExtractor)(nil)
var _ settingappsrv.SettingsListener = (*ZipExtractor)(nil)
var _ gameappsrv.GameInstallationProgressFetcher = (*ZipExtractor)(nil)

type ZipExtractor struct {
	OutputDir string
	mu        sync.RWMutex
	progress  map[string]gameapp.InstallationProgress
}

func NewZipExtractor(outputDir string) *ZipExtractor {
	return &ZipExtractor{
		OutputDir: outputDir,
		progress:  make(map[string]gameapp.InstallationProgress),
	}
}

func (extractor *ZipExtractor) Extract(ctx context.Context, slug string, filePath string, removeArchive bool) (string, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer zipReader.Close()

	extractor.mu.RLock()
	outputDir := extractor.OutputDir
	extractor.mu.RUnlock()

	extractionPath := filepath.Join(outputDir, slug)

	err = os.MkdirAll(extractionPath, 0755)
	if err != nil {
		return "", err
	}

	extractor.mu.Lock()
	extractor.progress[slug] = gameapp.InstallationProgress{
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
		return "", extractionErr
	}

	progress := extractor.progress[slug]
	progress.EndedAt = time.Now()
	progress.Succeeded = true
	extractor.progress[slug] = progress

	if removeArchive {
		os.Remove(filePath)
	}

	relativeExtractionPath, err := filepath.Rel(outputDir, extractionPath)
	if err != nil {
		return slug, err
	}

	return relativeExtractionPath, nil
}

func extractZipFile(zipFile *zip.File, dest string) error {
	cleanDest := filepath.Clean(dest)

	targetPath := filepath.Join(cleanDest, zipFile.Name)
	targetPath = filepath.Clean(targetPath)

	if !strings.HasPrefix(targetPath, cleanDest+string(os.PathSeparator)) {
		return errors.New("illegal file path in zip: " + zipFile.Name)
	}

	if zipFile.FileInfo().IsDir() {
		return os.MkdirAll(targetPath, zipFile.Mode())
	}

	err := os.MkdirAll(filepath.Dir(targetPath), 0755)
	if err != nil {
		return err
	}

	src, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (extractor *ZipExtractor) ApplySettings(settings setting.Settings) {
	extractor.mu.Lock()
	extractor.OutputDir = settings.GameDirectory
	extractor.mu.Unlock()
}

func (extractor *ZipExtractor) GetAllProgress() (map[string]gameapp.InstallationProgress, error) {
	extractor.mu.RLock()
	defer extractor.mu.RUnlock()

	return maps.Clone(extractor.progress), nil
}

func (extractor *ZipExtractor) GetProgress(slug string) (gameapp.InstallationProgress, error) {
	extractor.mu.RLock()
	defer extractor.mu.RUnlock()

	progress, ok := extractor.progress[slug]
	if !ok {
		return gameapp.InstallationProgress{}, errors.New("progress not found")
	}

	return progress, nil
}
