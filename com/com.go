package com

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modOle32                    = windows.NewLazySystemDLL("ole32.dll")
	procCoCreateInstanceEx      = modOle32.NewProc("CoCreateInstanceEx")
	procCoTaskMemFree           = modOle32.NewProc("CoTaskMemFree")
	modOleaut32                 = windows.NewLazySystemDLL("oleaut32.dll")
	procVariantClear            = modOleaut32.NewProc("VariantClear")
	procVariantTimeToSystemTime = modOleaut32.NewProc("VariantTimeToSystemTime")
	procSystemTimeToVariantTime = modOleaut32.NewProc("SystemTimeToVariantTime")
	procSafeArrayGetVarType     = modOleaut32.NewProc("SafeArrayGetVartype")
	procSafeArrayGetLBound      = modOleaut32.NewProc("SafeArrayGetLBound")
	procSafeArrayGetUBound      = modOleaut32.NewProc("SafeArrayGetUBound")
	procSafeArrayGetElement     = modOleaut32.NewProc("SafeArrayGetElement")
	procSysAllocStringLen       = modOleaut32.NewProc("SysAllocStringLen")
	procSafeArrayCreateVector   = modOleaut32.NewProc("SafeArrayCreateVector")
	procSafeArrayPutElement     = modOleaut32.NewProc("SafeArrayPutElement")
	procSysFreeString           = modOleaut32.NewProc("SysFreeString")
)

func CoTaskMemFree(pv unsafe.Pointer) {
	r0, _, _ := syscall.SyscallN(procCoTaskMemFree.Addr(), uintptr(pv))
	if int32(r0) < 0 {
		panic(syscall.Errno(r0))
	}
}

type CLSCTX uint32

const (
	CLSCTX_LOCAL_SERVER  CLSCTX = 0x4
	CLSCTX_REMOTE_SERVER CLSCTX = 0x10
)

type COAUTHIDENTITY struct {
	User           *uint16
	UserLength     uint32
	Domain         *uint16
	DomainLength   uint32
	Password       *uint16
	PasswordLength uint32
	Flags          uint32
}

type COAUTHINFO struct {
	DwAuthnSvc           uint32
	DwAuthzSvc           uint32
	PwszServerPrincName  *uint16
	DwAuthnLevel         uint32
	DwImpersonationLevel uint32
	PAuthIdentityData    *COAUTHIDENTITY
	DwCapabilities       uint32
}

type COSERVERINFO struct {
	DwReserved1 uint32
	PwszName    *uint16
	PAuthInfo   *COAUTHINFO
	DwReserved2 uint32
}

type MULTI_QI struct {
	PIID *windows.GUID
	PItf *IUnknown
	Hr   int32 // long
}

func CoCreateInstanceEx(Clsid *windows.GUID, punkOuter *IUnknown, dwClsCtx CLSCTX, pServerInfo *COSERVERINFO, dwCount uint32, pResults *MULTI_QI) (ret error) {
	r0, _, _ := syscall.SyscallN(procCoCreateInstanceEx.Addr(), uintptr(unsafe.Pointer(Clsid)), uintptr(unsafe.Pointer(punkOuter)), uintptr(dwClsCtx), uintptr(unsafe.Pointer(pServerInfo)), uintptr(dwCount), uintptr(unsafe.Pointer(pResults)))
	if r0 != 0 {
		ret = syscall.Errno(r0)
	}
	return
}

func VariantClear(pvarg *VARIANT) (err error) {
	r0, _, _ := syscall.SyscallN(procVariantClear.Addr(), uintptr(unsafe.Pointer(pvarg)))
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

func safeArrayGetVarType(safeArray *SafeArray) (varType uint16, err error) {
	r0, _, _ := syscall.SyscallN(procSafeArrayGetVarType.Addr(), uintptr(unsafe.Pointer(safeArray)), uintptr(unsafe.Pointer(&varType)))
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

func safeArrayGetLBound(safeArray *SafeArray, dimension uint32) (lowerBound int32, err error) {
	r0, _, _ := syscall.SyscallN(
		procSafeArrayGetLBound.Addr(),
		uintptr(unsafe.Pointer(safeArray)),
		uintptr(dimension),
		uintptr(unsafe.Pointer(&lowerBound)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

func safeArrayGetUBound(safeArray *SafeArray, dimension uint32) (upperBound int32, err error) {
	r0, _, _ := syscall.SyscallN(
		procSafeArrayGetUBound.Addr(),
		uintptr(unsafe.Pointer(safeArray)),
		uintptr(dimension),
		uintptr(unsafe.Pointer(&upperBound)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

func safeArrayGetElement(safeArray *SafeArray, index int32, pv unsafe.Pointer) (err error) {
	r0, _, _ := syscall.SyscallN(
		procSafeArrayGetElement.Addr(),
		uintptr(unsafe.Pointer(safeArray)),
		uintptr(unsafe.Pointer(&index)),
		uintptr(pv))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

func SysAllocStringLen(v string) (ss *uint16) {
	u := windows.StringToUTF16(v)
	pss, _, _ := procSysAllocStringLen.Call(uintptr(unsafe.Pointer(&u[0])), uintptr(len(u)-1))
	ss = (*uint16)(unsafe.Pointer(pss))
	return
}

func safeArrayCreateVector(variantType VT, lowerBound int32, length uint32) (safearray *SafeArray, err error) {
	r0, _, err := syscall.SyscallN(
		procSafeArrayCreateVector.Addr(),
		uintptr(variantType),
		uintptr(lowerBound),
		uintptr(length),
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	safearray = (*SafeArray)(unsafe.Pointer(r0))
	return safearray, nil
}

func safeArrayPutElement(safearray *SafeArray, index int64, element uintptr) (err error) {
	r0, _, _ := syscall.SyscallN(
		procSafeArrayPutElement.Addr(),
		uintptr(unsafe.Pointer(safearray)),
		uintptr(unsafe.Pointer(&index)),
		element,
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

func SysFreeString(v *uint16) (err error) {
	r0, _, _ := syscall.SyscallN(
		procSysFreeString.Addr(),
		uintptr(unsafe.Pointer(v)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}
