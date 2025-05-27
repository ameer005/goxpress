package goxpress

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// TODO handle body
type Request struct {
	method  string
	path    string
	proto   string
	headers map[string]string
	body    []byte
	params  map[string]string
	query   map[string]string
}

func parseReq(rawData []byte) (*Request, error) {
	req := &Request{
		headers: make(map[string]string),
		query:   make(map[string]string),
	}

	// separting body and header without converting it to string
	parts := bytes.SplitN(rawData, []byte("\r\n\r\n"), 2)

	if len(parts) < 2 {
		return nil, errors.New("Invalid request")
	}

	// converting header to string
	headerPart := strings.Split(string(parts[0]), "\r\n")
	bodyPart := parts[1]

	// parsing request line
	// method / uri / proto
	requesLine := strings.Split(headerPart[0], " ")
	if len(requesLine) < 3 {
		return nil, errors.New("Invalid request")
	}
	req.method = requesLine[0]
	req.path = requesLine[1]
	req.proto = requesLine[2]

	// parsing headers
	i := 1
	for ; i < len(headerPart); i++ {
		currLine := headerPart[i]
		if currLine == "" {
			i++
			break
		}

		values := strings.SplitN(currLine, ":", 2)
		if len(values) < 2 {
			return nil, errors.New("Invalid Headers")
		}
		req.headers[values[0]] = strings.TrimRight(values[1], " ")
	}

	// assiging body data
	req.body = bytes.TrimRight(bodyPart, "\x00 \n\r\t")
	req.params = make(map[string]string)

	// parsing path and queries
	pathParts := strings.SplitN(req.path, "?", 2)
	if len(pathParts) == 2 {
		req.path = pathParts[0]
		rawQuery := pathParts[1]
		pairs := strings.Split(rawQuery, "&")

		for _, query := range pairs {
			queryPairs := strings.SplitN(query, "=", 2)
			if len(queryPairs) < 2 {
				continue
			}
			req.query[queryPairs[0]] = queryPairs[1]
		}

	}

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

func (t *Request) RawBody() []byte {
	return t.body
}

func (t *Request) UntypedQuery() map[string]string {
	return t.query
}

func (t *Request) SetRequestParam(key, value string) {
	t.params[key] = value
}

func (t *Request) GetParams() map[string]string {

	return t.params
}

// Parse URL encoded form data
func (t *Request) ParseURLEncodedForm() (map[string]string, error) {
	formData := make(map[string]string)

	data := strings.Split(string(t.body), "&")

	for _, pair := range data {
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		key, err := url.QueryUnescape(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid key encoding: %w", err)
		}

		var value string
		if len(parts) > 1 {
			value, err = url.QueryUnescape(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid value encoding: %w", err)
			}
		}

		formData[key] = value
	}

	return formData, nil
}
