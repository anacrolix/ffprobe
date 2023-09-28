package ffprobe

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Converts possible values in Info.
func AnyAsFloat64(a any) (float64, error) {
	switch v := a.(type) {
	case float64:
		return v, nil
	case json.Number:
		return v.Float64()
	default:
		return 0, fmt.Errorf("unhandled type %T", a)
	}
}

// Converts possible values in Info.
func AnyAsInt64(a any) (int64, error) {
	switch v := a.(type) {
	case int64:
		return v, nil
	case json.Number:
		return v.Int64()
	case float64:
		return strconv.ParseInt(fmt.Sprintf("%f", a), 0, 64)
	default:
		return 0, fmt.Errorf("unhandled type %T", a)
	}
}
