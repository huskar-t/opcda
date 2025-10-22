package opcda

import (
	"testing"
	"time"

	"github.com/huskar-t/opcda/com"
	"github.com/stretchr/testify/assert"
)

func TestOPCItem_GetParent(t *testing.T) {
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
}

func TestOPCItem_GetItemID(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, TestWriteItem, item.GetItemID())
}

func TestOPCItem_GetAccessPath(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "", item.GetAccessPath())
}

func TestOPCItem_GetIsActive(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, true, item.GetIsActive())
}

func TestOPCItem_GetAccessRights(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, uint32(3), item.GetAccessRights())
}

func TestOPCItem_GetEUType(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	euType, err := item.GetEUType()
	assert.NoError(t, err)
	t.Log(euType)
}

func TestOPCItem_GetEUInfo(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	euInfo, err := item.GetEUInfo()
	assert.NoError(t, err)
	t.Log(euInfo)
}

func TestOPCItem_GetCanonicalDataType(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	vt := item.GetCanonicalDataType()
	t.Log(vt)
	assert.Equal(t, com.VT_I4, vt)
}

func TestOPCItem_WriteError(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_item_write_error")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestWriteErrorItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	time.Sleep(time.Millisecond * 100)
	v, q, ts, err := item.Read(OPC_DS_CACHE)
	assert.NoError(t, err)
	t.Log(v, q, ts)
	assert.Equal(t, v, item.GetValue())
	assert.Equal(t, q, item.GetQuality())
	assert.Equal(t, ts, item.GetTimestamp())
	err = item.Write(int32(12))
	t.Log(err)
	assert.Error(t, err)
}

func TestOPCItem_ReadError(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_item_read_error")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	item, err := items.AddItem(TestReadErrorItem)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	time.Sleep(time.Millisecond * 100)
	_, _, _, err = item.Read(OPC_DS_CACHE)
	assert.Error(t, err)
}

func TestOPCItemRead(t *testing.T) {
	tags := []string{
		"Random.ArrayOfReal8",
		"Random.ArrayOfString",
		"Random.Boolean",
		"Random.Int1",
		"Random.Int2",
		"Random.Int4",
		"Random.Int8",
		"Random.Qualities",
		"Random.Real4",
		"Random.Real8",
		"Random.String",
		"Random.Time",
		"Random.UInt1",
		"Random.UInt2",
		"Random.UInt4",
		"Random.UInt8",
	}
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_item_read")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	itemList, errs, err := items.AddItems(tags)
	if err != nil {
		t.Fatalf("add items failed: %s\n", err)
	}
	for i, err := range errs {
		if err != nil {
			t.Fatalf("add item %s failed: %s\n", tags[i], err)
		}
	}
	time.Sleep(time.Second * 2)
	for i, item := range itemList {
		value, quality, timestamp, err := item.Read(OPC_DS_CACHE)
		if err != nil {
			t.Fatalf("read item failed: %s\n", err)
		}
		if tags[i] != "Random.Qualities" {
			assert.Equal(t, uint16(192), quality)
		}
		t.Logf("%s:\t%s\t%d\t%v\n", tags[i], timestamp, quality, value)
	}
}

func TestOPCItemWrite(t *testing.T) {
	tags := []string{
		"Bucket Brigade.ArrayOfReal8",
		"Bucket Brigade.ArrayOfString",
		"Bucket Brigade.Boolean",
		"Bucket Brigade.Int1",
		"Bucket Brigade.Int2",
		"Bucket Brigade.Int4",
		"Bucket Brigade.Real4",
		"Bucket Brigade.Real8",
		"Bucket Brigade.String",
		"Bucket Brigade.Time",
		"Bucket Brigade.UInt1",
		"Bucket Brigade.UInt2",
		"Bucket Brigade.UInt4",
		".ArrayOfBool",
		".ArrayOfByte",
		".ArrayOfChar",
		".ArrayOfDate",
		".ArrayOfReal4",
		".ArrayOfReal8",
		".ArrayOfShort",
		".ArrayOfString",
		".ArrayOfUlong",
		".ArrayOfUshort",
		".ArrayOfLong",
		".ArrayOfInt",
		".ArrayOfUint",
		".ArrayOfInt64",
		".ArrayOfUint64",
		".Int",
		".Uint",
		".Int64",
		".Uint64",
	}
	now := time.Now()
	millisecond := now.UnixMilli()
	writeNow := time.UnixMilli(millisecond).UTC()
	values := []interface{}{
		[]float64{1.2, 2.3, 3.4},
		[]string{"hello", "world"},
		true,
		int8(1),
		int16(2),
		int32(3),
		float32(5.5777),
		float64(6.6777777777),
		"hello",
		writeNow,
		uint8(7),
		uint16(8),
		uint32(9),
		[]bool{true, false, true},
		[]byte{1, 2, 3},
		[]int8{'a', 'b', 'c'},
		[]time.Time{writeNow, writeNow},
		[]float32{1.2, 2.3, 3.4},
		[]float64{1.2, 2.3, 3.4},
		[]int16{1, 2, 3},
		[]string{"hello", "world"},
		[]uint32{1, 2, 3},
		[]uint16{1, 2, 3},
		[]int32{1, 2, 3},
		[]int{4, 5, 6},
		[]uint{7, 8, 9},
		[]int64{10, 11, 13},
		[]uint64{20, 21, 23},
		int(100),
		uint(200),
		int64(300),
		uint64(400),
	}
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_item_write")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.OPCItems()
	itemList, errs, err := items.AddItems(tags)
	if err != nil {
		t.Fatalf("add items failed: %s\n", err)
	}
	for i, err := range errs {
		if err != nil {
			t.Fatalf("add item %s failed: %s\n", tags[i], err)
		}
	}
	time.Sleep(time.Second * 2)
	for i, item := range itemList {
		var value interface{}
		var quality uint16
		for j := 0; j < 50; j++ {
			err := item.Write(values[i])
			if err != nil {
				t.Fatalf("write item %s failed: %s\n", tags[i], err)
			}
			var timestamp time.Time
			value, quality, timestamp, err = item.Read(OPC_DS_CACHE)
			t.Logf("%s:\t%s\t%d\t%v\n", tags[i], timestamp, quality, value)
			if err != nil {
				t.Fatalf("read item failed: %s\n", err)
			}
			if !assert.ObjectsAreEqual(values[i], value) {
				time.Sleep(time.Millisecond * 100)
				continue
			} else {
				break
			}
		}
		assert.Equal(t, values[i], value)
	}
}
