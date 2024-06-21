package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type IConnectionPointVtbl struct {
	IUnknownVtbl
	GetConnectionInterface      uintptr
	GetConnectionPointContainer uintptr
	Advise                      uintptr
	Unadvise                    uintptr
	EnumConnections             uintptr
}

type IConnectionPoint struct {
	*IUnknown
}

func (p *IConnectionPoint) Vtbl() *IConnectionPointVtbl {
	return (*IConnectionPointVtbl)(unsafe.Pointer(p.IUnknown.LpVtbl))
}

func (p *IConnectionPoint) Advise(pUnkSink *IUnknown) (cookie uint32, err error) {
	r0, _, _ := syscall.SyscallN(p.Vtbl().Advise, uintptr(unsafe.Pointer(p.IUnknown)), uintptr(unsafe.Pointer(pUnkSink)), uintptr(unsafe.Pointer(&cookie)))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

func (p *IConnectionPoint) Unadvise(dwCookie uint32) error {
	r0, _, _ := syscall.SyscallN(p.Vtbl().Unadvise, uintptr(unsafe.Pointer(p.IUnknown)), uintptr(dwCookie))
	if int32(r0) < 0 {
		return syscall.Errno(r0)
	}
	return nil
}

// B196B284-BAB4-101A-B69C-00AA00341D07
var IID_IConnectionPointContainer = windows.GUID{
	Data1: 0xB196B284,
	Data2: 0xBAB4,
	Data3: 0x101A,
	Data4: [8]byte{0xB6, 0x9C, 0x00, 0xAA, 0x00, 0x34, 0x1D, 0x07},
}

type IConnectionPointContainerVtbl struct {
	IUnknownVtbl
	EnumConnectionPoints uintptr
	FindConnectionPoint  uintptr
}

type IConnectionPointContainer struct {
	*IUnknown
}

func (c *IConnectionPointContainer) Vtbl() *IConnectionPointContainerVtbl {
	return (*IConnectionPointContainerVtbl)(unsafe.Pointer(c.IUnknown.LpVtbl))
}

func (c *IConnectionPointContainer) FindConnectionPoint(riid *windows.GUID) (*IConnectionPoint, error) {
	var iUnknown *IUnknown
	r0, _, _ := syscall.SyscallN(
		c.Vtbl().FindConnectionPoint,
		uintptr(unsafe.Pointer(c.IUnknown)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&iUnknown)),
	)
	if int32(r0) < 0 {
		return nil, syscall.Errno(r0)
	}
	return &IConnectionPoint{iUnknown}, nil
}
