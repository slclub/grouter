package grouter

// dev file
// this file for testing router.
// just for test.
// finish an simple http engine for test Exeucter interface
import (
	//"io"
	"github.com/slclub/gnet"
	"net/http"
	"sync"
)

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

func (eg *Engine) allocateContext() gnet.Contexter {
	ctx := gnet.NewContext()
	r := gnet.NewRequest()
	s := &gnet.Response{}
	ctx.SetRequest(r)
	ctx.SetResponse(s)

	return ctx
}

func (eg *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := eg.pool.Get().(gnet.Contexter)
	ctx.Request().InitWithHttp(req)
	ctx.Response().InitSelf(res)

	eg.router.Execute(ctx)
	ctx.Reset()
	eg.pool.Put(ctx)
}
