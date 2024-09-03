package cachex

func unmarshalString(buf string, val interface{}) bool {
	return unmarshal([]byte(buf), val)
}

func marshalString(val interface{}) string {
	res := marshal(val)
	return string(res)
}
