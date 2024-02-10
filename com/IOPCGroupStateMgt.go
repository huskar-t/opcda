package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IOPCGroupStateMgt = windows.GUID{
	Data1: 0x39c13a50,
	Data2: 0x011e,
	Data3: 0x11d0,
	Data4: [8]byte{0x96, 0x75, 0x00, 0x20, 0xaf, 0xd8, 0xad, 0xb3},
}

type IOPCGroupStateMgtVtbl struct {
	IUnknownVtbl
	GetState   uintptr
	SetState   uintptr
	SetName    uintptr
	CloneGroup uintptr
}

type IOPCGroupStateMgt struct {
	*IUnknown
}

func (sl *IOPCGroupStateMgt) Vtbl() *IOPCGroupStateMgtVtbl {
	return (*IOPCGroupStateMgtVtbl)(unsafe.Pointer(sl.IUnknown.LpVtbl))
}

func (sl *IOPCGroupStateMgt) GetState() (pUpdateRate uint32, pActive bool, ppName string, pTimeBias int32, pPercentDeadband float32, pLCID uint32, phClientGroup uint32, phServerGroup uint32, err error) {
	var pName *uint16
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().GetState,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(unsafe.Pointer(&pUpdateRate)),
		uintptr(unsafe.Pointer(&pActive)),
		uintptr(unsafe.Pointer(&pName)),
		uintptr(unsafe.Pointer(&pTimeBias)),
		uintptr(unsafe.Pointer(&pPercentDeadband)),
		uintptr(unsafe.Pointer(&pLCID)),
		uintptr(unsafe.Pointer(&phClientGroup)),
		uintptr(unsafe.Pointer(&phServerGroup)))
	if r0 != 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		if pName != nil {
			CoTaskMemFree(unsafe.Pointer(pName))
		}
	}()
	ppName = windows.UTF16PtrToString(pName)
	return
}

func (sl *IOPCGroupStateMgt) SetState(requestedUpdateRate *uint32, pActive *int32, pTimeBias *int32, pPercentDeadband *float32, pLCID *uint32, phClientGroup *uint32) (pRevisedUpdateRate uint32, err error) {
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().SetState,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(unsafe.Pointer(requestedUpdateRate)),
		uintptr(unsafe.Pointer(&pRevisedUpdateRate)),
		uintptr(unsafe.Pointer(pActive)),
		uintptr(unsafe.Pointer(pTimeBias)),
		uintptr(unsafe.Pointer(pPercentDeadband)),
		uintptr(unsafe.Pointer(pLCID)),
		uintptr(unsafe.Pointer(phClientGroup)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
		return
	}
	return
}

func (sl *IOPCGroupStateMgt) SetName(szName string) (err error) {
	var pName *uint16
	pName, err = syscall.UTF16PtrFromString(szName)
	if err != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().SetName,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(unsafe.Pointer(pName)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
		return
	}
	return
}
