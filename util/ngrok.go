package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"zlp-demo-golang/config"

	"github.com/tidwall/gjson"
)

type ngrok struct{}

// Ngrok util
var Ngrok = &ngrok{}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (ng *ngrok) GetPublicURL() string {
	// Lấy thông tin ngrok
	res, err := http.Get(config.Get("ngrok.tunnels"))
	logError(err)

	body, err := ioutil.ReadAll(res.Body)
	logError(err)

	data := string(body)
	publicURL := gjson.Get(data, "tunnels.0.public_url").String()

	return publicURL
}
