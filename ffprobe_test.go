package ffprobe

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/missinggo/leaktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyFile(t *testing.T) {
	if !exeFound() {
		t.SkipNow()
	}
	f, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	_, err = Run(f.Name())
	assert.EqualError(t, err, fmt.Sprintf("exit status 1: %s: Invalid data found when processing input", f.Name()))
}

func TestKilledWhileStuckReading(t *testing.T) {
	if !exeFound() {
		t.SkipNow()
	}
	time.Sleep(time.Second)
	defer leaktest.GoroutineLeakCheck(t)()
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
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
	require.NoError(t, err)
	require.NoError(t, cmd.Cmd.Process.Kill())
	s.Close()
	// time.Sleep(time.Second)
	// select {}
}
