package launcher

import (
	"os/exec"
	"path/filepath"
)

func RunDowser() error {
	return exec.Command(filepath.Join(CWD, "dowser.exe")).Start()
}
