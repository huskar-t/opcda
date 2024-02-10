package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IOPCBrowseServerAddressSpace = windows.GUID{
	Data1: 0x39c13a4f,
	Data2: 0x011e,
	Data3: 0x11d0,
	Data4: [8]byte{0x96, 0x75, 0x00, 0x20, 0xaf, 0xd8, 0xad, 0xb3},
}

type IOPCBrowseServerAddressSpace struct {
	*IUnknown
}

type IOPCBrowseServerAddressSpaceVtbl struct {
	IUnknownVtbl
	QueryOrganization    uintptr
	ChangeBrowsePosition uintptr
	BrowseOPCItemIDs     uintptr
	GetItemID            uintptr
	BrowseAccessPaths    uintptr
}

func (v *IOPCBrowseServerAddressSpace) Vtbl() *IOPCBrowseServerAddressSpaceVtbl {
	return (*IOPCBrowseServerAddressSpaceVtbl)(unsafe.Pointer(v.IUnknown.LpVtbl))
}

type OPCNAMESPACETYPE uint32

func (v *IOPCBrowseServerAddressSpace) QueryOrganization() (pNameSpaceType OPCNAMESPACETYPE, err error) {
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QueryOrganization,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&pNameSpaceType)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

type OPCBROWSEDIRECTION uint32

func (v *IOPCBrowseServerAddressSpace) ChangeBrowsePosition(dwBrowseDirection OPCBROWSEDIRECTION, szString string) (err error) {
	var pName *uint16
	pName, err = syscall.UTF16PtrFromString(szString)
	if err != nil {
		return
	}

	r0, _, _ := syscall.SyscallN(
		v.Vtbl().ChangeBrowsePosition,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(dwBrowseDirection),
		uintptr(unsafe.Pointer(pName)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

type OPCBROWSETYPE uint32

func (v *IOPCBrowseServerAddressSpace) BrowseOPCItemIDs(dwBrowseFilterType OPCBROWSETYPE, szFilterCriteria string, vtDataTypeFilter uint16, dwAccessRightsFilter uint32) (result []string, err error) {
	var pString *IUnknown
	var pName *uint16
	pName, err = syscall.UTF16PtrFromString(szFilterCriteria)
	if err != nil {
		return
	}

	r0, _, _ := syscall.SyscallN(
		v.Vtbl().BrowseOPCItemIDs,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(dwBrowseFilterType),
		uintptr(unsafe.Pointer(pName)),
		uintptr(vtDataTypeFilter),
		uintptr(dwAccessRightsFilter),
		uintptr(unsafe.Pointer(&pString)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	ppIEnumString := &IEnumString{pString}
	defer func() {
		ppIEnumString.Release()
	}()

	for {
		var batch []string
		batch, err = ppIEnumString.Next(100)
		if err != nil {
			return nil, err
		}
		if len(batch) < 100 {
			result = append(result, batch...)
			break
		}
		result = append(result, batch...)
	}
	return result, nil
}

func (v *IOPCBrowseServerAddressSpace) GetItemID(szItemDataID string) (szItemID string, err error) {
	var pString *uint16
	var pName *uint16
	pName, err = syscall.UTF16PtrFromString(szItemDataID)
	if err != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().GetItemID,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(pName)),
		uintptr(unsafe.Pointer(&pString)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		if pString != nil {
			CoTaskMemFree(unsafe.Pointer(pString))
		}
	}()
	szItemID = windows.UTF16PtrToString(pString)
	return
}
