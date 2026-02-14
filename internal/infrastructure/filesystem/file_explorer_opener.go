package filesystem

import (
	"fmt"
	"os/exec"
	"runtime"

	gameappsrv "github.com/seternate/go-lanty-client/internal/application/game/service"
)

var _ gameappsrv.GameExplorerOpener = (*FileExplorerOpener)(nil)

type FileExplorerOpener struct {
}

func NewFileExplorerOpener() *FileExplorerOpener {
	return &FileExplorerOpener{}
}

func (opener *FileExplorerOpener) OpenFileExplorer(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
