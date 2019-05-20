package common

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

type ngrok struct {
	Tunnels string
}

// Ngrok util
var Ngrok = &ngrok{
	Tunnels: "http://localhost:4040/api/tunnels",
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (ng *ngrok) GetPublicURL() string {
	// Lấy thông tin ngrok
	res, err := http.Get(ng.Tunnels)
	logError(err)

	body, err := ioutil.ReadAll(res.Body)
	logError(err)

	data := string(body)
	publicURL := gjson.Get(data, "tunnels.0.public_url").String()

	return publicURL
}
