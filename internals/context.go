package server

type Context struct {
	Req *Request
	Res *Response
}

func NewContext(req *Request, res *Response) *Context {
	return &Context{Req: req, Res: res}
}
