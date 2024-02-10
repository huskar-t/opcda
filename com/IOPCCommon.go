package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IOPCCommon = windows.GUID{
	Data1: 0xF31DFDE2,
	Data2: 0x07B6,
	Data3: 0x11d2,
	Data4: [8]byte{0xB2, 0xD8, 0x00, 0x60, 0x08, 0x3B, 0xA1, 0xFB},
}

type IOPCCommon struct {
	*IUnknown
}

type IOPCCommonVtbl struct {
	IUnknownVtbl
	SetLocaleID             uintptr
	GetLocaleID             uintptr
	QueryAvailableLocaleIDs uintptr
	GetErrorString          uintptr
	SetClientName           uintptr
}

func (v *IOPCCommon) Vtbl() *IOPCCommonVtbl {
	return (*IOPCCommonVtbl)(unsafe.Pointer(v.IUnknown.LpVtbl))
}

func (v *IOPCCommon) SetLocaleID(dwLcid uint32) (err error) {
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().SetLocaleID,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(dwLcid))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

func (v *IOPCCommon) GetLocaleID() (pdwLcid uint32, err error) {
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().GetLocaleID,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&pdwLcid)))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

func (v *IOPCCommon) QueryAvailableLocaleIDs() (result []uint32, err error) {
	var pLcid unsafe.Pointer
	var pdwCount uint32
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QueryAvailableLocaleIDs,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&pdwCount)),
		uintptr(unsafe.Pointer(&pLcid)))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		CoTaskMemFree(pLcid)
	}()
	if pdwCount == 0 {
		return
	}
	result = make([]uint32, pdwCount)
	for i := uint32(0); i < pdwCount; i++ {
		result[i] = *(*uint32)(unsafe.Pointer(uintptr(pLcid) + uintptr(i)*4))
	}
	return
}

func (v *IOPCCommon) GetErrorString(dwError uint32) (str string, err error) {
	var pString *uint16
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().GetErrorString,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(dwError),
		uintptr(unsafe.Pointer(&pString)))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	defer func() {
		if pString != nil {
			CoTaskMemFree(unsafe.Pointer(pString))
		}
	}()
	str = windows.UTF16PtrToString(pString)
	return
}

func (v *IOPCCommon) SetClientName(szName string) (err error) {
	var pName *uint16
	pName, err = syscall.UTF16PtrFromString(szName)
	if err != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().SetClientName,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(pName)))
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}
