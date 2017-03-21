package ffprobe

import (
	"errors"
	"fmt"
	"time"
)

type Info struct {
	Format  map[string]interface{}
	Streams []map[string]interface{}
}

// returns res attributes for the raw stream
func (info *Info) Bitrate() (bitrate uint, err error) {
	bit_rate, exist := info.Format["bit_rate"]
	if !exist {
		err = errors.New("no bit_rate key in format")
		return
	}
	_, err = fmt.Sscan(bit_rate.(string), &bitrate)
	return
}

func (info *Info) Duration() (duration time.Duration, err error) {
	di := info.Format["duration"]
	if di == nil {
		err = errors.New("missing value")
		return
	}
	ds := di.(string)
	if ds == "N/A" {
		err = errors.New("N/A")
		return
	}
	var f float64
	_, err = fmt.Sscan(ds, &f)
	if err != nil {
		return
	}
	duration = time.Duration(f * float64(time.Second))
	return
}
