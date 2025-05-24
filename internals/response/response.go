package response

import (
	"fmt"
	"http-server/internals/phrase"
	"net"
	"net/http"
	"strings"
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

func (t *Response) Send(message string) {
	// setting headers
	t.SetHeader("Content-Length", fmt.Sprintf("%d", len(message)))

	// building response text
	line := statusLine("HTTP/1.1", t.statusCode)
	headers := responseHeaders(t.headers)
	res := line + headers + message

	// sending response
	t.con.Write([]byte(res))
	return
}

// Helpers
func statusLine(version string, code int) string {
	phrase, ok := phrase.StatusPhrases[code]
	if !ok {
		phrase = "Unkown Status"
	}

	return version + " " + fmt.Sprintf("%d", code) + " " + phrase + "\r\n"
}

// method to build header in this format "key: value \r\n"
func responseHeaders(headers map[string]string) string {
	var headersString strings.Builder
	for k, v := range headers {
		headersString.WriteString(k)
		headersString.WriteString(": ")
		headersString.WriteString(v)
		headersString.WriteString("\r\n")
	}

	headersString.WriteString("\r\n") // End of headers

	return headersString.String()

}
