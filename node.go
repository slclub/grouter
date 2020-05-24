package grouter

import (
	"strings"
)

const (
	// method node . as every method root node.
	NODE_T_ROOT = 1
	// wild node. have no path matched.
	NODE_T_WILD = 2

	// valid path node.
	NODE_T_PATH = 3
)

// the place of store the nodes.
type nodeStore struct {
	nodes map[string]Node
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
	handle   HandleFunc
	indices  string
	children map[string]Node
	keys     []string
}

// ============================store router node function =====================
func (st *nodeStore) Save(n Node) {
	st.nodes[n.GetIndices()] = n
}

func (st *nodeStore) Create() Store {
	return &nodeStore{}
}

func (st *nodeStore) Lookup(method string) (n Node) {
	n = st.nodes[method]
	if n == nil {
		n = st.CreateNode(0)
		n.SetIndices(method)
		n.SetType(NODE_T_ROOT)
		st.nodes[method] = n
	}
	return
}

func (st *nodeStore) CreateNode(param_num int) Node {

	return &node{
		children: make(map[string]Node),
		keys:     make([]string, param_num),
	}
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

func (nd *node) SetHandleFunc(handle HandleFunc) {
	nd.handle = handle
}

func (nd *node) GetHandleFunc() HandleFunc {
	return nd.handle
}

func (nd *node) GetIndices() string {
	return nd.indices
}

func (nd *node) SetIndices(indices string) {
	nd.indices = indices
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
	if nd.children[one.GetIndices()] != nil {
		return false
	}
	nd.children[one.GetIndices()] = one
	return true
}

func (nd *node) AddKey(key string) {
	for _, k := range nd.keys {
		if key == k {
			return
		}
	}
	nd.keys = append(nd.keys, key)
}

func (nd *node) GetKeys() []string {
	return nd.keys
}

// Match the longest chain
func (nd *node) Lookup(path string) (Node, string) {
	var begin = 0

	next := nd
	head := nd
	len_path := len(path)

	if len_path == 1 {
		next = head.GetNodeAuto(path[begin:1]).(*node)
		if next == nil {
			return nd, path
		}
		return next, ""
	}

	for i:=1; i<= len_path; i++ {
		if path[1] != '/' {
			continue
		}
		next = head.GetNodeAuto(path[begin:i]).(*node)
		if next == nil {
			return head, path[begin:]
		}
		head = next
		begin = i
	}

	// perfect match.
	if head != nil {
		return head, ""
	}

	return nd, path[begin:]
}
// implement store.AddRoute
func (nd *node)AddRoute(path string, handle HandleFunc, param_keys []string) {
	f_node, path_l := nd.Lookup(path)
	// Second condition was supported for "/" .
	if path_l == "" || (path_l == "/" && path != path_l){
		panic("[ERROR][GROUTER][ADD_ROUTE][PATH_EXIST]" + path)
	}

	// It must have error before add route. please check if in route handle.
	if f_node == nil  {
		panic("[ERROR][GROUTE][FOUND_ROOT_NODE]")
	}

	path_l_slice := strings.Split(path_l, "/")
	lenp := len(path_l_slice)

	head := &node {
			children : make(map[string]Node),
		}
	//head.SetType(NODE_T_WILD)

	//// the lenp == 1 when path is "/"
	//// head just is node we needed.
	//if (lenp == 1 || (lenp == 2 && path_l_slice[1] == "" ) {
	//	f_node.AddNode(head)
	//	head.keys = param_keys
	//	head.SetPath(path)
	//	head.SetType(NODE_T_PATH)
	//	head.SetIndices("/" + path_l_slice[0])
	//	head.SetHandleFunc(handle)
	//	return
	//}
	//// rid of empty item.
	//path_l_slice = path_l_slice[:lenp -1]

	next := head

	for i, indices := range path_l_slice {
		next.SetType(NODE_T_WILD)
		next.SetIndices("/" + indices)
		// that is our need node.
		// first condition is ready for "/"
		// second condition checked whether it is a leaf node.
		if indices == "" || i + 2 == lenp{
			ok :=f_node.AddNode(head)
			// there should be no such situation. if there is, there is a problem with previous procedure.
			if !ok {
				panic ("[ERROR][GROUTER][ADD_NODE][CHILD_EXIST]CHILD_KEY[" + head.GetIndices()+ "]")
			}
			next.keys = param_keys
			next.SetPath(path)
			next.SetType(NODE_T_PATH)
			next.SetHandleFunc(handle)
		}
		// not a leaf node. create a wild node.
		next_tmp := &node {
			children : make(map[string]Node),
		}
		// link with before.
		next.AddNode(next_tmp)
		next = next_tmp
	}
}


func (nd *node) GetChildren() map[string]Node {
	return nil
}
