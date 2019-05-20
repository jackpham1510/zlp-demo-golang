package util

import (
	"encoding/json"
	"io"
	"log"
)

type _json struct{}

// JSON utils
var JSON _json

// ParseReader parse io.ReadCloser to json
// ex: request.Body
func (this *_json) ParseReader(body io.ReadCloser) map[string]string {
	defer body.Close()

	var result map[string]string

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&result); err != nil {
		log.Println("[util.JSON.ParseReader]", err)
	}

	return result
}

func (this *_json) Parse(jsonStr string) map[string]interface{} {
	var result map[string]interface{}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Println("[util.JSON.parse]", err)
		return nil
	}

	return result
}

func (this *_json) Add(jsonStr, key, val string) string {
	result := this.Parse(jsonStr)

	if result == nil {
		return jsonStr
	}

	result[key] = val
	resultBytes, _ := json.Marshal(result)
	return string(resultBytes)
}
