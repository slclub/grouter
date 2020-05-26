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
}
