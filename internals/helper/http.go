package helper

import "time"

type HTTPResponse struct {
	TimeStamp string      `json:"timestamp"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

func NewHTTPResponse(message string, data any) HTTPResponse {
	return HTTPResponse{
        TimeStamp: time.Now().Format(time.RFC3339),
        Message: message,
        Data:   data,
	}
}
