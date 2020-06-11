package grouter

import (
	"github.com/slclub/gcore/execute"
	"github.com/slclub/gnet"
)

const ROUTER_MIDDLER_KEY = "grouter_middler"

type StackMiddle map[int]execute.Middler

func (sm *StackMiddle) Use(access_id int, handle gnet.HandleFunc) {
	m := (*sm).Get(access_id)
	if m == nil {
		(*sm)[access_id] = execute.NewMiddle(ROUTER_MIDDLER_KEY)
		m = (*sm)[access_id]
	}
	m.Use(handle)
}

func (sm *StackMiddle) Deny(access_id int, handle gnet.HandleFunc) {
	m := (*sm).Get(access_id)
	if m == nil {
		(*sm)[access_id] = execute.NewMiddle(ROUTER_MIDDLER_KEY)
		m = (*sm)[access_id]
	}

	m.Deny(handle)
}

func (sm *StackMiddle) Get(access_id int) execute.Middler {
	m, _ := (*sm)[access_id]
	return m
}
