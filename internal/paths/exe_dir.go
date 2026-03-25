package paths

import (
	"os"
	"path/filepath"
)

func ExecutableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return ".", err
	}
	if resolved, e := filepath.EvalSymlinks(exe); e == nil {
		exe = resolved
	}
	return filepath.Dir(exe), nil
}
