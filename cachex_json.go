//go:build !arm64
// +build !arm64

package cachex

import "github.com/bytedance/sonic"

func unmarshal(buf []byte, val interface{}) bool {
	return sonic.Unmarshal(buf, val) == nil
}

func marshal(val interface{}) []byte {
	res, _ := sonic.Marshal(val)
	return res
}
