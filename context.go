package grouter

// dev file
// this file for testing router.
// just for test.
// finish an simple http engine for test Exeucter interface
import (
	//"io"
	"net/http"
	"sync"
)

// test Context
type Context struct {
	http_request *http.Request
	response     http.ResponseWriter
}

func (ctx *Context) SetRequest(value interface{}) bool {
	http_request, ok := value.(*http.Request)
	if ok {
		ctx.http_request = http_request
		return true
	}
	return false
}

func (ctx *Context) GetRequest(type_str string) interface{} {
	if type_str == "http" {
		return ctx.http_request
	}
	return nil
}

func (ctx *Context) SetResponseWriter(writer interface{}) {
	ctx.response = writer.(http.ResponseWriter)
}

func (ctx *Context) GetResponseWriter() interface{} {
	return ctx.response
}

func (ctx *Context) GetHttpResponse() http.ResponseWriter {
	return ctx.response.(http.ResponseWriter)
}

func (ctx *Context) Get(key string) interface{} {
	return nil
}

func (ctx *Context) GetAll() []Param {
	return nil
}

func (ctx *Context) SetParam(key string, value interface{}) {
}

func (ctx *Context) GetString(key string) string {
	return ""
}

func (ctx *Context) Reset() {
}

func (ctx *Context) Status(code int) {
	if ctx.response == nil {
		return
	}
	ctx.response.WriteHeader(code)
}

var _ Contexter = &Context{}

// -=================================================================================

type Engine struct {
	pool   sync.Pool
	router Router
}

func NewEngine() *Engine {
	eg := &Engine{
		router: NewRouter(),
	}
	eg.pool.New = func() interface{} {
		return eg.allocateContext()
	}
	return eg
}

func (eg *Engine) allocateContext() Contexter {
	return &Context{}
}

func (eg *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := eg.pool.Get().(*Context)
	ctx.SetRequest(req)
	ctx.SetResponseWriter(res)

	eg.router.Execute(ctx)
	ctx.Reset()
	eg.pool.Put(ctx)
}
