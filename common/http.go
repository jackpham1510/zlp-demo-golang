package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type _http struct{}

var Http = &_http{}

func (h *_http) MapToParams(m map[string]string) url.Values {
	params := make(url.Values)

	for key, value := range m {
		params.Add(key, value)
	}

	return params
}

func (h *_http) PostForm(url string, mapParams map[string]string) string {
	params := h.MapToParams(mapParams)
	res, err := http.PostForm(url, params)

	if err != nil {
		log.Println(err)
		return ""
	}

	body, _ := ioutil.ReadAll(res.Body)

	return string(body)
}
