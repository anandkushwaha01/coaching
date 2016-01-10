package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ParseBody(r *http.Request) (bodyJSON map[string]interface{}, req_error error) {
	debug := log.Println
	if r.Method == "POST" {
		req_error = errors.New("request must be POST")
		return
	}
	cn_type := r.Header.Get("Content-Type")
	if !strings.Contains(cn_type, "application/json") {
		debug("Onboard: invalid content type ", cn_type)
		req_error = errors.New("Not application/json type")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		debug("onboard: body is not readable")
		req_error = errors.New(" Not a valid type of form")
		return
	}
	bodyJSON = make(map[string]interface{}, 0)
	req_error = json.Unmarshal(body, &bodyJSON)
	return
}
