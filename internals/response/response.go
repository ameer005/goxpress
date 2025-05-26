package response

import (
	"encoding/json"
	"fmt"
	"http-server/internals/phrase"
	"net"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	con             net.Conn
	statusCode      int
	body            []byte
	headers         map[string]string
	ResponseWritten bool
}

/*Implement body*/
func NewResponse(con net.Conn) *Response {
	return &Response{
		con:             con,
		statusCode:      404,
		ResponseWritten: false,
		headers: map[string]string{
			"Content-Type": "text/plain",
			"Connection":   "close",
			"Date":         time.Now().UTC().Format(http.TimeFormat),
		},
	}
}

// Client side method for Setting response status code
func (t *Response) Status(code int) *Response {
	t.statusCode = code
	return t
}

// Client side method for Setting headers status code
func (t *Response) SetHeader(key, value string) *Response {
	t.headers[key] = value
	return t
}

// Client side basic send reponse method
func (t *Response) Send(message string) {
	// setting headers
	t.SetHeader("Content-Length", fmt.Sprintf("%d", len(message)))

	// building response text
	line := statusLine("HTTP/1.1", t.statusCode)
	headers := responseHeaders(t.headers)
	res := line + headers + message

	// sending response
	t.con.Write([]byte(res))
	t.ResponseWritten = true
	return
}

// send json response
func (t *Response) JSON(payload map[string]any) {
	t.SetHeader("Content-Type", "application/json")

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		t.Status(500).Send("Internal Server Error")
	}

	t.SetHeader("Content-Length", fmt.Sprintf("%d", len(jsonData)))

	line := statusLine("HTTP/1.1", t.statusCode)
	headers := responseHeaders(t.headers)
	res := line + headers + string(jsonData)

	t.con.Write([]byte(res))
	t.ResponseWritten = true
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
