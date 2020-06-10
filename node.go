package grouter

import (
	"bytes"
	//"fmt"
	"github.com/slclub/gnet"
	"net/http"
	"strings"
)

const (
	// method node . as every method root node.
	NODE_T_ROOT = uint8(1)
	// wild node. have no path matched.
	NODE_T_WILD = uint8(2)

	// valid path node.
	NODE_T_PATH = uint8(3)
)

// the place of store the nodes.
// TODO: auto node merge and split.for good performence.
type nodeStore struct {
	nodes        map[string]Node
	http_methods map[string]bool
}

// Node defined need impement these interface function.
//type Node interface {
//	GetPath() string
//	SetPath(string)
//
//	GetHandleFunc() HandleFunc
//	SetHandleFunc(HandleFunc)
//
//	GetType() uint8
//	SetType() uint8
//
//	GetIndices() string
//	SetIndices(string)
//
//	// Get node from children.
//	// Overide one of them will be ok.
//	GetNodeAuto(interface{}) Node
//	GetNodeInt(int) Node
//	GetNodeStr(string) Node
//	//GetNodes() map[string]Node
//
//	// Add a new node to the childrens.
//	AddNode(Node)
//}
//
type node struct {
	path     string
	n_type   uint8
	handle   gnet.HandleFunc
	indices  string
	children map[string]Node
	keys     []string
}

// ============================store router node function =====================
func NewStore() Store {
	st := &nodeStore{
		http_methods: map[string]bool{},
	}

	st.http_methods[http.MethodGet] = true
	st.http_methods[http.MethodHead] = true
	st.http_methods[http.MethodOptions] = true
	st.http_methods[http.MethodPost] = true
	st.http_methods[http.MethodPut] = true
	st.http_methods[http.MethodPatch] = true
	st.http_methods[http.MethodDelete] = true
	st.http_methods["ANY"] = true

	return st
}
func (st *nodeStore) Save(n Node) {
	if st.nodes == nil {
		st.nodes = make(map[string]Node)
	}
	st.nodes[n.GetIndices()] = n
}

func (st *nodeStore) Create() Store {
	return NewStore()
}

func (st *nodeStore) Lookup(method string) (n Node, nothing string) {
	n = st.nodes[method]
	// TODO:need to support more schema. not only http
	if !st.checkHttpMethod(method) {
		nothing = "NOT ALLOWED METHOD"
		return
	}
	if n == nil {
		n = st.CreateNode(0)
		n.SetIndices(method)
		n.SetType(NODE_T_ROOT)
		st.Save(n)
	}
	return
}

func (st *nodeStore) CreateNode(param_num int) Node {

	return &node{
		children: make(map[string]Node),
		keys:     make([]string, param_num),
	}
}

// http allowd methods check
func (st *nodeStore) checkHttpMethod(method string) bool {
	if st.http_methods[method] {
		return true
	}
	return false
}

// ============================store router node function =====================

func (nd *node) GetPath() string {
	return nd.path
}

func (nd *node) SetPath(path string) {
	nd.path = path
}

func (nd *node) GetType() uint8 {
	return nd.n_type
}

func (nd *node) SetType(n_type uint8) {
	nd.n_type = n_type
}

func (nd *node) SetHandleFunc(handle gnet.HandleFunc) {
	nd.handle = handle
}

func (nd *node) GetHandleFunc() gnet.HandleFunc {
	return nd.handle
}

func (nd *node) GetIndices() string {
	return nd.indices
}

func (nd *node) SetIndices(indices string) {
	nd.indices = (indices)
}

func (nd *node) GetNodeAuto(key interface{}) Node {
	if rel, ok := key.(int); ok {
		return nd.GetNodeInt(rel)
	}

	if rel, ok := key.(string); ok {
		return nd.GetNodeStr(rel)
	}
	return nil
}

func (nd *node) GetNodeInt(index int) Node {
	return nil
}

func (nd *node) GetNodeStr(indices string) Node {
	return nd.children[indices]
}

func (nd *node) AddNode(one Node) bool {
	indices := lowercase(one.GetIndices())

	//fmt.Println("--------AddNode------", indices, one.GetIndices())
	if nd.children[indices] != nil {
		return false
	}
	nd.children[indices] = one
	return true
}

func (nd *node) AddKey(keys []string) {
	//for _, k := range nd.keys {
	//	if key == k {
	//		return
	//	}
	//}
	//nd.keys = append(nd.keys, key)
	nd.keys = keys
}

func (nd *node) GetKeys() []string {
	return nd.keys
}

// router query.
// Match the longest chain
func (nd *node) Lookup(path string) (Node, string) {
	var begin = 0

	var next Node
	next = nd
	head := nd
	len_path := len(path)

	// len_path==1 path=/
	if len_path == 1 {
		next = head.GetNodeAuto(lowercase(path[begin:1]))
		if next == nil {
			return nd, path
		}
		return next, ""
	}

	// fmt.Println("--------Lookup------", path)
	// Program optimization lowercase.
	// avoid applying for buffer in the loop
	// 这里可以优化掉 3 alloc/op  尽可能的避免在循环里申请buffer 代替lowercase
	buf := bytes.NewBufferString(path)
	buf.Reset()
	for i := 0; i < len_path; i++ {
		tmp_b := uint8(path[i])
		lowercaseWithBuffer(buf, tmp_b)
		if path[i] != '/' || i == 0 {
			continue
		}

		//next = head.GetNodeAuto(lowercase(path[begin:i]))
		next = head.GetNodeAuto(string(buf.Bytes()[begin:i]))
		if next == nil {
			return head, path[begin:]
		}
		head = next.(*node)
		begin = i
	}

	// perfect match.
	//if head != nil {
	//	return head, ""
	//}

	// begin ==0 not match any node
	// begin >0 perfect match.
	return head, path[begin:]
}

