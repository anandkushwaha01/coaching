package common
import(
	"net/http"
	"encoding/json"
	"log"
	"net/url"
	"fmt"
)
type ResponseData map[string]interface{}

type Response struct {
	StatusCode         int 			//200 ok; 500 error
	StatusText         string
	ErrorStatusCode    int
	URL                string
	Output             ResponseData
	Headers            http.Header
	IsError            bool
	ErrorId            string
	InternalError      error
}
func NewResponse() *Response {
	r := &Response{
		StatusCode:      200,
		ErrorStatusCode: 0,
		Output:          make(ResponseData),
		Headers:         make(http.Header),
		IsError:         false,
	}
	r.Headers.Add(
		"Cache-Control",
		"no-cache, no-store, max-age=0, must-revalidate",
	)
	r.Headers.Add("Pragma", "no-cache")
	r.Headers.Add("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	return r
}

func (rs *Response) WriteJson(w http.ResponseWriter){
	// Add headers
	for i, k := range rs.Headers {
		for _, v := range k {
			w.Header().Add(i, v)
		}
	}
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(rs.StatusCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(rs.Output)
	if err != nil {
		log.Println("error in writing data. Error:", err)
		w.WriteHeader(500)
		w.Write([]byte("unknown error"))
	}
}

func (rs *Response) WriteError(w http.ResponseWriter){
	var msg string
	var err error
	if msg, err = GetErrorMessage(rs.ErrorStatusCode); err != nil {
		if rs.InternalError != nil{
			msg = rs.InternalError.Error()
		}else{
			msg = "unknown error"
		}
	}
	http.Error(w, msg, rs.ErrorStatusCode)
}
func (rs *Response) SetRedirectUrl(url string){
	rs.URL = url
}
func (rs *Response) RedirectUrl(w http.ResponseWriter){
	u, err := url.Parse(rs.URL)
	if err != nil {
		rs.ErrorStatusCode = ERROR_INVALID_REDIRECT_URI
		rs.WriteError(w)
		return
	}
	// add parameters
	q := u.Query()
	for n, v := range rs.Output {
		q.Set(n, fmt.Sprint(v))
	}
	u.RawQuery = q.Encode()
	w.Header().Add("Location", u.String())
	w.WriteHeader(302)
}