package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type IEnumGUIDVtbl struct {
	IUnknownVtbl
	Next  uintptr
	Skip  uintptr
	Reset uintptr
	Clone uintptr
}

type IEnumGUID struct {
	*IUnknown
}

func (ie *IEnumGUID) Vtbl() *IEnumGUIDVtbl {
	return (*IEnumGUIDVtbl)(unsafe.Pointer(ie.IUnknown.LpVtbl))
}

func (ie *IEnumGUID) Next(celt uint32, rgelt *windows.GUID, pceltFetched *uint32) error {
	r0, _, _ := syscall.SyscallN(ie.Vtbl().Next, uintptr(unsafe.Pointer(ie.IUnknown)), uintptr(celt), uintptr(unsafe.Pointer(rgelt)), uintptr(unsafe.Pointer(pceltFetched)))
	if r0 != 0 {
		return syscall.Errno(r0)
	}
	return nil
}
