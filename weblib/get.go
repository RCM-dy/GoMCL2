package weblib

import (
	"io"
	"net/http"
)

func GetBytesFromString(urls string) (rb Bytes, err error) {
	r, err := http.Get(urls)
	if err != nil {
		return
	}
	defer r.Body.Close()
	rb, err = io.ReadAll(r.Body)
	return
}
func GetBytesFromStringWithHeader(urls string, head map[string][]string) (rb Bytes, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, urls, nil)
	if err != nil {
		return
	}
	for k, v := range head {
		for _, v1 := range v {
			req.Header.Add(k, v1)
		}
	}
	r, err := client.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()
	rb, err = io.ReadAll(r.Body)
	return
}
