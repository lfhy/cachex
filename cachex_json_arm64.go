//go:build arm64
// +build arm64

package cachex

import "encoding/json"

func unmarshal(buf []byte, val interface{}) bool {
	return json.Unmarshal(buf, val) == nil
}

func marshal(val interface{}) []byte {
	res, _ := json.Marshal(val)
	return res
}
