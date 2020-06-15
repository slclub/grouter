package grouter

import (
	"github.com/slclub/gcore/execute"
	"github.com/slclub/gnet"
	"github.com/slclub/gnet/permission"
	"github.com/slclub/link"
	"github.com/slclub/utils"
)

type IGroup interface {
	Group(f func(IGroup))
	Use(f gnet.HandleFunc) IGroup
	Deny(f gnet.HandleFunc) IGroup
	Reset()
	GetStack() StackMiddle
	Execute(access permission.Accesser)
	GetMiddle(access permission.Accesser) execute.Middler
}

type groupObject struct {
	middle_group execute.Middler
	middle_node  execute.Middler
	stack_middle StackMiddle

	// true: use group add middle
	// false: add middler to node.
	group_state bool
}

var _ IGroup = &groupObject{}

func NewGroup() IGroup {
	return &groupObject{
		stack_middle: make(StackMiddle),
	}
}

func (gr *groupObject) Group(f func(IGroup)) {
	gr.middle_group = execute.NewMiddle(ROUTER_MIDDLER_KEY)
	gr.group_state = true
	f(gr)
	gr.middle_group = nil
	gr.group_state = false
}

func (gr *groupObject) Use(f gnet.HandleFunc) IGroup {
	link.DEBUG_PRINT("[GROUTE][GROUP][USE][STEP1]", utils.EOL)
	if gr.group_state {
		gr.middle_group.Use(f)
	} else {
		gr.autoNodeMiddle()
		gr.middle_node.Use(f)
	}
	return gr
}

func (gr *groupObject) Deny(f gnet.HandleFunc) IGroup {
	if gr.group_state {
		gr.middle_group.Deny(f)
	} else {
		gr.autoNodeMiddle()
		gr.middle_node.Deny(f)
	}
	return gr
}

func (gr *groupObject) autoNodeMiddle() {
	if gr.middle_node == nil {
		gr.middle_node = execute.NewMiddle(ROUTER_MIDDLER_KEY)
	}
}

func (gr *groupObject) Reset() {
	gr.middle_node = nil
}

func (gr *groupObject) GetStack() StackMiddle {
	return gr.stack_middle
}

func (gr *groupObject) Execute(access permission.Accesser) {
	if nil == access {
		gr.Reset()
		return
	}
	link.DEBUG_PRINT("[GROUTE][GROUP][Execute][STEP1]", gr.middle_group, utils.EOL)
	if gr.middle_node == nil && gr.middle_group == nil {
		gr.Reset()
		return
	}
	link.DEBUG_PRINT("[GROUTE][GROUP][Execute][STEP2]", utils.EOL)
	// get mid
	mid := gr.stack_middle.Get(access.GetAID())
	if mid == nil {
		mid = execute.NewMiddle(ROUTER_MIDDLER_KEY)
		gr.stack_middle[access.GetAID()] = mid
	}
	// combine
	mid.Combine(gr.middle_group)
	mid.Combine(gr.middle_node)

	// reset access
	i := 0
	for {
		_, name := mid.GetHandle(i)
		i++
		if name == "" {
			break
		}
		ind, _ := mid.Invoker().GetId(name)
		scope := mid.Invoker().GetScopeById(ind)
		access.Set(ind, scope)
	}

	gr.Reset()
}
func (gr *groupObject) GetMiddle(access permission.Accesser) execute.Middler {
	if access == nil {
		return nil
	}
	return gr.stack_middle.Get(access.GetAID())
}

var Group = NewGroup()
