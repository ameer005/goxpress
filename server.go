package goxpress

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ameer005/goxpress/httpmethod"
)

type Server struct {
	addr   string
	Router *router
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, Router: &router{routes: make(map[httpmethod.Method][]routeEntry), globalMiddlewares: []HandlerFunc{}}}
}

func (t *Server) Listen() error {
	ls, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	defer ls.Close()

	fmt.Println("Listening on", t.addr)

	for {
		con, err := ls.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go t.handleConnection(con)
	}
}

func (t *Server) handleConnection(con net.Conn) {
	// closing connection after successfully handleing this request
	defer con.Close()

	/*initializing bytes slice to store request data*/
	var rawData = make([]byte, 1024)

	/*Storing raw bytes to rawData slice*/
	_, err := con.Read(rawData)
	if err != nil {
		fmt.Println("Invalid request", err)
	}

	/*Creating request and response */
	req, err := parseReq(rawData)

	if err != nil {
		fmt.Println("failed to parse request ", err)
		return
	}

	res := NewResponse(con)

	ctx := NewContext(req, res)

	HandleRequest(ctx, t.Router)
}

// Parsing json body
// because it requires generics and you can't add generic in method
func JSONBody[T any](r *Request) (T, error) {
	var data T

	err := json.Unmarshal(r.RawBody(), &data)

	return data, err
}

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

// server method for assigning global middlewares
func (t *Server) Use(middleware HandlerFunc) {
	t.Router.Use(middleware)
}
