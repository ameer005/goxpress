package server

type Context struct {
	req *Request
	res *Response
}

func NewContext(req *Request, res *Response) *Context {
	return &Context{req: req, res: res}
}
