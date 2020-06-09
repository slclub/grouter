package grouter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoreCreate(t *testing.T) {
	// support fmt work! no use
	fmt.Println("support fmt work!")

	var store Store
	st := NewStore()
	store = st.Create()
	root, _ := store.Lookup("ANY")
	root_get, _ := store.Lookup("GET")

	//fmt.Println("STORE_TEST ROOT TYPE VALUE", root.GetType(), "SYS_ROOT", NODE_T_ROOT)
	assert.Equal(t, "ANY", root.GetIndices())
	assert.Equal(t, NODE_T_ROOT, root_get.GetType())
	assert.Equal(t, NODE_T_ROOT, root.GetType())

	root_my, _ := store.Lookup("myself")
	assert.Nil(t, root_my)

	assert.Nil(t, root.GetNodeAuto(1))
	assert.Nil(t, root.GetNodeAuto(int32(1)))
	assert.Nil(t, root.GetNodeAuto("huanghuacai"))
}

func TestNodeSomeMethodSure(t *testing.T) {
	st := NewStore()
	root, _ := st.Lookup("ANY")
	assert.Equal(t, "", root.GetPath())

	nd := st.CreateNode(0)
	nd.SetPath("/find/you")
	nd.SetType(NODE_T_PATH)
	assert.Equal(t, "/find/you", nd.GetPath())
	assert.Equal(t, NODE_T_PATH, nd.GetType())
	nd.AddKey([]string{"ab", "bc", "cd"})
	assert.Equal(t, 3, len(nd.GetKeys()))

	nd, _ = root.Lookup("/")
	assert.NotNil(t, nd)
	assert.Equal(t, "ANY", nd.GetIndices())

	nd = st.CreateNode(0)
	nd.SetPath("/")
	nd.SetIndices("/")
	root.AddNode(nd)
	nd, _ = root.Lookup("/")
	assert.Equal(t, "/", nd.GetIndices())

	nd = st.CreateNode(0)
	root.AddRoute("/best/", http_405_handle, []string{"ni", "hao"})
	nd, _ = root.Lookup("/best/")
	assert.Equal(t, "/best", nd.GetIndices())
	//k1, i1 := nd.getKey(0)
	//assert.Equal(t, "ni", k1)
	//assert.Equal(t, 0, i1)

	root.AddRoute("/Fest/SheM", http_405_handle, []string{"ni", "hao"})
	nd, _ = root.Lookup("/fest/shem/")
	//for i, v := range nd.GetChildren() {
	//	fmt.Println("TEST:node ", i, v.GetIndices())
	//}
	assert.Equal(t, "/shem", nd.GetIndices())
}
