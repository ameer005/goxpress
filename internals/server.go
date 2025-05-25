package server

import (
	"encoding/json"
	"fmt"
	"http-server/internals/httpmethod"
	"http-server/internals/request"
	"http-server/internals/response"
	"net"
)

type Server struct {
	addr   string
	Router *Router
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, Router: &Router{routes: make(map[httpmethod.Method][]RouteEntry)}}
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
	req, err := request.ParseReq(rawData)

	if err != nil {
		fmt.Println("failed to parse request ", err)
		return
	}

	res := response.NewResponse(con)

	ctx := NewContext(req, res)

	HandleRequest(ctx, t.Router)
}

func JSONBody[T any](r *request.Request) (T, error) {
	var data T

	err := json.Unmarshal(r.RawBody(), &data)

	return data, err
}
