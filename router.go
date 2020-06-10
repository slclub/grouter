package grouter

import (
	//"errors"
	"fmt"
	"github.com/slclub/gcore/flow"
	"github.com/slclub/gnet"
	"github.com/slclub/link"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	fmt.Println("[GROUTER][DEV]")
}

type router struct {
	flow.ExecuteNode
	store        Store
	decoder      Path
	code_handles map[int]gnet.HandleFunc
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
		code_handles: make(map[int]gnet.HandleFunc),
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

func (r *router) CodeHandle(error_code int) gnet.HandleFunc {
	return r.code_handles[error_code]
}

func (r *router) BindCodeHandle(error_code int, handle gnet.HandleFunc) {
	if handle == nil {
		return
	}
	r.code_handles[error_code] = handle
}

// ------------------------------------------shortcut-start------------------------------------------------
// shortcut for router.Handle
func (r *router) GET(path string, handle gnet.HandleFunc) {
	r.Handle(http.MethodGet, path, handle)
}

func (r *router) HEAD(path string, handle gnet.HandleFunc) {
	r.Handle(http.MethodHead, path, handle)
}
func (r *router) OPTIONS(path string, handle gnet.HandleFunc) {
	r.Handle(http.MethodOptions, path, handle)
}
func (r *router) POST(path string, handle gnet.HandleFunc) {
	r.Handle(http.MethodPost, path, handle)
}
func (r *router) PUT(path string, handle gnet.HandleFunc) {
	r.Handle(http.MethodPut, path, handle)
}
func (r *router) PATCH(path string, handle gnet.HandleFunc) {
	r.Handle(http.MethodPatch, path, handle)
}
func (r *router) DELETE(path string, handle gnet.HandleFunc) {
	r.Handle(http.MethodDelete, path, handle)
}

func (r *router) ANY(path string, handle gnet.HandleFunc) {
	r.Handle("ANY", path, handle)
}

// ------------------------------------------shortcut-end--------------------------------------------------

// implement Router.Handle
func (r *router) Handle(method, path string, handle gnet.HandleFunc) {
	if method == "" {
		panic("[ERROR][ROUTER][HANDLE] method is empty")
	}

	// not pass panic
	r.check(path)

	//if ok, err := r.check(path); !ok {
	//	link.ERROR(err)
	//	return
	//}

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

// static serve file system.
// and static file 404 demo.
// func error404Handler(w http.ResponseWriter, r *http.Request) {
// 	http.Error(w, "404 not found", http.StatusNotFound)
// }
//
// func staticHandler(w http.ResponseWriter, r *http.Request) {
// 	name := path.Clean(r.URL.Path)
// 	if _, err := os.Stat(name); err != nil {
// 		if os.IsNotExist(err) {
// 			error404Handler(w, r)
// 			return
// 		}
//
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}
//
// 	return http.ServeFile(w, r, name)
// }
func (r *router) ServerFile(path string, root_path string, args ...bool) {
	root_fs := http.Dir(root_path)
	file_server := http.FileServer(root_fs)
	if i := strings.IndexByte(path, ':'); i >= 1 {
		file_server = http.StripPrefix(path[:i], file_server)
	} else {
		file_server = http.StripPrefix(path, file_server)
	}
	tail := false
	if len(args) >= 1 {
		tail = args[0]
	}

	r.GET(path, func(ctx gnet.Contexter) {
		//if _, noListing := root_fs.(http.FileSystem); noListing {
		//	r.CodeHandle(404)(ctx)
		//	fmt.Println("GROUTER. static file 404", noListing)
		//	return
		//}

		// Use or deny floder list. decide by tail of request.URL.Path is '/'
		// 禁用或启用 静态目录浏览仅仅需要处理 请求的末尾时候有 '/'

		file := ctx.Request().GetHttpRequest().URL.Path
		file_new, _ := ctx.Request().GetString("filepath")
		// fmt.Println("REQUEST_PATH STEP1", file, "param", file_new, "http.FileSystem:", root_fs)
		lf := len(file)

		if file_new == "" {
			if lf > len(path) {
				file_new = file[len(path):]
			} else {
				file_new = ""
			}
		}

		// Request floder
		if gnet.IsDir(root_path + file_new) {
			if tail {
				if file[lf-1] != '/' {
					ctx.Request().GetHttpRequest().URL.Path += "/"
				}
			} else {
				//if !tail && file[lf-1] == '/' {
				//	ctx.Request().GetHttpRequest().URL.Path += file[:lf-1]
				//}
				r.CodeHandle(404)(ctx)
				return
			}
		} else {
			// request file
			if file_new[len(file_new)-1] == '/' {
				file_new = file_new[:len(file_new)-1]
				ctx.Request().GetHttpRequest().URL.Path = file[:lf-1]
			}
		}
		// Deny access floder. by your setting.

		// Can not found source file. use an public 404 handle.
		f, err := root_fs.Open(file_new)
		// fmt.Println("REQUEST_PATH STEP2", file, "param", file_new, "err:", err)
		if err != nil {
			r.CodeHandle(404)(ctx)
			return
		}
		f.Close()
		file_server.ServeHTTP(ctx.Response(), ctx.Request().GetHttpRequest())
	})
}

// request execute
// when a client request, this function will be called.
func (r *router) Execute(ctx gnet.Contexter) {
	var req *http.Request
	req = ctx.Request().GetHttpRequest()
	http_method := req.Method
	// Reduces empty request performance by half.
	// also need to QueryUnescape the param
	//path_type, path, params_str := r.GetDecoder().Decode(req.URL.String())
	//path_type, path, params_str := r.GetDecoder().Decode(req.URL.Path)
	path_type, path, _ := r.GetDecoder().Decode(req.URL.Path)

	var root Node
	var nothing string

WALK_AGAIN:
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
		ctx.SetHandler(handle)
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
			// here support ANY method.
			if http_method != "ANY" && http_method != http.MethodConnect && http_method != http.MethodOptions {
				http_method = "ANY"
				goto WALK_AGAIN
			}

			goto WALK_404
		}
		handle := node.GetHandleFunc()
		if handle == nil {
			// here support ANY method.
			if http_method != "ANY" && http_method != http.MethodConnect && http_method != http.MethodOptions {
				http_method = "ANY"
				goto WALK_AGAIN
			}
			goto WALK_404
		}
		if left_path != "/" {
			param_str, err := url.QueryUnescape(left_path)
			if err != nil {
				param_str = left_path
			}
			node.ParseParams(ctx, path_type, param_str)
		}
		// test
		//handle(ctx)

		ctx.SetHandler(handle)
	}
	// TODO: get param from net/url.URL.Query()
	// TODO: get param from req.ParseForm.
	// TODO: query param etc.
	// TODO: return process node. include handle param and scope.
	return
