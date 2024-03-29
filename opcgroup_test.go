package opcda

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOPCGroup_SetName(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	err = group.SetName("test2")
	assert.NoError(t, err)
	assert.Equal(t, "test2", group.GetName())
	assert.Equal(t, groups, group.GetParent())
}

func TestOPCGroup_IsActive(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	assert.Equal(t, true, group.GetIsActive())
	err = group.SetIsActive(false)
	assert.NoError(t, err)
	assert.Equal(t, false, group.GetIsActive())
	err = group.SetIsActive(true)
	assert.NoError(t, err)
}

func TestOPCGroup_Handle(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	assert.Equal(t, uint32(1), group.GetClientHandle())
	err = group.SetClientHandle(2)
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), group.GetClientHandle())
	serverHandle := group.GetServerHandle()
	assert.Greater(t, serverHandle, uint32(0))
}

func TestOPCGroup_LocaleID(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	localID, err := group.GetLocaleID()
	assert.NoError(t, err)
	assert.Greater(t, localID, uint32(0))
	err = group.SetLocaleID(1024)
	assert.NoError(t, err)
	localID, err = group.GetLocaleID()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1024), localID)
}

func TestOPCGroup_TimeBias(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	timeBias, err := group.GetTimeBias()
	assert.NoError(t, err)
	assert.Equal(t, int32(0), timeBias)
	err = group.SetTimeBias(1024)
	assert.NoError(t, err)
	timeBias, err = group.GetTimeBias()
	assert.NoError(t, err)
	assert.Equal(t, int32(1024), timeBias)
}

func TestOPCGroup_Deadband(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	deadband, err := group.GetDeadband()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), deadband)
	err = group.SetDeadband(0)
	assert.NoError(t, err)
	deadband, err = group.GetDeadband()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), deadband)
}

func TestOPCGroup_UpdateRate(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())
	updateRate, err := group.GetUpdateRate()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1000), updateRate)
	err = group.SetUpdateRate(5000)
	assert.NoError(t, err)
	updateRate, err = group.GetUpdateRate()
	assert.NoError(t, err)
	assert.Equal(t, uint32(5000), updateRate)
}

func TestOPCGroup_Items(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.GetName())

	items := group.items
	assert.NotNil(t, items)
	items2 := group.OPCItems()
	assert.Equal(t, items, items2)
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	err = item.SetIsActive(true)
	assert.NoError(t, err)
}

func TestOPCGroup_SyncRead(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_read")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test_group_read", group.GetName())

	items := group.items
	assert.NotNil(t, items)
	items2 := group.OPCItems()
	assert.Equal(t, items, items2)
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	match := false
	item.SetIsActive(true)
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		status, err := group.SyncRead(OPC_DS_CACHE, []uint32{item.GetServerHandle()})
		assert.NoError(t, err)
		t.Log(status[0])
		assert.Equal(t, 1, len(status))
		if status[0].Quality != uint16(192) {
			continue
		}
		value, quality, ts, err := item.Read(OPC_DS_CACHE)
		assert.NoError(t, err)
		assert.Equal(t, uint16(192), quality)
		if status[0].Timestamp != ts {
			continue
		}
		assert.Equal(t, status[0].Value, value)
		match = true
		break
	}
	assert.True(t, match)
}

func TestOPCGroup_ReadError(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_read_error")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.items
	assert.NotNil(t, items)
	items2 := group.OPCItems()
	assert.Equal(t, items, items2)
	item, err := items.AddItem(TestReadErrorItem)
	assert.NoError(t, err)
	_, err = group.SyncRead(OPC_DS_CACHE, []uint32{item.GetServerHandle()})
	assert.Error(t, err)
}

func TestOPCGroup_SyncWrite(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_write")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test_group_write", group.GetName())

	items := group.items
	assert.NotNil(t, items)
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	ch := make(chan *DataChangeCallBackData, 1)
	err = group.RegisterDataChange(ch)
	assert.NoError(t, err)
	errs, err := group.SyncWrite([]uint32{item.GetServerHandle()}, []interface{}{int32(11)})
	assert.NoError(t, err)
	for _, err := range errs {
		assert.NoError(t, err)
	}
	var value interface{}
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		value, _, _, err = item.Read(OPC_DS_CACHE)
		assert.NoError(t, err)
		if value == int32(11) {
			break
		}
	}
	assert.Equal(t, int32(11), value)
	timeout := time.NewTimer(time.Second * 5)
	defer timeout.Stop()
	for {
		select {
		case v := <-ch:
			t.Log(v)
			assert.Equal(t, 1, len(v.Qualities))
			return
		case <-timeout.C:
			t.Fatal("timeout")
		}
	}
}

