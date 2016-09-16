package main

func Clamp(value, min, max int32) int32 {
	switch {
	case value < min:
		return min
	case value > max:
		return max
	}
	return value
}
