package request

import (
	"errors"
	"strings"
)

// TODO handle body
type Request struct {
	method  string
	path    string
	proto   string
	headers map[string]string
}

func ParseReq(rawData []byte) (*Request, error) {
	req := &Request{
		headers: make(map[string]string),
	}
	lines := strings.Split(string(rawData), "\r\n")

	requesLine := strings.Split(lines[0], " ")
	if len(requesLine) < 3 {
		return nil, errors.New("Invalid request")
	}

	req.method = requesLine[0]
	req.path = requesLine[1]
	req.proto = requesLine[2]

	i := 1
	for ; i < len(lines); i++ {
		currLine := lines[i]
		if currLine == "" {
			i++
			break
		}

		values := strings.Split(currLine, ":")
		req.headers[values[0]] = strings.TrimRight(values[1], " ")
	}

	/* Parse body by using index i	 */

	return req, nil
}

// Getters
func (t *Request) Headers(key string) string {
	return t.headers[key]
}

func (t *Request) RequestMethod() string {
	return t.method
}

func (t *Request) RequestPath() string {
	return t.path
}
