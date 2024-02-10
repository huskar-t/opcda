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
	browse(t, browser)
	browser.MoveToRoot()
	browse(t, browser)
	err = browser.MoveTo([]string{"textual"})
	assert.NoError(t, err)
	assert.NoError(t, err)
	err = browser.ShowLeafs(false)
	assert.NoError(t, err)
	count := browser.GetCount()
	assert.Equal(t, 4, count)
	expectLeafs := []string{
		"color",
		"number",
		"random",
		"weekday",
	}

	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		assert.NoError(t, err)
		assert.Equal(t, expectLeafs[i], nextName)
		itermID, err := browser.GetItemID(nextName)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("textual.%s", nextName), itermID)
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
	assert.Equal(t, "textual", position)
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
	assert.Equal(t, 6, count)
	expectBranch := []string{"options", "numeric", "textual", "time", "enum", "storage"}
	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		assert.NoError(t, err)
		assert.Equal(t, expectBranch[i], nextName)
	}
	err = browser.MoveDown("textual")
	assert.NoError(t, err)
	err = browser.ShowLeafs(false)
	assert.NoError(t, err)
	count = browser.GetCount()
	assert.Equal(t, 4, count)
	expectLeafs := []string{
		"color",
		"number",
		"random",
		"weekday",
	}

	for i := 0; i < count; i++ {
		nextName, err := browser.Item(i)
		assert.NoError(t, err)
		assert.Equal(t, expectLeafs[i], nextName)
		itermID, err := browser.GetItemID(nextName)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("textual.%s", nextName), itermID)
	}
	err = browser.ShowBranches()
	assert.NoError(t, err)
	count = browser.GetCount()
	assert.Equal(t, 0, count)
}
