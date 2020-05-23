package grouter

import (
//"strings"
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

func (nd *node) AddNode(one Node) {
	nd.children[one.GetIndices()] = one
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

// root lookup node.
// Match the longest chain
func (nd *node) Lookup(path string) (Node, string) {
	var begin = 0

	next := nd
	head := nd

	for i, v := range path {
		if v != '/' {
			continue
		}
		next = head.GetNodeAuto(path[begin:i]).(*node)
		if next == nil {
			return head, path[i:]
		}
		head = next
		begin = i
	}

	return nd, path[begin:]
}

func (nd *node) GetChildren() map[string]Node {
	return nil
}
