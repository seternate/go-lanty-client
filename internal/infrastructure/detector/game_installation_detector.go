package detector

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/seternate/go-lanty-client/internal/application/game/service"
	"github.com/seternate/go-lanty-client/internal/application/setting"
	settingappsrv "github.com/seternate/go-lanty-client/internal/application/setting/service"
)

var _ service.InstallationDetector = (*GameInstallationDetector)(nil)
var _ settingappsrv.SettingsListener = (*GameInstallationDetector)(nil)

type GameInstallationDetector struct {
	InstallationDirectory string
	mu                    sync.RWMutex
}

func NewGameInstallationDetector(installationDirectory string) *GameInstallationDetector {
	return &GameInstallationDetector{
		InstallationDirectory: installationDirectory,
	}
}

func (detector *GameInstallationDetector) DetectInstalledGame(slug string, executableRelativePath string) (bool, string, error) {
	detector.mu.RLock()
	installationDirectory := detector.InstallationDirectory
	detector.mu.RUnlock()

	path := filepath.Join(installationDirectory, slug, executableRelativePath)
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		if err != nil {
			return true, slug, err
		}
	}

	depth := 3
	nodes := make([][]string, depth+1)
	nodes[0] = []string{installationDirectory}

	for nodedepth := 1; nodedepth <= depth; nodedepth++ {
		for _, node := range nodes[nodedepth-1] {
			childnodes, err := os.ReadDir(node)
			if err != nil {
				return false, "", err
			}
			for _, childnode := range childnodes {
				if childnode.IsDir() {
					childnodeInfo, err := childnode.Info()
					if err != nil {
						return false, "", err
					}
					childnodepath := filepath.Join(node, childnodeInfo.Name())
					nodes[nodedepth] = append(nodes[nodedepth], childnodepath)
				}
			}
		}

		if len(nodes[nodedepth]) == 0 {
			break
		}

		for _, node := range nodes[nodedepth] {
			path := filepath.Join(node, executableRelativePath)
			_, err := os.Stat(path)
			if !os.IsNotExist(err) {
				relativePath, err := filepath.Rel(installationDirectory, node)
				if err != nil {
					return false, "", err
				}
				return true, relativePath, nil
			}
		}
	}

	return false, "", nil
}

func (detector *GameInstallationDetector) ApplySettings(settings setting.Settings) {
	detector.mu.Lock()
	detector.InstallationDirectory = settings.GameDirectory
	detector.mu.Unlock()
}
