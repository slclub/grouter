package grouter

import (
	"bytes"
	"fmt"
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
	engine.router.GET("/ping/:uid", func(ctx Contexter) {
		ctx.Get("wo")
		ctx.GetAll()
		// dont implement to here so
		ctx.GetString("uid")
		//fmt.Println("/ping/:uid", "path:", ctx.GetRequest("http").(*http.Request).URL.Path)
		assert.Equal(t, "/ping/xiaoming", lowercase(ctx.GetRequest("http").(*http.Request).URL.Path))
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/Ping/xiaoming", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterHandleQuestion(t *testing.T) {
	r := NewRouter()
	var handle_index HandleFunc = func(ctx Contexter) {
		ctx.Status(http.StatusOK)
	}
	r.GET("/robot/list?a=1&b=2", handle_index)
	r.GET("/robot/book?a=1&b=2", handle_index)
	r.GET("/robot/test?a=1&b=2", handle_index)

}

func TestRouterQuestionHttp(t *testing.T) {
	engine := NewEngine()
	var handle_index HandleFunc = func(ctx Contexter) {
		ctx.Status(http.StatusOK)
	}
	engine.router.GET("/robot/list?a=1&b=2", handle_index)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/Ping/xiaoming", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterHandle(t *testing.T) {
	r := NewRouter()
	r.SetStore(NewStore())
	var handle_index HandleFunc = func(ctx Contexter) {
		ctx.Status(http.StatusOK)
	}
	r.GET("/index", handle_index)
	r.GET("/robot/list", handle_index)
	r.GET("/robot/get/:uid", handle_index)
	r.GET("/robot/update/:uid", handle_index)
	r.GET("/robot/not/:uid", handle_index)

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

	print_store_tree_root(t, r, http.MethodGet)
	print_store_tree_root(t, r, http.MethodPost)
	print_store_tree_root(t, r, http.MethodHead)
	print_store_tree_root(t, r, http.MethodOptions)
	print_store_tree_root(t, r, http.MethodPut)
	print_store_tree_root(t, r, http.MethodDelete)
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
	fmt.Println(get_copy_tree_depth([]byte(" "), depth*2), "|__", node.GetIndices(), "param:", node.GetKeys(), "nodeType:", node.GetType())
	children := node.GetChildren()
	if len(children) == 0 {
		return
	}
	for _, v := range children {
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
