package utils

import "github.com/toheart/functrace"

func Contains[T comparable](s []T, elem T) bool {
	defer functrace.Trace([]interface {
	}{s, elem})()
	for _, a := range s {
		if a == elem {
			return true
		}
	}
	return false
}

func RemoveAll[T comparable](s []T, elem T) []T {
	defer functrace.Trace([]interface {
	}{s, elem})()
	result := make([]T, 0, len(s))
	for _, v := range s {
		if v != elem {
			result = append(result, v)
		}
	}
	return result
}
