package opcda

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/huskar-t/opcda/com"

	"golang.org/x/sys/windows"
)

type OPCItems struct {
	itemMgt                  *com.IOPCItemMgt
	iCommon                  *com.IOPCCommon
	parent                   *OPCGroup
	itemID                   uint32
	defaultRequestedDataType com.VT
	defaultAccessPath        string
	defaultActive            bool
	items                    []*OPCItem
	sync.RWMutex
}

func NewOPCItems(
	parent *OPCGroup,
	itemMgt *com.IOPCItemMgt,
	iCommon *com.IOPCCommon,
) *OPCItems {
	return &OPCItems{
		parent:                   parent,
		itemMgt:                  itemMgt,
		defaultRequestedDataType: com.VT_EMPTY,
		defaultAccessPath:        "",
		defaultActive:            true,
		iCommon:                  iCommon,
	}
}

// GetParent Returns reference to the parent OPCGroup object
func (is *OPCItems) GetParent() *OPCGroup {
	return is.parent
}

// GetDefaultRequestedDataType get the requested data type that will be used in calls to Add
func (is *OPCItems) GetDefaultRequestedDataType() com.VT {
	return is.defaultRequestedDataType
}

// SetDefaultRequestedDataType set the requested data type that will be used in calls to Add
func (is *OPCItems) SetDefaultRequestedDataType(defaultRequestedDataType com.VT) {
	is.defaultRequestedDataType = defaultRequestedDataType
}

// GetDefaultAccessPath get the default AccessPath that will be used in calls to Add
func (is *OPCItems) GetDefaultAccessPath() string {
	return is.defaultAccessPath
}

// SetDefaultAccessPath set the default AccessPath that will be used in calls to Add
func (is *OPCItems) SetDefaultAccessPath(defaultAccessPath string) {
	is.Lock()
	defer is.Unlock()
	is.defaultAccessPath = defaultAccessPath
}

// GetDefaultActive get the default active state for OPCItems created using Items.Add
func (is *OPCItems) GetDefaultActive() bool {
	return is.defaultActive
}

// SetDefaultActive set the default active state for OPCItems created using Items.Add
func (is *OPCItems) SetDefaultActive(defaultActive bool) {
	is.defaultActive = defaultActive
}

// GetCount get the number of items in the collection
func (is *OPCItems) GetCount() int {
	return len(is.items)
}

// Item get the item by index
func (is *OPCItems) Item(index int32) (*OPCItem, error) {
	is.RLock()
	defer is.RUnlock()
	if index < 0 || index >= int32(len(is.items)) {
		return nil, errors.New("index out of range")
	}
	return is.items[index], nil
}

// ItemByName get the item by name
func (is *OPCItems) ItemByName(name string) (*OPCItem, error) {
	is.RLock()
	defer is.RUnlock()
	for _, v := range is.items {
		if v.tag == name {
			return v, nil
		}
	}
	return nil, errors.New("not found")
}

// GetOPCItem returns the OPCItem by serverHandle
func (is *OPCItems) GetOPCItem(serverHandle uint32) (*OPCItem, error) {
	is.RLock()
	defer is.RUnlock()
	for _, v := range is.items {
		if v.serverHandle == serverHandle {
			return v, nil
		}
	}
	return nil, errors.New("not found")
}

// AddItem adds an item to the group.
func (is *OPCItems) AddItem(tag string) (*OPCItem, error) {
	items, errs, err := is.AddItems([]string{tag})
	if err != nil {
		return nil, err
	}
	if errs[0] != nil {
		return nil, errs[0]
	}
	return items[0], nil
}

// AddItems adds items to the group.
func (is *OPCItems) AddItems(tags []string) ([]*OPCItem, []error, error) {
	is.Lock()
	defer is.Unlock()
	accessPath := is.defaultAccessPath
	active := is.defaultActive
	dt := is.defaultRequestedDataType
	items := is.createDefinitions(tags, accessPath, active, dt)
	results, errs, err := is.itemMgt.AddItems(items)
	if err != nil {
		return nil, nil, err
	}
	var resultErrors = make([]error, len(tags))
	var opcItems = make([]*OPCItem, len(tags))
	for j := 0; j < len(tags); j++ {
		if errs[j] < 0 {
			resultErrors[j] = is.getError(errs[j])
		} else {
			item := NewOPCItem(is, tags[j], results[j], items[j].HClient, accessPath, active)
			opcItems[j] = item
			is.items = append(is.items, item)
		}
	}
	return opcItems, resultErrors, nil
}

