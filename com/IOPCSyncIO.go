package com

import (
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IOPCSyncIO = windows.GUID{
	Data1: 0x39c13a52,
	Data2: 0x011e,
	Data3: 0x11d0,
	Data4: [8]byte{0x96, 0x75, 0x00, 0x20, 0xaf, 0xd8, 0xad, 0xb3},
}

type IOPCSyncIOVtbl struct {
	IUnknownVtbl
	Read  uintptr
	Write uintptr
}

type IOPCSyncIO struct {
	*IUnknown
}

func (sl *IOPCSyncIO) Vtbl() *IOPCSyncIOVtbl {
	return (*IOPCSyncIOVtbl)(unsafe.Pointer(sl.IUnknown.LpVtbl))
}

type OPCDATASOURCE int32

type TagOPCITEMSTATE struct {
	HClient    uint32
	FTimestamp windows.Filetime
	WQuality   uint16
	WReserved  uint16
	VDataValue VARIANT
}

type ItemState struct {
	Value        interface{}
	Quality      uint16
	Timestamp    time.Time
	ClientHandle int32
}

func (sl *IOPCSyncIO) Read(source OPCDATASOURCE, serverHandles []uint32) ([]*ItemState, []int32, error) {
	var pErrors unsafe.Pointer
	var pValues unsafe.Pointer
	count := len(serverHandles)
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().Read,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(source),
		uintptr(count),
		uintptr(unsafe.Pointer(&serverHandles[0])),
		uintptr(unsafe.Pointer(&pValues)),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pErrors)
		CoTaskMemFree(pValues)
	}()
	errors := make([]int32, count)
	returnValues := make([]*ItemState, count)
	for i := 0; i < count; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		value := *(*TagOPCITEMSTATE)(unsafe.Pointer(uintptr(pValues) + uintptr(i)*unsafe.Sizeof(TagOPCITEMSTATE{})))
		if errNo >= 0 {
			returnValues[i] = &ItemState{
				Value:        value.VDataValue.Value(),
				Quality:      value.WQuality,
				Timestamp:    time.Unix(0, value.FTimestamp.Nanoseconds()),
				ClientHandle: int32(value.HClient),
			}
		}
		value.VDataValue.Clear()
		errors[i] = int32(errNo)
	}
	return returnValues, errors, nil
}

func (sl *IOPCSyncIO) Write(serverHandles []uint32, values []VARIANT) ([]int32, error) {
	var pErrors unsafe.Pointer
	count := len(serverHandles)
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().Write,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(count),
		uintptr(unsafe.Pointer(&serverHandles[0])),
		uintptr(unsafe.Pointer(&values[0])),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pErrors)
	}()
	errors := make([]int32, count)
	for i := 0; i < count; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		errors[i] = int32(errNo)
	}
	return errors, nil
}
