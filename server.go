package goxpress

import (
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
