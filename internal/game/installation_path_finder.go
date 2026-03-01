package game

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

var _ InstallationPathFinder = (*installationPathFinder)(nil)

type installationPathFinder struct{}

func NewInstallationPathFinder() *installationPathFinder {
	return &installationPathFinder{}
}

func (f *installationPathFinder) FindPath(ctx context.Context, baseDirectory string, slug string, detectionHints InstallationDetectionHints) (installationPath string, found bool, err error) {
	if len(detectionHints.FilePaths) == 0 && len(detectionHints.FolderPaths) == 0 {
		return "", false, nil
	}

	candidate := filepath.Join(baseDirectory, slug)
	matches, err := candidateMatchesAllHints(candidate, detectionHints)
	if err != nil {
		return "", false, fmt.Errorf("failed to check if candidate matches all hints: %w", err)
	}
	if matches {
		absPath, absErr := filepath.Abs(candidate)
		if absErr != nil {
			return "", false, fmt.Errorf("failed to get absolute path of candidate: %w", absErr)
		}
		return absPath, true, nil
	}

	const depth = 3
	nodes := make([][]string, depth+1)
	nodes[0] = []string{baseDirectory}

	for nodedepth := 1; nodedepth <= depth; nodedepth++ {
		for _, node := range nodes[nodedepth-1] {
			entries, readErr := os.ReadDir(node)
			if readErr != nil {
				return "", false, fmt.Errorf("failed to read directory: %w", readErr)
			}
			for _, entry := range entries {
				if entry.IsDir() {
					info, infoErr := entry.Info()
					if infoErr != nil {
						return "", false, fmt.Errorf("failed to get info of entry: %w", infoErr)
					}
					childPath := filepath.Join(node, info.Name())
					nodes[nodedepth] = append(nodes[nodedepth], childPath)
				}
			}
		}

		if len(nodes[nodedepth]) == 0 {
			break
		}

		for _, node := range nodes[nodedepth] {
			matches, matchErr := candidateMatchesAllHints(node, detectionHints)
			if matchErr != nil {
				return "", false, fmt.Errorf("failed to check if node matches all hints: %w", matchErr)
			}
			if matches {
				absPath, absErr := filepath.Abs(node)
				if absErr != nil {
					return "", false, fmt.Errorf("failed to get absolute path of node: %w", absErr)
				}
				return absPath, true, nil
			}
		}
	}
	return "", false, nil
}

func candidateMatchesAllHints(candidateDir string, hints InstallationDetectionHints) (bool, error) {
	for _, hintFilePath := range hints.FilePaths {
		normalized := filepath.FromSlash(hintFilePath)
		fullPath := filepath.Join(candidateDir, normalized)
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, fmt.Errorf("failed to stat file: %w", err)
		}
		if !info.Mode().IsRegular() {
			return false, nil
		}
	}
	for _, hintFolderPath := range hints.FolderPaths {
		normalized := filepath.FromSlash(hintFolderPath)
		fullPath := filepath.Join(candidateDir, normalized)
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, fmt.Errorf("failed to stat folder: %w", err)
		}
		if !info.IsDir() {
			return false, nil
		}
	}
	return true, nil
}
