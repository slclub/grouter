package grouter

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestRouterHandle(t *testing.T) {
	r := NewRouter()
	r.SetStore(NewStore())
	var handle_index HandleFunc = func(ctx Context) {
	}
	r.GET("/index", handle_index)
	r.GET("/robot/list", handle_index)
	r.GET("/robot/get/:uid", handle_index)
	r.GET("/robot/update/:uid", handle_index)
	r.GET("/robot/not/:uid", handle_index)

	print_store_tree_root(t, r, http.MethodGet)
	r.POST("/post", handle_index)
	r.POST("/robot/list", handle_index)
	r.POST("/robot/get/:uid", handle_index)
	r.POST("/robot/update/:uid", handle_index)
	r.POST("/robot/not/:uid", handle_index)

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
