package opcda

import (
	"testing"

	"github.com/huskar-t/opcda/com"

	"github.com/stretchr/testify/assert"
)

func TestOPCItems_RequestedDataType(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	items := group.OPCItems()
	assert.Equal(t, group, items.GetParent())
	assert.Equal(t, com.VT_EMPTY, items.GetDefaultRequestedDataType())
	items.SetDefaultRequestedDataType(com.VT_I4)
	assert.Equal(t, com.VT_I4, items.GetDefaultRequestedDataType())
	items.SetDefaultRequestedDataType(com.VT_EMPTY)
}

func TestOPCItems_DefaultAccessPath(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	accessPath := items.GetDefaultAccessPath()
	assert.Equal(t, "", accessPath)
	items.SetDefaultAccessPath("test")
	assert.Equal(t, "test", items.GetDefaultAccessPath())
	items.SetDefaultAccessPath("")
}

func TestOPCItems_DefaultActive(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	items.SetDefaultActive(false)
	assert.Equal(t, false, items.GetDefaultActive())
	items.SetDefaultActive(true)
	assert.Equal(t, true, items.GetDefaultActive())
}

func TestOPCItems_AddItems(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()

	boolItem, err := items.AddItem(TestBoolItem)
	assert.NoError(t, err)
	assert.NotNil(t, boolItem)
	itemList, errors, err := items.AddItems([]string{TestFloatItem, TestWriteItem})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(itemList))
	assert.Equal(t, 2, len(errors))
	for i := 0; i < 1; i++ {
		assert.NoError(t, errors[i])
	}
	assert.Equal(t, 3, items.GetCount())
	item0, err := items.Item(0)
	assert.NoError(t, err)
	assert.Equal(t, boolItem, item0)
	item1, err := items.Item(1)
	assert.NoError(t, err)
	assert.Equal(t, itemList[0], item1)
	item2, err := items.Item(2)
	assert.NoError(t, err)
	assert.Equal(t, itemList[1], item2)
	item3, err := items.Item(3)
	assert.Error(t, err)
	assert.Nil(t, item3)
	item0, err = items.ItemByName(TestBoolItem)
	assert.NoError(t, err)
	assert.Equal(t, boolItem, item0)
	item1, err = items.ItemByName(TestFloatItem)
	assert.NoError(t, err)
	assert.Equal(t, itemList[0], item1)
	item2, err = items.ItemByName(TestWriteItem)
	assert.NoError(t, err)
	assert.Equal(t, itemList[1], item2)
	item3, err = items.ItemByName("test")
	assert.Error(t, err)
	assert.Nil(t, item3)
	item0, err = items.GetOPCItem(boolItem.GetServerHandle())
	assert.NoError(t, err)
	assert.Equal(t, boolItem, item0)
	item1, err = items.GetOPCItem(itemList[0].GetServerHandle())
	assert.NoError(t, err)
	assert.Equal(t, itemList[0], item1)
	item2, err = items.GetOPCItem(itemList[1].GetServerHandle())
	assert.NoError(t, err)
	assert.Equal(t, itemList[1], item2)
	item3, err = items.GetOPCItem(0)
	assert.Error(t, err)
	assert.Nil(t, item3)
	items.Remove([]uint32{boolItem.GetServerHandle()})
	assert.Equal(t, 2, items.GetCount())
	items.Remove([]uint32{itemList[0].GetServerHandle(), itemList[1].GetServerHandle()})
	assert.Equal(t, 0, items.GetCount())
}

func TestOPCItems_Validate(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()

	errs, err := items.Validate([]string{TestBoolItem, TestFloatItem, TestWriteItem}, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.Nil(t, e)
	}
	errs, err = items.Validate([]string{"xxx", "x"}, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(errs))
	for _, e := range errs {
		assert.Error(t, e)
	}
	errs, err = items.Validate([]string{"xxx", TestFloatItem, TestWriteItem}, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(errs))
	assert.Error(t, errs[0])
	assert.NoError(t, errs[1])
	assert.NoError(t, errs[2])
	errs, err = items.Validate([]string{TestFloatItem, TestFloatItem, TestWriteItem}, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.NoError(t, e)
	}
}

func TestOPCItems_SetActive(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()

	opcItems, errs, err := items.AddItems([]string{TestBoolItem, TestFloatItem, TestWriteItem})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(opcItems))
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.NoError(t, e)
	}
	errs = items.SetActive([]uint32{opcItems[0].GetServerHandle(), opcItems[1].GetServerHandle(), opcItems[2].GetServerHandle()}, false)
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.NoError(t, e)
	}
	for _, item := range opcItems {
		assert.False(t, item.GetIsActive())
	}
}

func TestOPCItems_SetClientHandles(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()

	opcItems, errs, err := items.AddItems([]string{TestBoolItem, TestFloatItem, TestWriteItem})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(opcItems))
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.NoError(t, e)
	}
	errs = items.SetClientHandles([]uint32{opcItems[0].GetServerHandle(), opcItems[1].GetServerHandle(), opcItems[2].GetServerHandle()}, []uint32{100, 200, 300})
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.NoError(t, e)
	}
	for i, item := range opcItems {
		assert.Equal(t, uint32(i+1)*100, item.GetClientHandle())
	}
}

func TestOPCItems_SetDataTypes(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()

	opcItems, errs, err := items.AddItems([]string{TestBoolItem, TestFloatItem, TestWriteItem})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(opcItems))
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.NoError(t, e)
	}
	errs = items.SetDataTypes([]uint32{opcItems[0].GetServerHandle(), opcItems[1].GetServerHandle(), opcItems[2].GetServerHandle()}, []com.VT{com.VT_I4, com.VT_I4, com.VT_I4})
	assert.Equal(t, 3, len(errs))
	for _, e := range errs {
		assert.NoError(t, e)
	}
	for _, item := range opcItems {
		assert.Equal(t, com.VT_I4, item.GetRequestedDataType())
	}
}
