package utils

import (
	"github.com/toheart/functrace"
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func RoomIdCreate(n int) string {
	defer functrace.Trace([]interface {
	}{n})()
	id := ""
	for i := 0; i < n; i++ {
		id += string(chars[rand.Intn(len(chars))])
	}
	return id
}
