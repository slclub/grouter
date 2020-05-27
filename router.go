package grouter

import (
	//"errors"
	"fmt"
	"github.com/slclub/link"
	"net/http"
	"net/url"
)

func init() {
	fmt.Println("[GROUTER][DEV]")
}

type router struct {
	store        Store
	decoder      Path
	code_handles map[int]HandleFunc
}

//type param struct {
//	key   string
//	value string
//}
//
//// implement gcontext.Param:GetKey
//func (p *param) GetKey() string {
//	return p.key
//}
//
//// implement gcontext.Param:GetValue
//func (p *param) GetValue() string {
//	return p.value
//}
//
//// impement gcontext.Param.Set
//func (p *param) Set(key, value string) {
//	p.key = key
//	p.value = value
//}

// -----------------------------------------router ----------------------------------------

func NewRouter() Router {

	r := &router{
		store:        NewStore(),
		decoder:      NewPath(),
		code_handles: make(map[int]HandleFunc),
	}
	//bind code handle
	r.BindCodeHandle(http.StatusNotFound, http_404_handle)
	//http.StatusMethodNotAllowed
	r.BindCodeHandle(http.StatusMethodNotAllowed, http_405_handle)
	r.BindCodeHandle(http.StatusInternalServerError, http_500_handle)
	return r
}

// Parse url and generate node and save it.
func (r *router) GetStore() Store {
	return r.store
}

// can change Store belong to your.
func (r *router) SetStore(st Store) {
	r.store = st
}

func (r *router) GetDecoder() Path {
	return r.decoder
}

// you can change Path belong to your.
func (r *router) SetDecoder(path_dec Path) {
	r.decoder = path_dec
}

func (r *router) check(path string) (bool, error) {

	// do redirect path check
	if len(path) < 1 || path[0] != '/' {
		//return false, errors.New("request url not found!")
		//http.Redirect(w, req, "/404/get", http.StatusPermanentRedirect)
		panic("[ERROR][GROUTER][URL][EMPTY] or [NOT START WITH /]")
	}

	return true, nil
}

func (r *router) CodeHandle(error_code int) HandleFunc {
	return r.code_handles[error_code]
}

func (r *router) BindCodeHandle(error_code int, handle HandleFunc) {
	if handle == nil {
		return
	}
	r.code_handles[error_code] = handle
}

// ------------------------------------------shortcut-start------------------------------------------------
// shortcut for router.Handle
func (r *router) GET(path string, handle HandleFunc) {
	r.Handle(http.MethodGet, path, handle)
}

func (r *router) HEAD(path string, handle HandleFunc) {
	r.Handle(http.MethodHead, path, handle)
}
func (r *router) OPTIONS(path string, handle HandleFunc) {
	r.Handle(http.MethodOptions, path, handle)
}
func (r *router) POST(path string, handle HandleFunc) {
	r.Handle(http.MethodPost, path, handle)
}
func (r *router) PUT(path string, handle HandleFunc) {
	r.Handle(http.MethodPut, path, handle)
}
func (r *router) PATCH(path string, handle HandleFunc) {
	r.Handle(http.MethodPatch, path, handle)
}
func (r *router) DELETE(path string, handle HandleFunc) {
	r.Handle(http.MethodDelete, path, handle)
}

// ------------------------------------------shortcut-end--------------------------------------------------

// implement Router.Handle
func (r *router) Handle(method, path string, handle HandleFunc) {
	if method == "" {
		panic("[ERROR][ROUTER][HANDLE] method is empty")
	}

	if ok, err := r.check(path); !ok {
		link.ERROR(err)
		return
	}

	if handle == nil {
		link.ERROR("[ROUTE]HandleFunc is empty!")
		panic("HandleFunc is empty!")
	}

	// here need path parse the path. get all param keys.
	// params_keys := make([]string, 1)
	path_l, param_keys := r.GetDecoder().Parse(path)

	// TODO:before lookup root node. we need to check method was allowed.
	// if not check
	root, _ := r.store.Lookup(method)
	if root == nil {
		panic("[ERROR][GROUTER][NOT_ALLOWD]METHOD[" + method + "]")
	}
	root.AddRoute(path_l, handle, param_keys)

}

// request execute
// when a client request, this function will be called.
func (r *router) Execute(ctx Contexter) {
	var req *http.Request
	req = ctx.GetRequest("http").(*http.Request)
	http_method := req.Method
	// Reduces empty request performance by half.
	// also need to QueryUnescape the param
	//path_type, path, params_str := r.GetDecoder().Decode(req.URL.String())
	//path_type, path, params_str := r.GetDecoder().Decode(req.URL.Path)
	path_type, path, _ := r.GetDecoder().Decode(req.URL.Path)

	var root Node
	var nothing string
	if path_type > 0 {
		root, nothing = r.GetStore().Lookup(http_method)
	}
	// method not allowed. handle 405
	if root == nil {
		handle := r.CodeHandle(http.StatusMethodNotAllowed)
		if handle != nil {
			//handle()
			return
		}
		link.ERROR("[GROUTER]", nothing)
		return
	}

	//fmt.Println("---------sys type", PATH_T_QUESTION, "-------------", path_type, "path", path, "URL.Path:", req.URL.Path)
	//if path_type == PATH_T_QUESTION {
	//	node, left_path := root.Lookup(path)
	//	if node == nil {
	//		goto WALK_404
	//	}
	//	handle := node.GetHandleFunc()
	//	if left_path != "" || handle == nil {
	//		goto WALK_404
	//	}

	//	node.ParseParams(ctx, path_type, params_str.(string))
	//	// test
	//	handle(ctx)
	//}

	if path_type == PATH_T_COMMON {
		node, left_path := root.Lookup(path)
		if node == nil {
			goto WALK_404
		}
		handle := node.GetHandleFunc()
		if handle == nil {
			goto WALK_404
		}
		param_str, err := url.QueryUnescape(left_path)
		if err != nil {
			param_str = left_path
		}
		node.ParseParams(ctx, path_type, param_str)
		// test
		handle(ctx)

	}
	// TODO: get param from net/url.URL.Query()
	// TODO: get param from req.ParseForm.
	// TODO: query param etc.
	// TODO: return process node. include handle param and scope.
	return
WALK_404:
	not_handle := r.CodeHandle(http.StatusNotFound)
	if not_handle == nil {
		return
	}
	not_handle(ctx)

}

// for test. cover test.
// can bind with Http.ListenAndServe
//func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
//
//}

// =========================================code handle func ===============================================
func http_404_handle(ctx Contexter) {
	ctx.Status(http.StatusNotFound)
	ctx.GetHttpResponse().Write([]byte("404 not found"))

	fmt.Println("---------------handle-not found--------------------")
}

func http_405_handle(ctx Contexter) {
	ctx.Status(http.StatusMethodNotAllowed)
	ctx.GetHttpResponse().Write([]byte("405 not not allowed!"))
}

func http_500_handle(ctx Contexter) {
	ctx.Status(http.StatusInternalServerError)
	ctx.GetHttpResponse().Write([]byte("500 server internal error!"))
	link.ERROR("[500] server internal error!")
}
