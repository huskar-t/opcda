package opcda

import (
	"errors"
	"unsafe"

	"github.com/huskar-t/opcda/com"
)

type OPCBrowser struct {
	iBrowseServerAddressSpace *com.IOPCBrowseServerAddressSpace
	filter                    string
	dataType                  uint16
	accessRights              uint32
	names                     []string
	parent                    *OPCServer
}

func NewOPCBrowser(parent *OPCServer) (*OPCBrowser, error) {
	var iBrowseServerAddressSpace *com.IUnknown
	err := parent.iServer.QueryInterface(&com.IID_IOPCBrowseServerAddressSpace, unsafe.Pointer(&iBrowseServerAddressSpace))
	if err != nil {
		return nil, NewOPCWrapperError("query interface IOPCBrowseServerAddressSpace", err)
	}
	return &OPCBrowser{
		iBrowseServerAddressSpace: &com.IOPCBrowseServerAddressSpace{IUnknown: iBrowseServerAddressSpace},
		parent:                    parent,
		accessRights:              OPC_READABLE | OPC_WRITEABLE,
	}, nil
}

// GetFilter get the filter that applies to ShowBranches and ShowLeafs methods
func (b *OPCBrowser) GetFilter() string {
	return b.filter
}

// SetFilter set the filter that applies to ShowBranches and ShowLeafs methods
func (b *OPCBrowser) SetFilter(filter string) {
	b.filter = filter
}

// GetDataType get the requested data type that applies to ShowLeafs methods. This property defaults to
// com.VT_EMPTY, which means that any data type is acceptable.
func (b *OPCBrowser) GetDataType() uint16 {
	return b.dataType
}

// SetDataType set the requested data type that applies to ShowLeafs methods.
func (b *OPCBrowser) SetDataType(dataType uint16) {
	b.dataType = dataType
}

// GetAccessRights get the requested access rights that apply to the ShowLeafs methods
func (b *OPCBrowser) GetAccessRights() uint32 {
	return b.accessRights
}

// SetAccessRights set the requested access rights that apply to the ShowLeafs methods
func (b *OPCBrowser) SetAccessRights(accessRights uint32) error {
	if accessRights&OPC_READABLE == 0 && accessRights&OPC_WRITEABLE == 0 {
		return errors.New("accessRights must be OPC_READABLE or OPC_WRITEABLE")
	}
	b.accessRights = accessRights
	return nil
}

// GetCurrentPosition Returns the current position in the tree
func (b *OPCBrowser) GetCurrentPosition() (string, error) {
	id, err := b.iBrowseServerAddressSpace.GetItemID("")
	return id, err
}

// GetOrganization Returns either OPCHierarchical or OPCFlat.
func (b *OPCBrowser) GetOrganization() (com.OPCNAMESPACETYPE, error) {
	return b.iBrowseServerAddressSpace.QueryOrganization()
}

// GetCount Required property for collections
func (b *OPCBrowser) GetCount() int {
	return len(b.names)
}

// Item returns the name of the item at the specified index. index is 0-based.
func (b *OPCBrowser) Item(index int) (string, error) {
	if index < 0 || index >= len(b.names) {
		return "", errors.New("index out of range")
	}
	return b.names[index], nil
}

// ShowBranches Fills the collection with names of the branches at the current browse position.
func (b *OPCBrowser) ShowBranches() error {
	b.names = nil
	var err error
	b.names, err = b.iBrowseServerAddressSpace.BrowseOPCItemIDs(OPC_BRANCH, b.filter, b.dataType, b.accessRights)
	return err
}

// ShowLeafs Fills the collection with the names of the leafs at the current browse position
func (b *OPCBrowser) ShowLeafs(flat bool) error {
	b.names = nil
	var err error
	browseType := OPC_LEAF
	if flat {
		browseType = OPC_FLAT
	}
	b.names, err = b.iBrowseServerAddressSpace.BrowseOPCItemIDs(browseType, b.filter, b.dataType, b.accessRights)
	return err
}

// MoveUp Move up one level in the tree.
func (b *OPCBrowser) MoveUp() error {
	return b.iBrowseServerAddressSpace.ChangeBrowsePosition(OPC_BROWSE_UP, "")
}

// MoveToRoot Move up to the first level in the tree.
func (b *OPCBrowser) MoveToRoot() {
	for {
		err := b.iBrowseServerAddressSpace.ChangeBrowsePosition(OPC_BROWSE_UP, "")
		if err != nil {
			break
		}
	}
}

// MoveDown Move down into this branch.
func (b *OPCBrowser) MoveDown(name string) error {
	return b.iBrowseServerAddressSpace.ChangeBrowsePosition(OPC_BROWSE_DOWN, name)
}

// MoveTo Move to an absolute position.
func (b *OPCBrowser) MoveTo(branches []string) error {
	b.MoveToRoot()
	for _, branch := range branches {
		err := b.MoveDown(branch)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetItemID Given a name, returns a valid ItemID that can be passed to OPCItems Add method.
func (b *OPCBrowser) GetItemID(leaf string) (string, error) {
	return b.iBrowseServerAddressSpace.GetItemID(leaf)
}

// Release release the OPCBrowser
func (b *OPCBrowser) Release() {
	b.iBrowseServerAddressSpace.Release()
}
