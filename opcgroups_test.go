package opcda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// parent
func TestOPCGroups_Parent(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	parent := groups.GetParent()
	assert.NoError(t, err)
	assert.Equal(t, server, parent)
}

func TestOPCGroups_GetDefaultGroupIsActive(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	isActive := groups.GetDefaultGroupIsActive()
	assert.NoError(t, err)
	assert.Equal(t, true, isActive)
}

func TestOPCGroups_SetDefaultGroupIsActive(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	groups.SetDefaultGroupIsActive(false)
	assert.NoError(t, err)
	isActive := groups.GetDefaultGroupIsActive()
	assert.NoError(t, err)
	assert.Equal(t, false, isActive)
}

// GetDefaultGroupUpdateRate
func TestOPCGroups_GetDefaultGroupUpdateRate(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	updateRate := groups.GetDefaultGroupUpdateRate()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1000), updateRate)
}

// SetDefaultGroupUpdateRate
func TestOPCGroups_SetDefaultGroupUpdateRate(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	groups.SetDefaultGroupUpdateRate(2000)
	assert.NoError(t, err)
	updateRate := groups.GetDefaultGroupUpdateRate()
	assert.NoError(t, err)
	assert.Equal(t, uint32(2000), updateRate)
}

// GetDefaultGroupDeadband
func TestOPCGroups_GetDefaultGroupDeadband(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	deadband := groups.GetDefaultGroupDeadband()
	assert.NoError(t, err)
	assert.Equal(t, float32(0.0), deadband)
}

// SetDefaultGroupDeadband
func TestOPCGroups_SetDefaultGroupDeadband(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	groups.SetDefaultGroupDeadband(1.0)
	assert.NoError(t, err)
	deadband := groups.GetDefaultGroupDeadband()
	assert.NoError(t, err)
	assert.Equal(t, float32(1.0), deadband)
}

// GetDefaultGroupLocaleID
func TestOPCGroups_GetDefaultGroupLocaleID(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	localeID := groups.GetDefaultGroupLocaleID()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0400), localeID)
}

// SetDefaultGroupLocaleID
func TestOPCGroups_SetDefaultGroupLocaleID(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	groups.SetDefaultGroupLocaleID(0x0401)
	assert.NoError(t, err)
	localeID := groups.GetDefaultGroupLocaleID()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x0401), localeID)
}

// GetDefaultGroupTimeBias
func TestOPCGroups_GetDefaultGroupTimeBias(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()

	assert.NotNil(t, groups)
	timeBias := groups.GetDefaultGroupTimeBias()
	assert.NoError(t, err)
	assert.Equal(t, int32(0), timeBias)
}

// SetDefaultGroupTimeBias
func TestOPCGroups_SetDefaultGroupTimeBias(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()

	assert.NotNil(t, groups)
	groups.SetDefaultGroupTimeBias(1)
	assert.NoError(t, err)
	timeBias := groups.GetDefaultGroupTimeBias()
	assert.NoError(t, err)
	assert.Equal(t, int32(1), timeBias)
}

// GetCount
func TestOPCGroups_GetCount(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()

	assert.NotNil(t, groups)
	count := groups.GetCount()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestOPCGroups_AddGroup(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.groupName)
	count := groups.GetCount()
	assert.Equal(t, 1, count)
	g, err := groups.Item(0)
	assert.NoError(t, err)
	assert.Equal(t, group, g)
	_, err = groups.Item(1)
	assert.Error(t, err)
	g2, err := groups.ItemByName("test")
	assert.NoError(t, err)
	assert.Equal(t, group, g2)
	_, err = groups.ItemByName("test2")
	assert.Error(t, err)
	g3, err := groups.GetOPCGroupByName("test")
	assert.NoError(t, err)
	assert.Equal(t, group, g3)
	_, err = groups.GetOPCGroupByName("test2")
	assert.Error(t, err)
	g4, err := groups.GetOPCGroup(group.GetServerHandle())
	assert.NoError(t, err)
	assert.Equal(t, group, g4)
	_, err = groups.GetOPCGroup(group.GetServerHandle() + 1)
	assert.Error(t, err)
	err = groups.Remove(group.GetServerHandle())
	assert.NoError(t, err)
	err = groups.Remove(group.GetServerHandle())
	assert.Error(t, err)
	_, err = groups.Add("test")
	assert.NoError(t, err)
	assert.Equal(t, 1, groups.GetCount())
	err = groups.RemoveByName("test")
	assert.NoError(t, err)
	assert.Equal(t, 0, groups.GetCount())
	err = groups.RemoveByName("test")
	assert.Error(t, err)
	_, err = groups.Add("test")
	assert.NoError(t, err)
	assert.Equal(t, 1, groups.GetCount())
	err = groups.RemoveAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, groups.GetCount())
}
