package utils

import (
	"encoding/base64"
	"github.com/toheart/functrace"
)

func Encode(s string) string {
	defer functrace.Trace([]interface {
	}{s})()
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func Decode(s string) (string, error) {
	defer functrace.Trace([]interface {
	}{s})()
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
