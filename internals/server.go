package server

import (
	"fmt"
	"net"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
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

	var rawData = make([]byte, 1024)
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
	fmt.Println(ctx.res.statusCode)
}
