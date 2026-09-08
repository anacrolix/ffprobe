package ffprobe

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/missinggo/leaktest"
	"github.com/go-quicktest/qt"
)

func TestEmptyFile(t *testing.T) {
	if !exeFound() {
		t.SkipNow()
	}
	f, err := ioutil.TempFile("", "")
	qt.Assert(t, qt.IsNil(err))
	defer os.Remove(f.Name())
	_, err = Run(f.Name())
	qt.Check(t, qt.ErrorMatches(err, regexp.QuoteMeta(
		fmt.Sprintf("exit status 1: %s: Invalid data found when processing input", f.Name()))))
}

func TestKilledWhileStuckReading(t *testing.T) {
	if !exeFound() {
		t.SkipNow()
	}
	time.Sleep(time.Second)
	defer leaktest.GoroutineLeakCheck(t)()
	l, err := net.Listen("tcp", "localhost:0")
	qt.Assert(t, qt.IsNil(err))
	s := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print("got request")
			<-r.Context().Done()
		}),
	}
	go func() {
		log.Printf("serve returned: %s", s.Serve(l))
	}()
	defer s.Close()
	cmd, err := Start("http://" + l.Addr().String())
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.IsNil(cmd.Cmd.Process.Kill()))
	s.Close()
	// time.Sleep(time.Second)
	// select {}
}
