package helpers

func Ternary[T interface{}](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
