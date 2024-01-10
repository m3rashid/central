package helpers

import "golang.org/x/exp/constraints"

func Ternary[T interface{}](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}

func RemoveDuplicates[T constraints.Ordered](arr []T) []T {
	var newArr []T
	numFound := false

	for _, num := range arr {
		numFound = false
		for _, newNum := range newArr {
			if num == newNum {
				numFound = true
			}
		}

		if !numFound {
			newArr = append(newArr, num)
		}
	}

	return newArr
}
