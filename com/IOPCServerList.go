package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var CLSID_OpcServerList = windows.GUID{
	Data1: 0x13486D51,
	Data2: 0x4821,
	Data3: 0x11D2,
	Data4: [8]byte{0xA4, 0x94, 0x3C, 0xB3, 0x06, 0xC1, 0x00, 0x00},
}

var IID_IOPCServerList2 = windows.GUID{
	Data1: 0x9DD0B56C,
	Data2: 0xAD9E,
	Data3: 0x43ee,
	Data4: [8]byte{0x83, 0x05, 0x48, 0x7F, 0x31, 0x88, 0xBF, 0x7A},
}

type IOPCServerListVtbl struct {
	IUnknownVtbl
	EnumClassesOfCategories uintptr
	GetClassDetails         uintptr
	CLSIDFromProgID         uintptr
}

type IOPCServerList2 struct {
	*IUnknown
}

func (sl *IOPCServerList2) Vtbl() *IOPCServerListVtbl {
	return (*IOPCServerListVtbl)(unsafe.Pointer(sl.IUnknown.LpVtbl))
}

func (sl *IOPCServerList2) EnumClassesOfCateGories(rgcatidImpl []windows.GUID, rgcatidReq []windows.GUID) (ppenumClsid *IEnumGUID, err error) {
	var r0 uintptr
	cImplemented := uint32(len(rgcatidImpl))
	cRequired := uint32(len(rgcatidReq))
	var iUnknown *IUnknown
	if cRequired == 0 {
		r0, _, _ = syscall.SyscallN(sl.Vtbl().EnumClassesOfCategories, uintptr(unsafe.Pointer(sl.IUnknown)), uintptr(cImplemented), uintptr(unsafe.Pointer(&rgcatidImpl[0])), uintptr(0), uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(&iUnknown)))
	} else {
		r0, _, _ = syscall.SyscallN(sl.Vtbl().EnumClassesOfCategories, uintptr(unsafe.Pointer(sl.IUnknown)), uintptr(cImplemented), uintptr(unsafe.Pointer(&rgcatidImpl[0])), uintptr(cRequired), uintptr(unsafe.Pointer(&rgcatidReq[0])), uintptr(unsafe.Pointer(&iUnknown)))
	}
	if r0 != 0 {
		err = syscall.Errno(r0)
		return
	}
	ppenumClsid = &IEnumGUID{IUnknown: iUnknown}
	return
}

func (sl *IOPCServerList2) GetClassDetails(guid *windows.GUID) (*uint16, *uint16, *uint16, error) {
	var ppszProgID, ppszUserType, ppszVerIndProgID *uint16
	r0, _, _ := syscall.SyscallN(sl.Vtbl().GetClassDetails, uintptr(unsafe.Pointer(sl.IUnknown)), uintptr(unsafe.Pointer(guid)), uintptr(unsafe.Pointer(&ppszProgID)), uintptr(unsafe.Pointer(&ppszUserType)), uintptr(unsafe.Pointer(&ppszVerIndProgID)))
	if r0 != 0 {
		return nil, nil, nil, syscall.Errno(r0)
	}
	return ppszProgID, ppszUserType, ppszVerIndProgID, nil
}
