package opcda

import (
	"testing"
	"time"

	"github.com/huskar-t/opcda/com"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const TestProgID = "Matrikon.OPC.Simulation.1"
const TestHost = "localhost"
const TestServiceName = "MatrikonOPC Server for Simulation and Testing"
const TestBoolItem = "Bucket Brigade.Boolean"
const TestFloatItem = "Bucket Brigade.Real4"
const TestWriteItem = "Bucket Brigade.Int4"
const TestWriteErrorItem = "Write Error.Int4"
const TestReadErrorItem = "Write Only.Int4"

func TestMain(m *testing.M) {
	com.Initialize()
	com.Uninitialize()
	com.Initialize()
	defer com.Uninitialize()
	m.Run()
}
func TestServers(t *testing.T) {
	serverInfos, err := GetOPCServers(TestHost)
	assert.NoError(t, err)
	assert.Greater(t, len(serverInfos), 0)
	t.Log(serverInfos[0].ProgID)
}

func TestOpcServer_GetLocaleID(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	localID, err := server.GetLocaleID()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x800), localID)
}

func TestOpcServer_GetStartTime(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	startTime, err := server.GetStartTime()
	assert.NoError(t, err)
	assert.False(t, startTime.IsZero())
	t.Log("startTime", startTime)
}

func TestOpcServer_GetCurrentTime(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	currentTime, err := server.GetCurrentTime()
	assert.NoError(t, err)
	assert.False(t, currentTime.IsZero())
	t.Log("currentTime", currentTime)
}

func TestOpcServer_GetLastUpdateTime(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	lastUpdateTime, err := server.GetLastUpdateTime()
	assert.NoError(t, err)
	assert.False(t, lastUpdateTime.IsZero())
	t.Log("lastUpdateTime", lastUpdateTime)
}

func TestOpcServer_GetMajorVersion(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	majorVersion, err := server.GetMajorVersion()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, majorVersion, uint16(0))
	t.Log("majorVersion", majorVersion)
}

func TestOpcServer_GetMinorVersion(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	minorVersion, err := server.GetMinorVersion()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, minorVersion, uint16(0))
	t.Log("minorVersion", minorVersion)
}

func TestOpcServer_GetBuildNumber(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	buildNumber, err := server.GetBuildNumber()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, buildNumber, uint16(0))
	t.Log("buildNumber", buildNumber)
}

func TestOpcServer_GetVendorInfo(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	vendorInfo, err := server.GetVendorInfo()
	assert.NoError(t, err)
	assert.NotEmpty(t, vendorInfo)
	t.Log("vendorInfo", vendorInfo)
}

func TestOpcServer_GetServerState(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	status, err := server.GetServerState()
	assert.NoError(t, err)
	assert.Equal(t, OPC_STATUS_RUNNING, status)
	t.Log("status", status)
}

// SetLocaleID
func TestOpcServer_SetLocaleID(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	ids, err := server.QueryAvailableLocaleIDs()
	assert.NoError(t, err)
	assert.Greater(t, len(ids), 0)
	err = server.SetLocaleID(ids[0])
	assert.NoError(t, err)
	localID, err := server.GetLocaleID()
	assert.NoError(t, err)
	assert.Equal(t, ids[0], localID)
}

func TestOpcServer_GetBandwidth(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	bandwidth, err := server.GetBandwidth()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, bandwidth, uint32(0))
	t.Log("bandwidth", bandwidth)
}

func TestOpcServer_OPCGroups(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
}

// GetServerName
func TestOpcServer_GetServerName(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	serverName := server.GetServerName()
	assert.Equal(t, TestProgID, serverName)
}

// GetServerNode
func TestOpcServer_GetServerNode(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	serverNode := server.GetServerNode()
	assert.Equal(t, TestHost, serverNode)
}

// GetClientName
func TestOpcServer_GetClientName(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	err = server.SetClientName("test")
	assert.NoError(t, err)
	clientName := server.GetClientName()
	assert.Equal(t, "test", clientName)
}

func TestOpcServer_QueryAvailableProperties(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	ppPropertyIDs, ppDescriptions, ppvtDataTypes, err := server.QueryAvailableProperties(TestWriteItem)
	assert.NoError(t, err)
	assert.Greater(t, len(ppPropertyIDs), 0)
	assert.Greater(t, len(ppDescriptions), 0)
	assert.Greater(t, len(ppvtDataTypes), 0)
	t.Log(ppPropertyIDs, ppDescriptions, ppvtDataTypes)
}

