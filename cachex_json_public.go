package cachex

func unmarshalString(buf string, val interface{}) bool {
	return unmarshal([]byte(buf), val)
}

func marshalString(val interface{}) string {
	res := marshal(val)
	return string(res)
}

func UnmarshalString(buf string, val interface{}) bool {
	return unmarshalString(buf, val)
}

func MarshalString(val interface{}) (string, bool) {
	res := marshalString(val)
	return res, res != ""
}

func Unmarshal(buf []byte, val interface{}) bool {
	return unmarshal(buf, val)
}

func Marshal(val interface{}) ([]byte, bool) {
	res := marshal(val)
	return res, res != nil
}
