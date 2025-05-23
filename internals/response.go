package server

import (
	"fmt"
	"http-server/internals/phrase"
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

func (t *Response) Status(code int) *Response {
	t.statusCode = code
	return t
}

func (t *Response) SetHeader(key, value string) *Response {
	t.headers[key] = value
	return t
}

func (t *Response) Send(body string) {
	line := statusLine("HTTP/1.1", t.statusCode)
	fmt.Println(line)

}

// Helpers
func statusLine(version string, code int) string {
	phrase, ok := phrase.StatusPhrases[code]
	if !ok {
		phrase = "Unkown Status"
	}

	return version + " " + fmt.Sprintf("%d", code) + " " + phrase + "\r\n"
}