// url param pase to ParamterArray.
func (nd *node) ParseParams(pa gnet.Contexter, path_type int, param_str string) {
	// url: /xxx?param1=v  url question mark request
	// it is not necessary to sort  by  keys of node.
	// net/url had also parse this situation.
	// when the code reachs here, the path had no ? and its parameters.
	//if path_type == PATH_T_QUESTION {
	//	begin := 0
	//	var key, value = "", ""
	//	// DEAL_ADD_MARK as symbol. use this sign explain the code using.
	//	// To deal mark of + in here can have much better performance.
	//	// deal_add_mark := false
	//	for i, v := range param_str {

	//		if v == '=' {
	//			key = param_str[begin:i]
	//			begin = i + 1
	//			continue
	//		}

	//		if v == '&' {
	//			value = param_str[begin:i]
	//			begin = i + 1
	//			pa.SetParam(key, value)
	//			key, value = "", ""
	//			continue
	//		}
	//	}
	//	// add last param
	//	if key != "" {
	//		pa.SetParam(key, value)
	//	}
	//	return
	//}

	// url: /xxx/xxx/:param1/:param2
	if path_type == PATH_T_COMMON && len(param_str) > 0 {
		lenk := len(nd.GetKeys())
		if lenk == 0 {
			return
		}
		begin, record := 0, 0
		if param_str[0] == '/' {
			param_str = param_str[1:]
		}
		var key, value = "", ""
		if len(nd.GetKeys()) == 1 {
			key, _ = nd.getKey(0)
			if param_str[len(param_str)-1] == '/' {
				param_str = param_str[:len(param_str)-1]
			}
			pa.SetParam(key, param_str)
			return
		}

		for i, v := range param_str {
			if v != '/' {
				continue
			}
			value = param_str[begin:i]
			begin = i + 1
			key, _ = nd.getKey(record)
			record++
			//fmt.Println("TEST.node.ParseParam", key, value, len(nd.GetKeys()))
			// key could not be empty.
			if key == "" {
				continue
			}
			pa.SetParam(key, value)
			key, value = "", ""
		}
		if key != "" {
			pa.SetParam(key, value)
		}
	}
}

func (nd *node) getKey(key interface{}) (string, int) {
	if k, ok := key.(int); ok {
		if k >= len(nd.keys) {
			return "", -1
		}
		return nd.keys[k], k
	}

	// here not was invoked.
	//if k, ok := key.(string); ok {
	//	for i, v := range nd.keys {
	//		if v == k {
	//			return v, i
	//		}
	//	}
	//}
	return "", -1

}

// implement store.AddRoute
func (nd *node) AddRoute(path string, handle gnet.HandleFunc, param_keys []string) {
	//fmt.Println("--------AddRoute", path)
	f_node, path_l := nd.Lookup(path)

	// It must have error before add route. please check if in route handle.
	// Lookup will return itself if it dose not found any suit node.
	// if f_node == nil {
	// 	panic("[ERROR][GROUTE][FOUND_ROOT_NODE]")
	// }

	// it is wild node and perfect match the path.
	// this node can convert to path node.
	if f_node.GetType() == NODE_T_WILD && path_l == "" {
		f_node.SetPath(path)
		f_node.SetType(NODE_T_PATH)
		f_node.SetHandleFunc(handle)
		f_node.AddKey(param_keys)
		return
	}

	// fmt.Println("--------AddRoute", path, "left path:", param_keys)
	// Second condition was supported for "/" .
	if path_l == "" || (path_l == "/" && path != path_l) {
		panic("[ERROR][GROUTER][ADD_ROUTE]PATH_EXIST[" + path + "]LEFT_PATH[" + path_l + "]")
	}

	path_l_slice := strings.Split(path_l, "/")
	lenp := len(path_l_slice)

	head := &node{
		children: make(map[string]Node),
	}

	next := head
	next.SetType(NODE_T_WILD)
	if lenp >= 2 && path_l_slice[0] == "" {
		path_l_slice = path_l_slice[1:]
		lenp--
	}

	//fmt.Println("--------AddRoute", path, "left path:", path_l)
	for i, indices := range path_l_slice {
		if indices == "" {
			break
		}
		if i == 0 {
			head.SetType(NODE_T_WILD)
			head.SetIndices("/" + lowercase(indices))
			continue
		}

		// that is our need node.
		// first condition is ready for "/"
		// second condition checked whether it is a leaf node.

		// not a leaf node. create a wild node.
		next_tmp := &node{
			children: make(map[string]Node),
		}
		next_tmp.SetType(NODE_T_WILD)
		next_tmp.SetIndices("/" + lowercase(indices))

		// link with before.
		next.AddNode(next_tmp)
		next = next_tmp
		//next.SetIndices("/" + indices)
	}

	ok := f_node.AddNode(head)
	// there should be no such situation. if there is, there is a problem with previous procedure.
	if !ok {
		panic("[ERROR][GROUTER][ADD_NODE][CHILD_EXIST]CHILD_KEY[" + head.GetIndices() + "]")
	}
	next.AddKey(param_keys)
	next.SetPath(path)
	next.SetType(NODE_T_PATH)
	next.SetHandleFunc(handle)

}

func (nd *node) GetChildren() map[string]Node {
	return nd.children
}

// =======================================================================
