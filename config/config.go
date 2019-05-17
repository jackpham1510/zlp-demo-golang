package config

import (
	"io/ioutil"
	"log"

	"github.com/tidwall/gjson"
)

var requiredConfigs = []string{
	"appid",
	"key1",
	"key2",
}

var _config = make(map[string]string)
var _json string

// Get config by given name
func Get(name string) string {
	return gjson.Get(_json, name).String()
}

func init() {
	jsonStr, err := ioutil.ReadFile("config.json")

	if err != nil {
		log.Fatal(err)
	}

	_json = string(jsonStr)
}
