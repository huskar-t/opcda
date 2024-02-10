package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IUnknown = &windows.GUID{
	Data1: 0x00000000,
	Data2: 0x0000,
	Data3: 0x0000,
	Data4: [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
}

type IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

type IUnknown struct {
	LpVtbl *[1024]uintptr
}

func (v *IUnknown) Vtbl() *IUnknownVtbl {
	return (*IUnknownVtbl)(unsafe.Pointer(v.LpVtbl))
}

// QueryInterface
// HRESULT QueryInterface(
//
//	REFIID riid,
//	void   **ppvObject
//
// );
func (v *IUnknown) QueryInterface(riid *windows.GUID, ppvObject unsafe.Pointer) (ret error) {
	r0, _, _ := syscall.SyscallN(v.Vtbl().QueryInterface, uintptr(unsafe.Pointer(v)), uintptr(unsafe.Pointer(riid)), uintptr(ppvObject))
	if r0 != 0 {
		ret = syscall.Errno(r0)
	}
	return
}

// Release
// ULONG Release();
func (v *IUnknown) Release() uint32 {
	ret, _, _ := syscall.SyscallN(v.Vtbl().Release, uintptr(unsafe.Pointer(v)))
	return uint32(ret)
}
