package opcda

import (
	"testing"
	"time"

	"github.com/huskar-t/opcda/com"

	"github.com/stretchr/testify/assert"
)

func TestOPCItem(t *testing.T) {
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
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, items, item.GetParent())
	assert.Equal(t, TestWriteItem, item.GetItemID())
	assert.Equal(t, "", item.GetAccessPath())
	assert.Equal(t, true, item.GetIsActive())
	assert.Equal(t, uint32(3), item.GetAccessRights())
	euType, err := item.GetEUType()
	assert.NoError(t, err)
	t.Log(euType)
	euInfo, err := item.GetEUInfo()
	assert.NoError(t, err)
	t.Log(euInfo)
	time.Sleep(time.Millisecond * 100)
	v, q, ts, err := item.Read(OPC_DS_CACHE)
	assert.NoError(t, err)
	t.Log(v, q, ts)
	assert.Equal(t, v, item.GetValue())
	assert.Equal(t, q, item.GetQuality())
	assert.Equal(t, ts, item.GetTimestamp())
	err = item.Write(int32(12))
	assert.NoError(t, err)
	v, _, _, err = item.Read(OPC_DS_CACHE)
	assert.NoError(t, err)
	assert.Equal(t, int32(12), v)
	vt := item.GetCanonicalDataType()
	t.Log(vt)
	assert.Equal(t, com.VT_I4, vt)
}
