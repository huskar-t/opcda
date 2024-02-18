package com

import (
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IOPCServer = windows.GUID{
	Data1: 0x39c13a4d,
	Data2: 0x011e,
	Data3: 0x11d0,
	Data4: [8]byte{0x96, 0x75, 0x00, 0x20, 0xaf, 0xd8, 0xad, 0xb3},
}

type IOPCServer struct {
	*IUnknown
}

type IOPCServerVtbl struct {
	IUnknownVtbl
	AddGroup              uintptr
	GetErrorString        uintptr
	GetGroupByName        uintptr
	GetStatus             uintptr
	RemoveGroup           uintptr
	CreateGroupEnumerator uintptr
}

func (v *IOPCServer) Vtbl() *IOPCServerVtbl {
	return (*IOPCServerVtbl)(unsafe.Pointer(v.IUnknown.LpVtbl))
}

func (v *IOPCServer) AddGroup(
	szName string,
	bActive bool,
	dwRequestedUpdateRate uint32,
	hClientGroup uint32,
	pTimeBias *int32,
	pPercentDeadband *float32,
	dwLCID uint32,
	riid *windows.GUID,
) (phServerGroup uint32, pRevisedUpdateRate uint32, ppUnk *IUnknown, err error) {
	var pUnk *IUnknown
	var pName *uint16
	pName, err = syscall.UTF16PtrFromString(szName)
	if err != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().AddGroup,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(pName)),
		uintptr(BoolToBOOL(bActive)),
		uintptr(dwRequestedUpdateRate),
		uintptr(hClientGroup),
		uintptr(unsafe.Pointer(pTimeBias)),
		uintptr(unsafe.Pointer(pPercentDeadband)),
		uintptr(dwLCID),
		uintptr(unsafe.Pointer(&phServerGroup)),
		uintptr(unsafe.Pointer(&pRevisedUpdateRate)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&pUnk)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	ppUnk = pUnk
	return
}

func BoolToBOOL(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

type OPCServerState uint32

type OPCSERVERSTATUS struct {
	FtStartTime      windows.Filetime
	FtCurrentTime    windows.Filetime
	FtLastUpdateTime windows.Filetime
	DwServerState    OPCServerState
	DwGroupCount     uint32
	DwBandWidth      uint32
	WMajorVersion    uint16
	WMinorVersion    uint16
	WBuildNumber     uint16
	WReserved        uint16
	SzVendorInfo     *uint16
}

type ServerStatus struct {
	StartTime      time.Time
	CurrentTime    time.Time
	LastUpdateTime time.Time
	ServerState    OPCServerState
	GroupCount     uint32
	BandWidth      uint32
	MajorVersion   uint16
	MinorVersion   uint16
	BuildNumber    uint16
	Reserved       uint16
	VendorInfo     string
}

func (v *IOPCServer) GetStatus() (status *ServerStatus, err error) {
	var pStatus *OPCSERVERSTATUS
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().GetStatus,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&pStatus)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		if pStatus != nil {
			if pStatus.SzVendorInfo != nil {
				CoTaskMemFree(unsafe.Pointer(pStatus.SzVendorInfo))
			}
			CoTaskMemFree(unsafe.Pointer(pStatus))
		}
	}()
	status = &ServerStatus{
		StartTime:      time.Unix(0, pStatus.FtStartTime.Nanoseconds()),
		CurrentTime:    time.Unix(0, pStatus.FtCurrentTime.Nanoseconds()),
		LastUpdateTime: time.Unix(0, pStatus.FtLastUpdateTime.Nanoseconds()),
		ServerState:    pStatus.DwServerState,
		GroupCount:     pStatus.DwGroupCount,
		BandWidth:      pStatus.DwBandWidth,
		MajorVersion:   pStatus.WMajorVersion,
		MinorVersion:   pStatus.WMinorVersion,
		BuildNumber:    pStatus.WBuildNumber,
		Reserved:       pStatus.WReserved,
		VendorInfo:     windows.UTF16PtrToString(pStatus.SzVendorInfo),
	}
	return
}

func (v *IOPCServer) RemoveGroup(hServerGroup uint32, bForce bool) (err error) {
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().RemoveGroup,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(hServerGroup),
		uintptr(BoolToBOOL(bForce)),
		0,
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	return
}
