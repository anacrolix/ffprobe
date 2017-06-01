package ffprobe

import (
	"log"
	"os/exec"
)

var (
	exePath          string
	outputFormatFlag = "-of"
)

func exeFound() bool {
	return exePath != ""
}

func init() {
	var err error
	exePath, err = exec.LookPath("ffprobe")
	if err == nil {
		outputFormatFlag = "-print_format"
		return
	}
	if !isExecErrNotFound(err) {
		log.Print(err)
	}
	exePath, err = exec.LookPath("avprobe")
	if err == nil {
		return
	}
	if isExecErrNotFound(err) {
		err = ExeNotFound
	}
	log.Print(err)
}

func isExecErrNotFound(err error) bool {
	if err == exec.ErrNotFound {
		return true
	}
	execErr, ok := err.(*exec.Error)
	if !ok {
		return false
	}
	return execErr.Err == exec.ErrNotFound
}
