package utils

import (
	"encoding/json"
	"github.com/toheart/functrace"
	"io"
	"io/ioutil"
	"net/http"
)

func toJSON(r interface{}) ([]byte, error) {
	defer functrace.Trace([]interface {
	}{r})()
	val, err := json.Marshal(r)
	if err != nil {
		return []byte{}, err
	}
	return val, nil
}

func ParseJSON(reader io.Reader, val interface{}) error {
	defer functrace.Trace([]interface {
	}{reader, val})()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &val)
}

func WriteJSON(w http.ResponseWriter, status int, r interface{}) error {
	defer functrace.Trace([]interface {
	}{w, status, r})()
	val, err := toJSON(r)

	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, string(val))
	return nil
}
