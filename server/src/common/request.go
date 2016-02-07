package common
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ParseJsonRequest(r *http.Request) (bodyJSON map[string]interface{}, req_error AppError) {
	debug := log.Println
	req_error = AppError{Code:200, Msg:"",}
	if r.Method != "POST" {
		debug("ParseJsonRequest: invalid method", r.Method)
		req_error.Code = ERROR_INVALID_METHOD_TYPE
		req_error.Msg  = "request must be post"
		return
	}
	cn_type := r.Header.Get("Content-Type")
	if !strings.Contains(cn_type, "application/json") {
		debug("ParseJsonRequest: invalid content type ", cn_type)
		req_error.Code = ERROR_INVALID_CONTENT_TYPE
		req_error.Msg  = "request Content-Type must be application/json type"
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		debug("ParseJsonRequest: body is not readable Error:", err)
		req_error.Code = ERROR_JSON_BODY_NOT_READABLE
		req_error.Msg  = "request body is not readable"
		return
	}
	bodyJSON = make(map[string]interface{}, 0)
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil{
		log.Println("ParseJsonRequest: error in parsing the request: ", err)
		req_error.Code = ERROR_JSON_BODY_NOT_PARSABLE
		req_error.Msg ="error in parsing the request body"
	}
	return
}
func ParseFormRequest(r *http.Request) (req_error AppError) {
	debug := log.Println
	req_error = AppError{Code:200, Msg:"",}
	if r.Method != "POST" {
		debug("ParseFormRequest: invalid method", r.Method)
		req_error.Code = ERROR_INVALID_METHOD_TYPE
		req_error.Msg  = "request must be post"
		return
	}
	cn_type := r.Header.Get("Content-Type")
	if !strings.Contains(cn_type, "application/x-www-form-urlencoded") {
		debug("ParseFormRequest: invalid content type ", cn_type)
		req_error.Code = ERROR_INVALID_CONTENT_TYPE
		req_error.Msg  = "request Content-Type must be application/x-www-form-urlencoded"
	}
	return
}
func ParseGetRequest(r *http.Request) (req_error AppError) {
	debug := log.Println
	req_error = AppError{Code:200, Msg:"",}
	if r.Method != "GET" {
		debug("ParseGetRequest: invalid method", r.Method)
		req_error.Code = ERROR_INVALID_METHOD_TYPE
		req_error.Msg  = "request must be get"
		return
	}
	return
}