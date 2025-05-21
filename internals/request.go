package server

import "fmt"

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
	fmt.Printf("%v", rawData)

	return req, nil
}
