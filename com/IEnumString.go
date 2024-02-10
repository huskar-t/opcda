package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type IEnumStringVtbl struct {
	IUnknownVtbl
	Next  uintptr
	Skip  uintptr
	Reset uintptr
	Clone uintptr
}

type IEnumString struct {
	*IUnknown
}

func (sl *IEnumString) Vtbl() *IEnumStringVtbl {
	return (*IEnumStringVtbl)(unsafe.Pointer(sl.IUnknown.LpVtbl))
}

//        virtual /* [local] */ HRESULT STDMETHODCALLTYPE Next(
//            /* [in] */ ULONG celt,
//            /* [annotation] */
//            _Out_writes_to_(celt,*pceltFetched)  LPOLESTR *rgelt,
//            /* [annotation] */
//            _Out_opt_  ULONG *pceltFetched) = 0;

func (sl *IEnumString) Next(celt uint32) (result []string, err error) {
	pRgelt := make([]*uint16, celt)
	var pceltFetched uint32
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().Next,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(celt),
		uintptr(unsafe.Pointer(&pRgelt[0])),
		uintptr(unsafe.Pointer(&pceltFetched)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	result = make([]string, pceltFetched)
	for i := uint32(0); i < pceltFetched; i++ {
		pwstr := pRgelt[i]
		result[i] = windows.UTF16PtrToString(pwstr)
		CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	return
}
