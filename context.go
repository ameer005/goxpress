package goxpress

type Context struct {
	Req  *Request
	Res  *Response
	Data map[string]any
}

func NewContext(req *Request, res *Response) *Context {
	return &Context{Req: req, Res: res, Data: make(map[string]any)}
}
