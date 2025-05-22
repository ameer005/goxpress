package server

import (
	"fmt"
	"net"
)

type Server struct {
	addr   string
	Router *Router
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, Router: &Router{}}
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
	/*initializing bytes slice to store request data*/
	var rawData = make([]byte, 1024)

	/*Storing raw bytes to rawData slice*/
	_, err := con.Read(rawData)
	if err != nil {
		fmt.Println("Invalid request", err)
	}

	/*Creating request and response */
	req, err := ParseReq(rawData)

	if err != nil {
		fmt.Println("failed to parse request ", err)
		return
	}

	res := NewResponse(con)

	ctx := NewContext(req, res)
	HandleRequest(ctx, t.Router)
}
