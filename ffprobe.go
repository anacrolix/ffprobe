// Package ffprobe wraps and interprets ffmpeg's ffprobe for Go.
package ffprobe

import "errors"

var ExeNotFound = errors.New("ffprobe and avprobe not found in $PATH")

// Runs ffprobe or avprobe or similar on the given file path.
func Run(path string) (info *Info, err error) {
	pc, err := Start(path)
	if err != nil {
		return
	}
	<-pc.Done
	info, err = pc.Info, pc.Err
	return
}
