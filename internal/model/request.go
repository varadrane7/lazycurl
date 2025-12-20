package model

import "net/http"

// Request represents an HTTP request to be executed by curl.
type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// NewRequest creates a default request.
func NewRequest() Request {
	return Request{
		Method:  http.MethodGet,
		URL:     "https://httpbin.org/get",
		Headers: make(map[string]string),
		Body:    "",
	}
}
