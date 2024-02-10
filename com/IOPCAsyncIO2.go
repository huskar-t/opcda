package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IOPCAsyncIO2 = windows.GUID{
	Data1: 0x39c13a71,
	Data2: 0x011e,
	Data3: 0x11d0,
	Data4: [8]byte{0x96, 0x75, 0x00, 0x20, 0xaf, 0xd8, 0xad, 0xb3},
}

type IOPCAsyncIO2Vtbl struct {
	IUnknownVtbl
	Read      uintptr
	Write     uintptr
	Refresh2  uintptr
	Cancel2   uintptr
	SetEnable uintptr
	GetEnable uintptr
}

type IOPCAsyncIO2 struct {
	*IUnknown
}

func (sl *IOPCAsyncIO2) Vtbl() *IOPCAsyncIO2Vtbl {
	return (*IOPCAsyncIO2Vtbl)(unsafe.Pointer(sl.IUnknown.LpVtbl))
}

//        virtual HRESULT STDMETHODCALLTYPE Read(
//            /* [in] */ DWORD dwCount,
//            /* [size_is][in] */ OPCHANDLE * phServer,
//            /* [in] */ DWORD dwTransactionID,
//            /* [out] */ DWORD * pdwCancelID,
//            /* [size_is][size_is][out] */ HRESULT * *ppErrors) = 0;

func (sl *IOPCAsyncIO2) Read(phServer []uint32, dwTransactionID uint32) (pdwCancelID uint32, ppErrors []int32, err error) {
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().Read,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(len(phServer)),
		uintptr(unsafe.Pointer(&phServer[0])),
		uintptr(dwTransactionID),
		uintptr(unsafe.Pointer(&pdwCancelID)),
		uintptr(unsafe.Pointer(&pErrors)))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		if pErrors != nil {
			CoTaskMemFree(pErrors)
		}
	}()
	ppErrors = make([]int32, len(phServer))
	for i := uint32(0); i < uint32(len(phServer)); i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		ppErrors[i] = int32(errNo)
	}
	return
}

//        virtual HRESULT STDMETHODCALLTYPE Write(
//            /* [in] */ DWORD dwCount,
//            /* [size_is][in] */ OPCHANDLE * phServer,
//            /* [size_is][in] */ VARIANT * pItemValues,
//            /* [in] */ DWORD dwTransactionID,
//            /* [out] */ DWORD * pdwCancelID,
//            /* [size_is][size_is][out] */ HRESULT * *ppErrors) = 0;

func (sl *IOPCAsyncIO2) Write(phServer []uint32, pItemValues []VARIANT, dwTransactionID uint32) (pdwCancelID uint32, ppErrors []int32, err error) {
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().Write,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(len(phServer)),
		uintptr(unsafe.Pointer(&phServer[0])),
		uintptr(unsafe.Pointer(&pItemValues[0])),
		uintptr(dwTransactionID),
		uintptr(unsafe.Pointer(&pdwCancelID)),
		uintptr(unsafe.Pointer(&pErrors)))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		if pErrors != nil {
			CoTaskMemFree(pErrors)
		}
	}()
	ppErrors = make([]int32, len(phServer))
	for i := uint32(0); i < uint32(len(phServer)); i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		ppErrors[i] = int32(errNo)
	}
	return
}

//        virtual HRESULT STDMETHODCALLTYPE Refresh2(
//            /* [in] */ OPCDATASOURCE dwSource,
//            /* [in] */ DWORD dwTransactionID,
//            /* [out] */ DWORD * pdwCancelID) = 0;

func (sl *IOPCAsyncIO2) Refresh2(dwSource OPCDATASOURCE, dwTransactionID uint32) (pdwCancelID uint32, err error) {
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().Refresh2,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwSource),
		uintptr(dwTransactionID),
		uintptr(unsafe.Pointer(&pdwCancelID)))
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	return
}

//        virtual HRESULT STDMETHODCALLTYPE Cancel2(
//            /* [in] */ DWORD dwCancelID) = 0;

func (sl *IOPCAsyncIO2) Cancel2(dwCancelID uint32) (err error) {
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().Cancel2,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwCancelID),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	return
}
