package xhttp

import (
	"bytes"
	"github.com/toheart/functrace"
	"io/ioutil"
	"net/http"
)

func sendRequest(method, url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{method, url, data})()
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token 5b59c8a43dfb287afba7df5ae145f41e9b26a537")
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Post(url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{url, data})()
	return sendRequest("POST", url, data)
}

func Put(url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{url, data})()
	return sendRequest("PUT", url, data)
}

func Delete(url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{url, data})()
	return sendRequest("DELETE", url, data)
}

func Head(url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{url, data})()
	return sendRequest("HEAD", url, data)
}

func Connect(url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{url, data})()
	return sendRequest("CONNECT", url, data)
}

func Options(url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{url, data})()
	return sendRequest("OPTIONS", url, data)
}

func Trace(url string, data []byte) (string, error) {
	defer functrace.Trace([]interface {
	}{url, data})()
	return sendRequest("TRACE", url, data)
}
