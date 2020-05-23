package grouter

import (
	//"errors"
	"fmt"
	"github.com/slclub/link"
	//"strings"
)

func init() {
	fmt.Println("[GROUTER][DEV]")
}

type router struct {
	store Store
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
// Parse url and generate node and save it.
func (r *router) GetStore() Store {
	return r.store
}

func (r *router) SetStore(st Store) {
	r.store = st
}

func (r *router) check(path string) (bool, error) {

	// path check
	if len(path) < 1 || path[0] != '/' {
		panic("")
	}

	return true, nil
}

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

}

// request execute
// when a client request, this function will be called.
func (r *router) Execute(Context) {
}
