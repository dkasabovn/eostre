package operator

import "eostre/packet"

func convertToI64(argType packet.Type, val any) (intv int64, ok bool) {
	switch argType {
	case packet.Int64:
		intv, ok = val.(int64)
	case packet.Int32:
		v, o := val.(int32)
		intv = int64(v)
		ok = o
	case packet.Int:
		v, o := val.(int)
		intv = int64(v)
		ok = o
	case packet.Int16:
		v, o := val.(int16)
		intv = int64(v)
		ok = o
	case packet.Int8:
		v, o := val.(int8)
		intv = int64(v)
		ok = o
	default:
		intv = 0
		ok = false
	}
	return intv, ok
}

func convertToFloat64(argType packet.Type, val any) (floatv float64, ok bool) {
	switch argType {
	case packet.Float64:
		floatv, ok = val.(float64)
	case packet.Float32:
		v, o := val.(float32)
		floatv = float64(v)
		ok = o
	default:
		floatv = 0.0
		ok = false
	}
	return floatv, ok
}
