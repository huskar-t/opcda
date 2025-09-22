package com

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modOle32                    = windows.NewLazySystemDLL("ole32.dll")
	procCoCreateInstanceEx      = modOle32.NewProc("CoCreateInstanceEx")
	procCoInitializeSecurity    = modOle32.NewProc("CoInitializeSecurity")
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
	if pv == nil {
		return
	}
	windows.CoTaskMemFree(pv)
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

func MakeCOMObjectEx(hostname string, serverLocation CLSCTX, requestedClass *windows.GUID, requestedInterface *windows.GUID) (*IUnknown, error) {
	reqInterface := MULTI_QI{
		PIID: requestedInterface,
		PItf: nil,
		Hr:   0,
	}
	var serverInfoPtr *COSERVERINFO = nil
	if serverLocation != CLSCTX_LOCAL_SERVER {
		serverInfoPtr = &COSERVERINFO{
			PwszName: windows.StringToUTF16Ptr(hostname),
		}
	}
	err := CoCreateInstanceEx(requestedClass, nil, serverLocation, serverInfoPtr, 1, &reqInterface)
	if err != nil {
		return nil, err
	}
	if reqInterface.Hr != 0 {
		return nil, syscall.Errno(reqInterface.Hr)
	}
	return reqInterface.PItf, nil
}

func IsLocal(host string) bool {
	if host == "" || host == "localhost" || host == "127.0.0.1" {
		return true
	}
	name, err := windows.ComputerName()
	if err != nil {
		return false
	}
	return strings.ToLower(name) == strings.ToLower(host)
}

// Initialize initialize COM with COINIT_MULTITHREADED
func Initialize() error {
	err := windows.CoInitializeEx(0, windows.COINIT_MULTITHREADED)
	if err != nil {
		return fmt.Errorf("call CoInitializeEx error: %s", err)
	}
	err = CoInitializeSecurity(RPC_C_AUTHN_LEVEL_NONE, RPC_C_IMP_LEVEL_IMPERSONATE, EOAC_NONE)
	if err != nil {
		Uninitialize()
		return fmt.Errorf("call CoInitializeSecurity error: %s", err)
	}
	return nil
}

// Uninitialize uninitialize COM
func Uninitialize() {
	windows.CoUninitialize()
}

func IsEqualGUID(guid1 *windows.GUID, guid2 *windows.GUID) bool {
	return guid1.Data1 == guid2.Data1 &&
		guid1.Data2 == guid2.Data2 &&
		guid1.Data3 == guid2.Data3 &&
		guid1.Data4[0] == guid2.Data4[0] &&
		guid1.Data4[1] == guid2.Data4[1] &&
		guid1.Data4[2] == guid2.Data4[2] &&
		guid1.Data4[3] == guid2.Data4[3] &&
		guid1.Data4[4] == guid2.Data4[4] &&
		guid1.Data4[5] == guid2.Data4[5] &&
		guid1.Data4[6] == guid2.Data4[6] &&
		guid1.Data4[7] == guid2.Data4[7]
}

func CoInitializeSecurity(authnLevel, impLevel, capabilities uint32) (err error) {
	cAuthSvc := int32(-1)
	r0, _, _ := procCoInitializeSecurity.Call(
		uintptr(0),
		uintptr(cAuthSvc),
		uintptr(0),
		uintptr(0),
		uintptr(authnLevel),
		uintptr(impLevel),
		uintptr(0),
		uintptr(capabilities),
		uintptr(0))
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}
