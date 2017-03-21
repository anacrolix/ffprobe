package ffprobe

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyFile(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	_, err = Run(f.Name())
	assert.EqualError(t, err, fmt.Sprintf("exit status 1: %s: Invalid data found when processing input", f.Name()))
}
