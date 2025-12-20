package model

import "time"

// Response represents the result of a curl execution.
type Response struct {
	StatusCode int               `json:"status_code"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	TimeTaken  time.Duration     `json:"time_taken"`
	Error      error             `json:"error,omitempty"`
}
