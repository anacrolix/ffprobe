package ffprobe

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type Cmd struct {
	Cmd  *exec.Cmd
	Done chan struct{}
	mu   sync.Mutex
	Info *Info
	Err  error
}

func Start(path string) (ret *Cmd, err error) {
	if exePath == "" {
		err = ExeNotFound
		return
	}
	cmd := exec.Command(exePath,
		"-loglevel", "error",
		"-show_format",
		"-show_streams",
		outputFormatFlag, "json",
		path)
	setHideWindow(cmd)
	var stdout, stderr *io.PipeReader
	stdout, cmd.Stdout = io.Pipe()
	stderr, cmd.Stderr = io.Pipe()
	err = cmd.Start()
	if err != nil {
		return
	}
	ret = &Cmd{
		Cmd:  cmd,
		Done: make(chan struct{}),
	}
	go ret.runner(stdout, stderr)
	return
}

func (me *Cmd) runner(stdout, stderr io.ReadCloser) {
	defer close(me.Done)
	lastErrLineCh := lastLineCh(stderr)
	d := json.NewDecoder(bufio.NewReader(stdout))
	decodeErr := d.Decode(&me.Info)
	stdout.Close()
	waitErr := me.Cmd.Wait()
	stderr.Close()
	if waitErr == nil {
		me.Err = decodeErr
		return
	}
	lastErrLine, lastErrLineOk := <-lastErrLineCh
	if lastErrLineOk {
		me.Err = fmt.Errorf("%s: %s", waitErr, lastErrLine)
	} else {
		me.Err = waitErr
	}
	return
}

// Returns the last line in r. ok is false if there are no lines. err is any
// error that occurs during scanning.
func lastLine(r io.Reader) (line string, ok bool, err error) {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line = s.Text()
		ok = true
	}
	err = s.Err()
	return
}

// Returns a channel that receives the last line in r.
func lastLineCh(r io.Reader) <-chan string {
	ch := make(chan string, 1)
	go func() {
		defer close(ch)
		line, ok, err := lastLine(r)
		switch err {
		case nil:
		case io.ErrClosedPipe:
		default:
			panic(err)
		}
		if ok {
			ch <- line
		}
	}()
	return ch
}
