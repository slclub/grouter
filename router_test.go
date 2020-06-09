package grouter

import (
	"bytes"
	"fmt"
	"github.com/slclub/gnet"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCodeFunctions(t *testing.T) {
	r := NewRouter()
	r.SetDecoder(NewPath())

	h := r.CodeHandle(http.StatusNotFound)
	assert.Equal(t, 404, http.StatusNotFound)
	assert.NotNil(t, h)
	if h == nil {
		fmt.Println("404 handle  not exist")
	}

	h = r.CodeHandle(405)
	assert.NotNil(t, h)
	h = r.CodeHandle(500)
	assert.NotNil(t, h)

}

// testing http request
func TestHttpRouterListen(t *testing.T) {
	engine := NewEngine()
	// and support URL path captial leter.
	engine.router.GET("/test/ping/:uid/:sex", func(ctx gnet.Contexter) {
		ctx.Get("wo")
		// dont implement to here so
		uid, _ := ctx.Request().GetString("uid")
		sex, _ := ctx.Request().GetString("sex")
		assert.Equal(t, "xiaoming", uid)
		fmt.Println("/ping/:uid/:sex", "path:", uid, sex)
		//assert.Equal(t, "/ping/xiaoming", lowercase(ctx.Request().GetHttpRequest().URL.Path))
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/Ping/xiaoming/boy", nil)
	engine.ServeHTTP(w, req)
	req, _ = http.NewRequest("GET", "/test/ping/xiaoming/girl", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterHandleQuestion(t *testing.T) {
	r := NewRouter()
	var handle_index gnet.HandleFunc = func(ctx gnet.Contexter) {
		ctx.Response().WriteHeader(http.StatusOK)
	}
	r.GET("/robot/list?a=1&b=2", handle_index)
	r.GET("/robot/book?a=1&b=2", handle_index)
	r.GET("/robot/test?a=1&b=2", handle_index)

}

func TestRouterQuestionHttp(t *testing.T) {
	engine := NewEngine()
	var handle_index gnet.HandleFunc = func(ctx gnet.Contexter) {
		ctx.Response().WriteHeader(http.StatusOK)
	}
	engine.router.GET("/robot/list", handle_index)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/robot/list?a=1&b=2", nil)
	///news/analysis?name=%E6%B5%8B%E8%AF%95%E6%95%B0%E6%8D%AE
	engine.ServeHTTP(w, req)

	print_store_tree_root(t, engine.router, http.MethodGet)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestRouterQuestionEncodeUrl(t *testing.T) {
	engine := NewEngine()
	var handle_index gnet.HandleFunc = func(ctx gnet.Contexter) {
		ctx.Response().WriteHeader(http.StatusOK)
	}
	engine.router.GET("/news/analysis", handle_index)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news/analysis?name=%E6%B5%8B%E8%AF%95%E6%95%B0%E6%8D%AE", nil)
	///news/analysis?name=%E6%B5%8B%E8%AF%95%E6%95%B0%E6%8D%AE
	engine.ServeHTTP(w, req)

	print_store_tree_root(t, engine.router, http.MethodGet)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter404(t *testing.T) {
	engine := NewEngine()
	var handle_list gnet.HandleFunc = func(ctx gnet.Contexter) {
		ctx.Response().WriteHeader(404)
		ctx.Response().Flush()
	}
	engine.router.GET("/robot/list?a=1&b=2", handle_list)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/Ping/notfound", nil)
	engine.ServeHTTP(w, req)

	// because ResponseWriter was overide by Context.Response.
	// assert.Equal(t, 404, w.Code)

	// test not found static file
	engine.router.ServerFile("/st/", http.Dir("."), true)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/st/m", nil)
	engine.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/", nil)
	engine.ServeHTTP(w, req)

}

func TestRouterHandle(t *testing.T) {
	r := NewRouter()
	r.SetStore(NewStore())
	var handle_index gnet.HandleFunc = func(ctx gnet.Contexter) {
		ctx.Response().WriteHeader(http.StatusOK)
	}
	r.GET("/index", handle_index)
	r.GET("/robot/list", handle_index)
	r.GET("/robot/get/:uid", handle_index)
	r.GET("/robot/update/:uid", handle_index)
	r.GET("/robot/not/:uid", handle_index)
	r.ServerFile("/st/", http.Dir("."))

	r.POST("/post", handle_index)
	r.POST("/robot/list", handle_index)
	r.POST("/robot/get/:uid", handle_index)
	r.POST("/robot/update/:uid", handle_index)
	r.POST("/robot/not/:uid", handle_index)

	r.PUT("/post", handle_index)
	r.PUT("/robot/list", handle_index)
	r.PUT("/robot/get/:uid", handle_index)
	r.PUT("/robot/update/:uid", handle_index)
	r.PUT("/robot/not/:uid", handle_index)

	r.DELETE("/robot/not/:uid", handle_index)
	r.HEAD("/robot/not/:uid", handle_index)
	r.OPTIONS("/robot/not/:uid", handle_index)
	r.ANY("/robot/not/:uid", handle_index)

	print_store_tree_root(t, r, http.MethodGet)
	print_store_tree_root(t, r, "ANY")
	//print_store_tree_root(t, r, http.MethodHead)
	//print_store_tree_root(t, r, http.MethodOptions)
	//print_store_tree_root(t, r, http.MethodPut)
	//print_store_tree_root(t, r, http.MethodDelete)
}

func TestRouterErrorANDPanic(t *testing.T) {
	r := NewRouter()
	var handle_index gnet.HandleFunc = func(ctx gnet.Contexter) {
		ctx.Response().WriteHeader(http.StatusOK)
	}
	assert.Panics(t, func() { r.GET("", nil) })
	assert.Panics(t, func() { r.GET("/ok", nil) })
	assert.Panics(t, func() { r.Handle("", "/ok", nil) })
	assert.Panics(t, func() {
		r.Handle("FK", "/ok", handle_index)
	})
	r.BindCodeHandle(400, nil)
}

func print_store_tree_root(t *testing.T, r Router, method string) {

	var store Store = r.GetStore()
	fmt.Println("---", method, "---------------------------------------------------------")
	var node Node
	node, _ = store.Lookup(method)
	print_store_tree_node(t, node, 1)
}

func print_store_tree_node(t *testing.T, node Node, depth int) {
	if node == nil {
		return
	}
	fmt.Println(get_copy_tree_depth([]byte(" "), depth*2), "|__", node.GetIndices(), "param:", node.GetKeys(), "nodeType:", node.GetType(), "path", node.GetPath())
	children := node.GetChildren()
	if len(children) == 0 {
		return
	}
	for _, v := range children {
		//fmt.Println("children key -----", key)
		print_store_tree_node(t, v, depth+1)
	}
}

func get_copy_tree_depth(s []byte, depth int) string {
	var buf bytes.Buffer
	for i := 0; i < depth; i++ {
		buf.Write(s)
	}
	return string(buf.Bytes())
}