WALK_404:
	not_handle := r.CodeHandle(http.StatusNotFound)
	//if not_handle == nil {
	//	return
	//}
	//not_handle(ctx)
	ctx.SetHandler(not_handle)
}

//func (r router) againAny(ctx gnet.Contexter) {
//	var req *http.Request
//	req = ctx.Request().GetHttpRequest()
//
//	req.Method = "ANY"
//	r.Execute(ctx)
//}

// for test. cover test.
// can bind with Http.ListenAndServe
//func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
//
//}

// =========================================code handle func ===============================================
func http_404_handle(ctx gnet.Contexter) {
	ctx.Response().WriteHeader(http.StatusNotFound)
	ctx.Response().Write([]byte("grouter 404 not found"))
	//ctx.Response().Flush()
	//ctx.Response().InitSelf(nil)

	// fmt.Println("---------------handle-not found--------------------")
}

func http_405_handle(ctx gnet.Contexter) {
	ctx.Response().WriteHeader(http.StatusMethodNotAllowed)
	ctx.Response().Write([]byte("grouter 405 not not allowed!"))
	//ctx.Response().Flush()
}

func http_500_handle(ctx gnet.Contexter) {
	ctx.Response().WriteHeader(http.StatusInternalServerError)
	ctx.Response().Write([]byte("grouter 500 server internal error!"))

	//ctx.Response().Flush()
	link.ERROR("[500] server internal error!")
}
