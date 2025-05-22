package server

import (
	"net"
	"net/http"
	"time"
)

type Response struct {
	con        net.Conn
	statusCode int
	body       []byte
	headers    map[string]string
}

/*Implement body*/
func NewResponse(con net.Conn) *Response {
	return &Response{
		con:        con,
		statusCode: 404,
		headers: map[string]string{
			"Content-Type": "text/plain",
			"Connection":   "close",
			"Date":         time.Now().UTC().Format(http.TimeFormat),
		},
	}
}
