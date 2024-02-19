package opcda

import (
	"fmt"
	"testing"

	"github.com/huskar-t/opcda/com"

	"github.com/stretchr/testify/assert"
)

func TestOPCBrowser(t *testing.T) {
	server, err := Connect(TestProgID, TestHost)
	assert.NoError(t, err)
	defer server.Disconnect()

	browser, err := server.CreateBrowser()
	if err != nil {
		t.Fatal(err)
	}
	defer browser.Release()

	browser.MoveToRoot()
	browse(t, browser)
	err = browser.MoveUp()
	assert.NoError(t, err)
	err = browser.MoveUp()
	assert.NoError(t, err)
	browse(t, browser)
	browser.MoveToRoot()
	browse(t, browser)
	err = browser.MoveTo([]string{"Simulation Items", "Bucket Brigade"})
	assert.NoError(t, err)
	err = browser.ShowLeafs(false)
	assert.NoError(t, err)
	count := browser.GetCount()
	assert.Equal(t, 14, count)
	expectLeafs := []string{
		"ArrayOfReal8",
		"ArrayOfString",
		"Boolean",
		"Int1",
		"Int2",
		"Int4",
		"Money",
		"Real4",
		"Real8",
		"String",
		"Time",
		"UInt1",
		"UInt2",
		"UInt4",
	}

	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		assert.NoError(t, err)
		assert.Equal(t, expectLeafs[i], nextName)
		itermID, err := browser.GetItemID(nextName)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Bucket Brigade.%s", nextName), itermID)
	}
	browser.SetFilter("")
	filter := browser.GetFilter()
	assert.Equal(t, "", filter)
	browser.SetDataType(uint16(com.VT_BOOL))
	dataType := browser.GetDataType()
	assert.Equal(t, uint16(com.VT_BOOL), dataType)
	err = browser.SetAccessRights(OPC_READABLE)
	assert.NoError(t, err)
	accessRights := browser.GetAccessRights()
	assert.Equal(t, OPC_READABLE, accessRights)
	err = browser.SetAccessRights(4)
	assert.Error(t, err)
	position, err := browser.GetCurrentPosition()
	assert.NoError(t, err)
	assert.Equal(t, "Bucket Brigade", position)
	organization, err := browser.GetOrganization()
	assert.NoError(t, err)
	assert.Equal(t, OPC_NS_HIERARCHIAL, organization)
	browser.MoveToRoot()
	position, err = browser.GetCurrentPosition()
	assert.NoError(t, err)
	assert.Equal(t, "", position)
}

func browse(t *testing.T, browser *OPCBrowser) {
	err := browser.ShowBranches()
	assert.NoError(t, err)
	count := browser.GetCount()
	assert.Equal(t, 2, count)
	expectBranch := []string{"Simulation Items", "Configured Aliases"}
	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		assert.NoError(t, err)
		assert.Equal(t, expectBranch[i], nextName)
	}
	err = browser.MoveDown("Simulation Items")
	assert.NoError(t, err)
	err = browser.ShowBranches()
	assert.NoError(t, err)
	count = browser.GetCount()
	assert.Equal(t, 8, count)
	expectBranch = []string{
		"Bucket Brigade",
		"Random",
		"Read Error",
		"Saw-toothed Waves",
		"Square Waves",
		"Triangle Waves",
		"Write Error",
		"Write Only",
	}
	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		assert.NoError(t, err)
		assert.Equal(t, expectBranch[i], nextName)
	}
	err = browser.MoveDown("Bucket Brigade")
	assert.NoError(t, err)
	err = browser.ShowLeafs(false)
	assert.NoError(t, err)
	count = browser.GetCount()
	assert.Equal(t, 14, count)
	expectLeafs := []string{
		"ArrayOfReal8",
		"ArrayOfString",
		"Boolean",
		"Int1",
		"Int2",
		"Int4",
		"Money",
		"Real4",
		"Real8",
		"String",
		"Time",
		"UInt1",
		"UInt2",
		"UInt4",
	}

	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		assert.NoError(t, err)
		assert.Equal(t, expectLeafs[i], nextName)
		itermID, err := browser.GetItemID(nextName)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Bucket Brigade.%s", nextName), itermID)
	}
	err = browser.ShowBranches()
	assert.NoError(t, err)
	count = browser.GetCount()
	assert.Equal(t, 0, count)
}
