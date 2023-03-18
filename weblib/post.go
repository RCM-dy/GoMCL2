package weblib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func PostMapGotByteInStr(urls string, v map[string]any) (Bytes, error) {
	client := &http.Client{}
	value, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urls, bytes.NewReader(value))
	if err != nil {
		return nil, err
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
func PostMapGotByteInStrWithHeaders(urls string, v map[string]any, headers map[string][]string) (Bytes, error) {
	client := &http.Client{}
	value, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urls, bytes.NewReader(value))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		for _, v1 := range v {
			req.Header.Add(k, v1)
		}
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
func PostMapGotByteInStrWithHeader(urls string, v map[string]any, header map[string]string) (Bytes, error) {
	client := &http.Client{}
	value, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urls, bytes.NewReader(value))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
func PostMapGotStrInStrWithHeader(urls string, v map[string]any, header map[string]string) (string, error) {
	client := &http.Client{}
	value, err := json.Marshal(&v)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, urls, bytes.NewReader(value))
	if err != nil {
		return "", err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	r, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	rb, err := io.ReadAll(r.Body)
	return string(rb), err
}
