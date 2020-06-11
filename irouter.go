package grouter

import (
	//"io"
	"github.com/slclub/gcore/flow"
	"github.com/slclub/gnet"
	"github.com/slclub/gnet/permission"
	//"net/http"
)

//====================================================================================
//

//=====================================================================================

// router haddle func

type Router interface {
	flow.IExecuteNode
	// path router store get and set.
	GetStore() Store
	SetStore(Store)

	// path decoder binding and get
	GetDecoder() Path
	SetDecoder(Path)

	Handle(string, string, gnet.HandleFunc)
	//Execute(gnet.Contexter)
	// Handle shortcut
	// Handle(string, string, HandleFunc)
	GET(string, gnet.HandleFunc)
	HEAD(string, gnet.HandleFunc)
	OPTIONS(string, gnet.HandleFunc)
	POST(string, gnet.HandleFunc)
	PATCH(string, gnet.HandleFunc)
	PUT(string, gnet.HandleFunc)
	DELETE(string, gnet.HandleFunc)
	ANY(string, gnet.HandleFunc)
	//just for test. can directly bind with http lisen
	//ServeHTTP(res http.ResponseWriter, req *http.Request)

	// except 200 3xx
	// the rest of the error code should have handle function defined.
	CodeHandle(int) gnet.HandleFunc
	BindCodeHandle(int, gnet.HandleFunc)
	// static file server.
	ServerFile(string, string, ...bool)
}

// --------------------------------------------------------------------------------------
// this interface was invoked by router directly.
type Lookuper interface {
	// lookup path return lasest node and left not valid path
	Lookup(string) (Node, string)
}

// this interface was invoked by router directly.
type RouterWriter interface {
	AddRoute(string, gnet.HandleFunc, []string)
}

// --------------------------------------------------------------------------------------

// router node interface.
// if you overide the node and struct is very different. you should use both the your defined node interface and this node interface.
type Node interface {
	Lookuper
	RouterWriter
	permission.IAccess

	GetPath() string
	SetPath(string)

	GetHandleFunc() gnet.HandleFunc
	SetHandleFunc(gnet.HandleFunc)

	GetType() uint8
	SetType(uint8)

	GetIndices() string
	SetIndices(string)

	// Get node from children.
	// Overide one of them will be ok.
	GetNodeAuto(interface{}) Node
	GetNodeInt(int) Node
	GetNodeStr(string) Node
	//GetNodes() map[string]Node

	// Add a new node to the childrens.
	AddNode(Node) bool

	AddKey([]string)
	GetKeys() []string

	GetChildren() map[string]Node
	// parse url param
	// @param	@1	ParamArray	save the params
	// @param	@2	url type	PATH_T_COMMON==1: /xxx/:param	PATH_T_QUESTION==2:/xxx?param1=value1&param2=value2
	//							PATH_T_FILE==3: file type
	// @param	@3	left path as param string.
	ParseParams(gnet.Contexter, int, string)
}

// Store the nodes generated by url.
type Store interface {
	// paramter should use url path or method or msgid(int) .
	// look up from store return node.
	// @param	1	string	method or path
	Lookuper
	// save node.
	Save(Node)

	// create an new store
	Create() Store

	CreateNode(int) Node
}

// path process interface
type Path interface {
	// it was invoked when grouter add router provied by devloper
	// @param	path	string
	// @return  @1		proessed path
	// @return  @2		param array.
	Parse(string) (string, []string)

	// it was invoked whne client request.
	// @param	path	string
	// @return	@1		processd path.
	// @return	@2		params was returned. maybe was a file etc.
	Decode(string) (int, string, interface{})

	//GetType() int

	// you can use this method to set the router is supported case sensitive.
	CaseSensitive(bool)
}