// GetItemProperties
func TestOpcServer_GetItemProperties(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	ppPropertyIDs, ppDescriptions, ppvtDataTypes, err := server.QueryAvailableProperties(TestWriteItem)
	assert.NoError(t, err)
	assert.Greater(t, len(ppPropertyIDs), 0)
	assert.Greater(t, len(ppDescriptions), 0)
	assert.Greater(t, len(ppvtDataTypes), 0)
	t.Log(ppPropertyIDs, ppDescriptions, ppvtDataTypes)
	properties, errors, err := server.GetItemProperties(TestWriteItem, ppPropertyIDs)
	assert.NoError(t, err)
	assert.Greater(t, len(properties), 0)
	assert.Greater(t, len(errors), 0)
	assert.Equal(t, len(properties), len(errors))
	for i := 0; i < len(properties); i++ {
		assert.NoError(t, errors[i])
	}
	t.Log(properties)
}

// LookupItemIDs The simulator does not support
func TestOpcServer_LookupItemIDs(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	ppPropertyIDs, ppDescriptions, ppvtDataTypes, err := server.QueryAvailableProperties(TestBoolItem)
	assert.NoError(t, err)
	assert.Equal(t, len(ppPropertyIDs), 14)
	assert.Equal(t, len(ppDescriptions), 14)
	assert.Equal(t, len(ppvtDataTypes), 14)
	//t.Log(ppPropertyIDs, ppDescriptions, ppvtDataTypes)
	itemIDs, errors, err := server.LookupItemIDs(TestBoolItem, ppPropertyIDs)
	assert.NoError(t, err)
	assert.Equal(t, len(itemIDs), 14)
	assert.Equal(t, len(errors), 14)
	assert.Equal(t, len(itemIDs), len(errors))
	for i := 0; i < 9; i++ {
		//[0xc0040203]: The server does not recognise the passed property ID or the string was not recognized as an area name.
		assert.Error(t, errors[i])
	}
	expected := []string{
		"Triangle Waves.Boolean",
		"Square Waves.Boolean",
		"Saw-toothed Waves.Boolean",
		"Random.Boolean",
		"Bucket Brigade.Boolean",
	}
	for i := 9; i < 14; i++ {
		assert.Equal(t, expected[i-9], itemIDs[i])
	}
}

func TestOPCGroup_AddItems(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.groupName)
	items := group.OPCItems()
	itemList, errors, err := items.AddItems([]string{TestBoolItem, "x.x"})
	assert.NoError(t, err)
	hasError := false

	for i := 0; i < len(errors); i++ {
		if errors[i] != nil {
			hasError = true
		}
	}
	assert.Equal(t, 2, len(itemList))
	assert.NotNil(t, itemList[0])
	assert.Nil(t, itemList[1])
	assert.True(t, hasError)
}

func TestOPCGroup_AddItems_Success(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	assert.NotNil(t, groups)
	group, err := groups.Add("test")
	assert.NoError(t, err)
	assert.NotNil(t, group)
	assert.Equal(t, "test", group.groupName)
	items := group.OPCItems()
	itemList, errors, err := items.AddItems([]string{TestBoolItem, TestFloatItem})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(itemList))
	assert.Equal(t, 2, len(errors))
	for i := 0; i < 1; i++ {
		assert.NoError(t, errors[i])
	}
	assert.NoError(t, err)
	assert.NotNil(t, items)
	assert.Equal(t, TestBoolItem, itemList[0].tag)
	time.Sleep(time.Millisecond * 10)
	value, quality, ts, err := itemList[1].Read(OPC_DS_CACHE)
	assert.NoError(t, err)
	t.Log(value)
	t.Log(quality)
	t.Log(ts)
}

// Can be tested manually, but cannot be tested automatically
func TestOPCServer_RegisterServerShutDown(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()
	ch := make(chan string, 1)
	err = server.RegisterServerShutDown(ch)
	assert.NoError(t, err)
	done := make(chan struct{})
	go func() {
		manager, err := mgr.Connect()
		if err != nil {
			t.Fatalf("Failed to connect to service manager: %v", err)
		}
		defer manager.Disconnect()
		serviceObj, err := manager.OpenService(TestServiceName)
		if err != nil {
			t.Fatalf("Failed to open service %s: %v", TestServiceName, err)
		}
		defer serviceObj.Close()
		defer func() {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Second)
				status, err := serviceObj.Query()
				if err != nil {
					t.Fatalf("Failed to query service %s: %v", TestServiceName, err)
				}
				t.Log(status.State)
				if status.State == svc.Stopped {
					err = serviceObj.Start()
					if err != nil {
						t.Fatalf("Failed to start service %s: %v", TestServiceName, err)
					}
					break
				}
			}
			close(done)
		}()
		_, err = serviceObj.Control(svc.Stop)
		if err != nil {
			t.Fatalf("Failed to stop service %s: %v", TestServiceName, err)
		}

		t.Logf("Service %s stopped", TestServiceName)
	}()
	select {
	case reason := <-ch:
		t.Log(reason)
	}
	<-done
}
