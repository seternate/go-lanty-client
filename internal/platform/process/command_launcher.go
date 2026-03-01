package process

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/seternate/go-lanty-client/internal/game"
)

var _ game.ProcessLauncher = (*ProcessLaunch)(nil)

type ProcessLaunch struct{}

func NewProcessLaunch() *ProcessLaunch {
	return &ProcessLaunch{}
}

func (launcher *ProcessLaunch) Launch(ctx context.Context, executable string, args []string, workingDir string) error {
	if filepath.Ext(executable) == ".bat" {
		args = append([]string{"/c", "start", "cmd.exe", "/k", executable}, args...)
		executable = "cmd.exe"
	}

	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Dir = workingDir

	if runtime.GOOS != "windows" {
		//TODO: Log this
		fmt.Println("Not on windows: Command that would run: ", cmd.String())
		return nil
	}

	return cmd.Start()
}
