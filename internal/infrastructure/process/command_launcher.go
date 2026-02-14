package process

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	gameappsrv "github.com/seternate/go-lanty-client/internal/application/game/service"
)

var _ gameappsrv.CommandLauncher = (*CommandLauncher)(nil)

type CommandLauncher struct {
}

func NewCommandLauncher() *CommandLauncher {
	return &CommandLauncher{}
}

func (launcher *CommandLauncher) Run(command gameappsrv.Command) error {
	if filepath.Ext(command.Executable) == ".bat" {
		command.Args = append([]string{"/c", "start", "cmd.exe", "/k", command.Executable}, command.Args...)
		command.Executable = "cmd.exe"
	}

	cmd := exec.Command(command.Executable, command.Args...)
	cmd.Dir = command.WorkingDir

	if runtime.GOOS == "windows" {
		return cmd.Start()
	}

	fmt.Println("Not on windows: Command that would run: ", cmd.String())

	return nil
}