func TestOPCGroup_WriteError(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_write_error")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	items := group.items
	assert.NotNil(t, items)
	item, err := items.AddItem(TestWriteErrorItem)
	assert.NoError(t, err)
	ch := make(chan *DataChangeCallBackData, 1)
	err = group.RegisterDataChange(ch)
	assert.NoError(t, err)
	errs, err := group.SyncWrite([]uint32{item.GetServerHandle()}, []interface{}{int32(11)})
	assert.NoError(t, err)
	for _, err := range errs {
		assert.Error(t, err)
	}
}

func TestOPCGroup_AsyncRead(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_async_read")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	ch := make(chan *ReadCompleteCallBackData, 1)
	out := make(chan *ReadCompleteCallBackData, 1)
	go func() {
		select {
		case data := <-ch:
			out <- data
		}
	}()
	err = group.RegisterReadComplete(ch)
	assert.NoError(t, err)
	items := group.items
	item, err := items.AddItem(TestBoolItem)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	cancelID, errs, err := group.AsyncRead([]uint32{item.serverHandle}, 100)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(errs))
	assert.Nil(t, errs[0])
	t.Log(cancelID)
	timeout := time.NewTimer(time.Second)
	select {
	case data := <-out:
		t.Log(data)
		assert.Equal(t, group.GetClientHandle(), data.GroupHandle)
		assert.Equal(t, uint32(100), data.TransID)
		assert.Equal(t, 1, len(data.ItemClientHandles))
		assert.Equal(t, item.clientHandle, data.ItemClientHandles[0])
	case <-timeout.C:
		t.Fatal("timeout")
	}
}

func TestOPCGroup_AsyncWrite(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_async_write")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	ch := make(chan *WriteCompleteCallBackData, 1)
	out := make(chan *WriteCompleteCallBackData, 1)
	go func() {
		select {
		case data := <-ch:
			out <- data
		}
	}()
	err = group.RegisterWriteComplete(ch)
	assert.NoError(t, err)
	items := group.items
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	cancelID, errs, err := group.AsyncWrite([]uint32{item.serverHandle}, []interface{}{int32(14)}, 100)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(errs))
	assert.Nil(t, errs[0])
	t.Log(cancelID)
	timeout := time.NewTimer(time.Second)
	select {
	case data := <-out:
		t.Log(data)
		assert.Equal(t, group.GetClientHandle(), data.GroupHandle)
		assert.Equal(t, uint32(100), data.TransID)
		assert.Equal(t, 1, len(data.ItemClientHandles))
		assert.Equal(t, item.clientHandle, data.ItemClientHandles[0])
		assert.Equal(t, 1, len(data.Errors))
		assert.Nil(t, data.Errors[0])
	case <-timeout.C:
		t.Fatal("timeout")
	}
}

func TestOPCGroup_AsyncRefresh(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_async_refresh")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	ch := make(chan *DataChangeCallBackData, 1)
	out := make(chan *DataChangeCallBackData, 1)
	go func() {
		select {
		case data := <-ch:
			t.Log(data)
			out <- data
		}
	}()
	err = group.RegisterDataChange(ch)
	assert.NoError(t, err)
	items := group.items
	item, err := items.AddItem(TestBoolItem)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	cancelID, err := group.AsyncRefresh(OPC_DS_CACHE, 100)
	assert.NoError(t, err)
	t.Log(cancelID)
	timeout := time.NewTimer(time.Second)
	select {
	case data := <-out:
		t.Log(data)
		assert.Equal(t, group.GetClientHandle(), data.GroupHandle)
		assert.Equal(t, 1, len(data.ItemClientHandles))
		assert.Equal(t, item.clientHandle, data.ItemClientHandles[0])
	case <-timeout.C:
		t.Fatal("timeout")
	}
}

func TestOPCGroup_AsyncCancel(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test_group_async_cancel")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	ch := make(chan *CancelCompleteCallBackData, 1)
	out := make(chan *CancelCompleteCallBackData, 1)
	go func() {
		select {
		case data := <-ch:
			out <- data
		}
	}()
	err = group.RegisterCancelComplete(ch)
	assert.NoError(t, err)
	items := group.items
	item, err := items.AddItem(TestWriteItem)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	for i := 0; i < 300; i++ {
		cancelID, errs, err := group.AsyncWrite([]uint32{item.serverHandle}, []interface{}{int32(14)}, 100)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(errs))
		assert.Nil(t, errs[0])
		group.AsyncCancel(cancelID)
	}
	timeout := time.NewTimer(time.Second * 5)
	select {
	case data := <-out:
		t.Log(data)
		assert.Equal(t, group.GetClientHandle(), data.GroupHandle)
		assert.Equal(t, uint32(100), data.TransID)
	case <-timeout.C:
		t.Fatal("timeout")
	}
}
