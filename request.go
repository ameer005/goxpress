package goxpress

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type FormFile struct {
	Name     string
	Path     string
	Size     int64
	Mimetype string
}

// TODO handle body
type Request struct {
	method  string
	path    string
	proto   string
	headers map[string]string
	body    []byte
	params  map[string]string
	query   map[string]string
	con     net.Conn
	files   map[string]FormFile
}

func parseReq(rawData []byte, con net.Conn) (*Request, error) {
	req := &Request{
		headers: make(map[string]string),
		query:   make(map[string]string),
		con:     con,
		files:   make(map[string]FormFile),
	}

	// separting body and header without converting it to string
	parts := bytes.SplitN(rawData, []byte("\r\n\r\n"), 2)

	if len(parts) == 0 {
		return nil, errors.New("Malformed Request")
	}

	// converting header to string
	headerPart := strings.Split(string(parts[0]), "\r\n")

	// parsing request line
	// method / uri / proto
	requesLine := strings.Split(headerPart[0], " ")
	if len(requesLine) < 3 {
		return nil, errors.New("Malformed request line")
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
		req.headers[values[0]] = strings.TrimLeft(strings.TrimRight(values[1], " "), " ")
	}

	// only parsing body if it exist
	if len(parts) > 1 {

		// assiging body data
		bodyPart := parts[1]

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
func (t *Request) GetReqPath() string {

	return t.path
}

// Parse URL encoded form data

// For getting typesafe query
func QueryData[T any](r *Request) (T, error) {
	var data T

	queryMap := r.UntypedQuery() // map[string]string

	// convert map to json
	jsonBytes, err := json.Marshal(queryMap)
	if err != nil {
		return data, err
	}

	// convert json to struct
	err = json.Unmarshal(jsonBytes, &data)
	return data, err
}

// internal method for parsing multipart data it takes sizeLimit (byes)
// TODO add size limit
func (t *Request) parseMultipart() error {

	if !strings.Contains(t.headers["Content-Type"], "multipart/form-data") {
		return errors.New("Incorrect content type")
	}

	contentType := strings.SplitN(t.headers["Content-Type"], "=", 2)
	boundary := strings.TrimSpace(contentType[1])

	// method to send client permission to send data
	if t.headers["Expect"] == "100-continue" {
		status := statusLine(t.proto, 100)
		t.con.Write([]byte(status))
	}

	mr := multipart.NewReader(bufio.NewReader(t.con), boundary)

	var formFields map[string]string = make(map[string]string)

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() != "" {
			// handling file upload
			dst, err := os.Create("/tmp/" + part.FileName())
			if err != nil {
				log.Println("Error creating file:", err)
				continue
			}
			defer dst.Close()

			_, err = io.Copy(dst, part)
			if err != nil {
				log.Println("Error saving file:", err)
				continue
			}

			dstFile, err := os.Open(dst.Name())
			if err != nil {
				log.Println("Error reopening file:", err)
				continue
			}
			defer dstFile.Close()

			fileInfo, err := os.Stat(dst.Name())
			if err != nil {
				log.Println("File does not exist:", err)
			}

			// Read first 512 bytes for mimetype detection
			buffer := make([]byte, 512)
			n, err := dstFile.Read(buffer)
			if err != nil && err != io.EOF {
				log.Println("Error reading file for MIME type:", err)
				continue
			}

			metadata := FormFile{
				Name:     part.FileName(),
				Path:     dst.Name(),
				Size:     fileInfo.Size(),
				Mimetype: http.DetectContentType(buffer[:n]),
			}

			t.files[part.FormName()] = metadata

		} else {
			// parsing form fields
			if part.FormName() != "" {
				data, err := io.ReadAll(part)

				if err != nil {
					log.Println("Error reading form field:", err)
					continue
				}

				formFields[part.FormName()] = string(data)
			}

		}

	}

	// converting form fields to bytes then storing it to body
	formBytes, err := json.Marshal(formFields)
	if err != nil {
		log.Println("Error marshaling form fields:", err)
	} else {
		t.body = formBytes
	}

	return nil
}

// client side parsing functions
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

// Parsing json body
// because it requires generics and you can't add generic in method
func JSONBody[T any](r *Request) (T, error) {
	var data T

	err := json.Unmarshal(r.RawBody(), &data)

	return data, err
}

// client side fucntion to get file metadata
func (t *Request) FormFile(filename string) (*FormFile, error) {
	if len(t.body) == 0 && len(t.files) == 0 {
		t.parseMultipart()
	}

	metadata, ok := t.files[filename]

	if !ok {
		return &FormFile{}, errors.New("no file found")
	}

	return &metadata, nil
}

// client side function to get form data
func FormBody[T any](r *Request) (T, error) {
	if len(r.body) == 0 && len(r.files) == 0 {
		r.parseMultipart()
	}

	var data T

	err := json.Unmarshal(r.body, &data)

	return data, err
}