// Remove Removes an OPCItem
func (is *OPCItems) Remove(serverHandles []uint32) {
	is.Lock()
	defer is.Unlock()
	for _, v := range serverHandles {
		for j, w := range is.items {
			if w.serverHandle == v {
				w.Release()
				w.itemMgt.RemoveItems([]uint32{v})
				is.items = append(is.items[:j], is.items[j+1:]...)
				break
			}
		}
	}
}

// Validate Determines if one or more OPCItems could be successfully created via the Add method (but does not add them).
func (is *OPCItems) Validate(tags []string, requestedDataTypes *[]com.VT, accessPaths *[]string) ([]error, error) {
	var definitions []com.TagOPCITEMDEF
	for i, v := range tags {
		cHandle := atomic.AddUint32(&is.itemID, 1)
		item := com.TagOPCITEMDEF{
			SzAccessPath: windows.StringToUTF16Ptr(""),
			SzItemID:     windows.StringToUTF16Ptr(v),
			BActive:      com.BoolToBOOL(false),
			HClient:      cHandle,
			DwBlobSize:   0,
			PBlob:        nil,
			VtRequested:  uint16(is.defaultRequestedDataType),
		}
		if requestedDataTypes != nil {
			item.VtRequested = uint16((*requestedDataTypes)[i])
		}
		if accessPaths != nil {
			item.SzAccessPath = windows.StringToUTF16Ptr((*accessPaths)[i])
		}
		definitions = append(definitions, item)
	}
	_, errs, err := is.itemMgt.ValidateItems(definitions, false)
	if err != nil {
		return nil, err
	}
	var resultErrors = make([]error, len(errs))
	for j := 0; j < len(errs); j++ {
		if errs[j] < 0 {
			resultErrors[j] = is.getError(errs[j])
		}
	}
	return resultErrors, nil
}

// SetActive Allows Activation and deactivation of individual OPCItem’s in the OPCItems Collection
func (is *OPCItems) SetActive(serverHandles []uint32, active bool) []error {
	resultErrors := make([]error, len(serverHandles))
	for i, handle := range serverHandles {
		item, err := is.GetOPCItem(handle)
		if err != nil {
			resultErrors[i] = err
			continue
		}
		err = item.SetIsActive(active)
		if err != nil {
			resultErrors[i] = err
		}
	}
	return resultErrors
}

// SetClientHandles Changes the client handles or one or more Items in a Group.
func (is *OPCItems) SetClientHandles(serverHandles []uint32, clientHandles []uint32) []error {
	resultErrors := make([]error, len(serverHandles))
	for i, handle := range serverHandles {
		item, err := is.GetOPCItem(handle)
		if err != nil {
			resultErrors[i] = err
			continue
		}
		err = item.SetClientHandle(clientHandles[i])
		if err != nil {
			resultErrors[i] = err
		}
	}
	return resultErrors
}

// SetDataTypes Changes the requested data type for one or more Items
func (is *OPCItems) SetDataTypes(serverHandles []uint32, requestedDataTypes []com.VT) []error {
	resultErrors := make([]error, len(serverHandles))
	for i, handle := range serverHandles {
		item, err := is.GetOPCItem(handle)
		if err != nil {
			resultErrors[i] = err
			continue
		}
		err = item.SetRequestedDataType(requestedDataTypes[i])
		if err != nil {
			resultErrors[i] = err
		}
	}
	return resultErrors
}

// Release Releases the OPCItems collection and all associated resources.
func (is *OPCItems) Release() {
	for _, item := range is.items {
		item.Release()
	}
	is.itemMgt.Release()
}

func (is *OPCItems) createDefinitions(tags []string, accessPath string, active bool, requestedDataType com.VT) []com.TagOPCITEMDEF {
	var definitions []com.TagOPCITEMDEF
	for _, v := range tags {
		cHandle := atomic.AddUint32(&is.itemID, 1)
		definitions = append(definitions, com.TagOPCITEMDEF{
			SzAccessPath: windows.StringToUTF16Ptr(accessPath),
			SzItemID:     windows.StringToUTF16Ptr(v),
			BActive:      com.BoolToBOOL(active),
			HClient:      cHandle,
			DwBlobSize:   0,
			PBlob:        nil,
			VtRequested:  uint16(requestedDataType),
		})
	}
	return definitions
}

func (is *OPCItems) getError(errorCode int32) error {
	errStr, _ := is.iCommon.GetErrorString(uint32(errorCode))
	return &OPCError{
		ErrorCode:    errorCode,
		ErrorMessage: errStr,
	}
}
