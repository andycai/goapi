package string

import "github.com/spf13/cast"

func defaultString(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func ToString(value string, defaultValue string) string {
	return defaultString(value, defaultValue)
}

func ToInt(value string, defaultValue int) int {
	if len(value) == 0 {
		return defaultValue
	}
	return cast.ToInt(value)
}

func ToUin(value string, defaultValue uint) uint {
	if len(value) == 0 {
		return defaultValue
	}
	return cast.ToUint(value)
}

func ToU32(value string, defaultValue uint32) uint32 {
	if len(value) == 0 {
		return defaultValue
	}
	return cast.ToUint32(value)
}

func ToI32(value string, defaultValue int32) int32 {
	if len(value) == 0 {
		return defaultValue
	}
	return cast.ToInt32(value)
}

func ToU64(value string, defaultValue uint64) uint64 {
	if len(value) == 0 {
		return defaultValue
	}
	return cast.ToUint64(value)
}

func ToI64(value string, defaultValue int64) int64 {
	if len(value) == 0 {
		return defaultValue
	}
	return cast.ToInt64(value)
}
