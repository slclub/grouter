package grouter

import (
	"github.com/slclub/gnet"
	"github.com/slclub/link"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupNew(t *testing.T) {
	var group = Group
	r := NewRouter()
	f1 := func(ctx gnet.Contexter) {}
	f2 := func(ctx gnet.Contexter) {}
	f3 := func(ctx gnet.Contexter) {}
	f4 := func(ctx gnet.Contexter) {}
	Group.Use(f3)
	group.Group(func(this IGroup) {
		this.Use(f1).Use(f2)
		r.GET("/f1/now", func(ctx gnet.Contexter) {})
	})

	//fmt.Println("[TEST][GROUP_NEW]", Group.GetStack()[0])

	assert.Equal(t, 1, len(Group.GetStack()))
	assert.Equal(t, 3, Group.GetStack()[0].Size())

	stack := Group.GetStack()
	stack.Use(0, f4)

	assert.Equal(t, 4, stack[0].Size())

	// repeat test
	stack.Use(0, f4)
	assert.Equal(t, 4, stack[0].Size())

	stack.Deny(0, f4)
	assert.Equal(t, 4, stack[0].Size())

	scope := stack[0].Invoker().GetScopeById(3)
	assert.Equal(t, uint8(1), scope)
	link.DEBUG_PRINT("[TEST:GROUP][f4][SCOPE]", scope, "\n")

	scope = stack[0].Invoker().GetScopeById(2)
	assert.Equal(t, uint8(2), scope)

	Group.Execute(nil)

}

func TestGroupDeny(t *testing.T) {
	f1 := func(ctx gnet.Contexter) {}
	Group.Use(f1)
	r := NewRouter()
	r.GET("/f1/now", func(ctx gnet.Contexter) {})
	Group.Deny(f1)

	stack := Group.GetStack()
	f2 := func(ctx gnet.Contexter) {}
	Group.Use(f2)
	r.GET("/f2/now", func(ctx gnet.Contexter) {})
	assert.Equal(t, 1, stack[1].Size())
}
