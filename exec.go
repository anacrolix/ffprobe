// +build !windows

package ffprobe

import (
	"os/exec"
)

func setHideWindow(cmd *exec.Cmd) {}
