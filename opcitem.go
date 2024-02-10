package opcda

import (
	"errors"
	"time"

	"github.com/huskar-t/opcda/com"
)

type OPCItem struct {
	itemMgt           *com.IOPCItemMgt
	syncIO            *com.IOPCSyncIO
	iCommon           *com.IOPCCommon
	value             interface{}
	quality           uint16
	timestamp         time.Time
	serverHandle      uint32
	clientHandle      uint32
	tag               string
	accessPath        string
	accessRights      uint32
	isActive          bool
	requestedDataType com.VT
	nativeDataType    com.VT
	parent            *OPCItems
}

// GetParent Returns reference to the parent OPCItems object.
func (i *OPCItem) GetParent() *OPCItems {
	return i.parent
}

// GetClientHandle get the client handle for the item.
func (i *OPCItem) GetClientHandle() uint32 {
	return i.clientHandle
}

// SetClientHandle set the client handle for the item.
func (i *OPCItem) SetClientHandle(clientHandle uint32) error {
	errs, err := i.itemMgt.SetClientHandles([]uint32{i.serverHandle}, []uint32{clientHandle})
	if err != nil {
		return err
	}
	if errs[0] != 0 {
		return i.getError(errs[0])
	}
	i.clientHandle = clientHandle
	return nil
}

// GetServerHandle get the server handle for the item.
func (i *OPCItem) GetServerHandle() uint32 {
	return i.serverHandle
}

// GetAccessPath get the access path for the item.
func (i *OPCItem) GetAccessPath() string {
	return i.accessPath
}

// GetAccessRights get the access rights for the item.
func (i *OPCItem) GetAccessRights() uint32 {
	return i.accessRights
}

// GetItemID get the item ID for the item.
func (i *OPCItem) GetItemID() string {
	return i.tag
}

// GetIsActive get the active state for the item.
func (i *OPCItem) GetIsActive() bool {
	return i.isActive
}

// GetRequestedDataType get the requested data type for the item.
func (i *OPCItem) GetRequestedDataType() com.VT {
	return i.requestedDataType
}

// SetRequestedDataType set the requested data type for the item.
func (i *OPCItem) SetRequestedDataType(requestedDataType com.VT) error {
	errs, err := i.itemMgt.SetDatatypes([]uint32{i.serverHandle}, []com.VT{requestedDataType})
	if err != nil {
		return err
	}
	if errs[0] != 0 {
		return i.getError(errs[0])
	}
	i.requestedDataType = requestedDataType
	return nil
}

// SetIsActive set the active state for the item.
func (i *OPCItem) SetIsActive(isActive bool) error {
	errs, err := i.itemMgt.SetActiveState([]uint32{i.serverHandle}, isActive)
	if err != nil {
		return err
	}
	if errs[0] < 0 {
		return i.getError(errs[0])
	}
	i.isActive = isActive
	return nil
}

// GetValue Returns the latest value read from the server
func (i *OPCItem) GetValue() interface{} {
	return i.value
}

// GetQuality Returns the latest quality read from the server
func (i *OPCItem) GetQuality() uint16 {
	return i.quality
}

// GetTimestamp Returns the latest timestamp read from the server
func (i *OPCItem) GetTimestamp() time.Time {
	return i.timestamp
}

// GetCanonicalDataType Returns the canonical data type for the item.
func (i *OPCItem) GetCanonicalDataType() com.VT {
	return i.nativeDataType
}

// GetEUType Returns the EU type for the item.
func (i *OPCItem) GetEUType() (int, error) {
	data, errs, err := i.parent.parent.parent.parent.GetItemProperties(i.tag, []uint32{7})
	if err != nil {
		return 0, err
	}
	if errs[0] != nil {
		return 0, errs[0]
	}
	return (int)(data[0].(int32)), nil
}

// GetEUInfo Returns the EU info for the item.
func (i *OPCItem) GetEUInfo() (interface{}, error) {
	euType, err := i.GetEUType()
	if err != nil {
		return nil, err
	}
	if euType == 0 {
		return nil, nil
	}
	if euType > 2 {
		return nil, errors.New("not valid")
	}
	data, errs, err := i.parent.parent.parent.parent.GetItemProperties(i.tag, []uint32{8})
	if err != nil {
		return nil, err
	}
	if errs[0] != nil {
		return nil, errs[0]
	}
	return data[0], nil
}

func NewOPCItem(
	parent *OPCItems,
	tag string,
	result com.TagOPCITEMRESULTStruct,
	clientHandle uint32,
	accessPath string,
	isActive bool,
) *OPCItem {
	return &OPCItem{
		itemMgt:        parent.itemMgt,
		syncIO:         parent.parent.syncIO,
		iCommon:        parent.iCommon,
		parent:         parent,
		tag:            tag,
		accessPath:     accessPath,
		serverHandle:   result.Server,
		clientHandle:   clientHandle,
		accessRights:   result.AccessRights,
		nativeDataType: com.VT(result.NativeType),
		isActive:       isActive,
	}
}

// Read makes a blocking call to read this item from the server.
func (i *OPCItem) Read(source com.OPCDATASOURCE) (value interface{}, quality uint16, timestamp time.Time, err error) {
	values, errs, err := i.syncIO.Read(source, []uint32{i.serverHandle})
	if err != nil {
		return nil, 0, time.Time{}, err
	}
	if errs[0] < 0 {
		err = i.getError(errs[0])
		return
	}
	value = values[0].Value
	quality = values[0].Quality
	timestamp = values[0].Timestamp
	i.value = values[0].Value
	i.quality = values[0].Quality
	i.timestamp = values[0].Timestamp
	return
}

// Write makes a blocking call to write this value to the server
func (i *OPCItem) Write(value interface{}) error {
	variant, err := com.NewVariant(value)
	if err != nil {
		return err
	}
	defer variant.Clear()
	errList, err := i.syncIO.Write([]uint32{i.serverHandle}, []com.VARIANT{*variant.Variant})
	if err != nil {
		return err
	}
	if errList[0] < 0 {
		return i.getError(errList[0])
	}
	return nil
}

func (i *OPCItem) getError(errorCode int32) error {
	errStr, _ := i.iCommon.GetErrorString(uint32(errorCode))
	return &OPCError{
		ErrorCode:    errorCode,
		ErrorMessage: errStr,
	}
}

// Release Releases the OPCItem object
func (i *OPCItem) Release() {
}
