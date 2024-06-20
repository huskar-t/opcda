package com

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var IID_IOPCItemMgt = windows.GUID{
	Data1: 0x39c13a54,
	Data2: 0x011e,
	Data3: 0x11d0,
	Data4: [8]byte{0x96, 0x75, 0x00, 0x20, 0xaf, 0xd8, 0xad, 0xb3},
}

type IOPCItemMgtVtbl struct {
	IUnknownVtbl
	AddItems         uintptr
	ValidateItems    uintptr
	RemoveItems      uintptr
	SetActiveState   uintptr
	SetClientHandles uintptr
	SetDatatypes     uintptr
	CreateEnumerator uintptr
}

type IOPCItemMgt struct {
	*IUnknown
}

func (sl *IOPCItemMgt) Vtbl() *IOPCItemMgtVtbl {
	return (*IOPCItemMgtVtbl)(unsafe.Pointer(sl.IUnknown.LpVtbl))
}

type TagOPCITEMDEF struct {
	SzAccessPath *uint16
	SzItemID     *uint16
	BActive      int32
	HClient      uint32
	DwBlobSize   uint32
	PBlob        *byte
	VtRequested  uint16
	WReserved    uint16
}

type TagOPCITEMRESULT struct {
	HServer        uint32
	VtCanonical    uint16
	WReserved      uint16
	DwAccessRights uint32
	DwBlobSize     uint32
	PBlob          *byte
}

type TagOPCITEMRESULTStruct struct {
	Server       uint32
	NativeType   uint16
	AccessRights uint32
	Blob         []byte
}

func (result *TagOPCITEMRESULT) CloneToStruct() TagOPCITEMRESULTStruct {
	var blob []byte
	if result.DwBlobSize > 0 {
		blob = make([]byte, result.DwBlobSize)
		for i := uint32(0); i < result.DwBlobSize; i++ {
			blob[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(result.PBlob)) + uintptr(i)))
		}
	}
	return TagOPCITEMRESULTStruct{
		Server:       result.HServer,
		NativeType:   result.VtCanonical,
		AccessRights: result.DwAccessRights,
		Blob:         blob,
	}
}

func (sl *IOPCItemMgt) AddItems(items []TagOPCITEMDEF) ([]TagOPCITEMRESULTStruct, []int32, error) {
	dwCount := uint32(len(items))
	var pAddResults unsafe.Pointer
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().AddItems,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwCount),
		uintptr(unsafe.Pointer(&items[0])),
		uintptr(unsafe.Pointer(&pAddResults)),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pAddResults)
		CoTaskMemFree(pErrors)
	}()
	addResults := make([]TagOPCITEMRESULTStruct, dwCount)
	addErrors := make([]int32, dwCount)
	for i := uint32(0); i < dwCount; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		if errNo >= 0 {
			addResults[i] = (*(*TagOPCITEMRESULT)(unsafe.Pointer(uintptr(pAddResults) + uintptr(i)*unsafe.Sizeof(TagOPCITEMRESULT{})))).CloneToStruct()
		}
		addErrors[i] = int32(errNo)
	}
	return addResults, addErrors, nil
}

//        virtual HRESULT STDMETHODCALLTYPE ValidateItems(
//            /* [in] */ DWORD dwCount,
//            /* [size_is][in] */ OPCITEMDEF *pItemArray,
//            /* [in] */ BOOL bBlobUpdate,
//            /* [size_is][size_is][out] */ OPCITEMRESULT **ppValidationResults,
//            /* [size_is][size_is][out] */ HRESULT **ppErrors) = 0;

func (sl *IOPCItemMgt) ValidateItems(items []TagOPCITEMDEF, bBlobUpdate bool) ([]TagOPCITEMRESULTStruct, []int32, error) {
	dwCount := uint32(len(items))
	var pValidationResults unsafe.Pointer
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().ValidateItems,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwCount),
		uintptr(unsafe.Pointer(&items[0])),
		uintptr(BoolToComBOOL(bBlobUpdate)),
		uintptr(unsafe.Pointer(&pValidationResults)),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pValidationResults)
		CoTaskMemFree(pErrors)
	}()
	validationResults := make([]TagOPCITEMRESULTStruct, dwCount)
	validationErrors := make([]int32, dwCount)
	for i := uint32(0); i < dwCount; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		if errNo >= 0 {
			validationResults[i] = (*(*TagOPCITEMRESULT)(unsafe.Pointer(uintptr(pValidationResults) + uintptr(i)*unsafe.Sizeof(TagOPCITEMRESULT{})))).CloneToStruct()
		}
		validationErrors[i] = int32(errNo)
	}
	return validationResults, validationErrors, nil
}

//        virtual HRESULT STDMETHODCALLTYPE RemoveItems(
//            /* [in] */ DWORD dwCount,
//            /* [size_is][in] */ OPCHANDLE *phServer,
//            /* [size_is][size_is][out] */ HRESULT **ppErrors) = 0;

func (sl *IOPCItemMgt) RemoveItems(phServer []uint32) ([]int32, error) {
	dwCount := uint32(len(phServer))
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().RemoveItems,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwCount),
		uintptr(unsafe.Pointer(&phServer[0])),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pErrors)
	}()
	errors := make([]int32, dwCount)
	for i := uint32(0); i < dwCount; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		errors[i] = int32(errNo)
	}
	return errors, nil
}

//        virtual HRESULT STDMETHODCALLTYPE SetActiveState(
//            /* [in] */ DWORD dwCount,
//            /* [size_is][in] */ OPCHANDLE *phServer,
//            /* [in] */ BOOL bActive,
//            /* [size_is][size_is][out] */ HRESULT **ppErrors) = 0;

func (sl *IOPCItemMgt) SetActiveState(phServer []uint32, bActive bool) ([]int32, error) {
	dwCount := uint32(len(phServer))
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().SetActiveState,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwCount),
		uintptr(unsafe.Pointer(&phServer[0])),
		uintptr(BoolToComBOOL(bActive)),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pErrors)
	}()
	errors := make([]int32, dwCount)
	for i := uint32(0); i < dwCount; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		errors[i] = int32(errNo)
	}
	return errors, nil
}

//        virtual HRESULT STDMETHODCALLTYPE SetClientHandles(
//            /* [in] */ DWORD dwCount,
//            /* [size_is][in] */ OPCHANDLE *phServer,
//            /* [size_is][in] */ OPCHANDLE *phClient,
//            /* [size_is][size_is][out] */ HRESULT **ppErrors) = 0;

func (sl *IOPCItemMgt) SetClientHandles(phServer []uint32, phClient []uint32) ([]int32, error) {
	dwCount := uint32(len(phServer))
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().SetClientHandles,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwCount),
		uintptr(unsafe.Pointer(&phServer[0])),
		uintptr(unsafe.Pointer(&phClient[0])),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pErrors)
	}()
	errors := make([]int32, dwCount)
	for i := uint32(0); i < dwCount; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		errors[i] = int32(errNo)
	}
	return errors, nil
}

//        virtual HRESULT STDMETHODCALLTYPE SetDatatypes(
//            /* [in] */ DWORD dwCount,
//            /* [size_is][in] */ OPCHANDLE *phServer,
//            /* [size_is][in] */ VARTYPE *pRequestedDatatypes,
//            /* [size_is][size_is][out] */ HRESULT **ppErrors) = 0;

func (sl *IOPCItemMgt) SetDatatypes(phServer []uint32, pRequestedDatatypes []VT) ([]int32, error) {
	dwCount := uint32(len(phServer))
	var pErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		sl.Vtbl().SetDatatypes,
		uintptr(unsafe.Pointer(sl.IUnknown)),
		uintptr(dwCount),
		uintptr(unsafe.Pointer(&phServer[0])),
		uintptr(unsafe.Pointer(&pRequestedDatatypes[0])),
		uintptr(unsafe.Pointer(&pErrors)),
	)
	if int32(r0) < 0 {
		return nil, syscall.Errno(r0)
	}
	defer func() {
		CoTaskMemFree(pErrors)
	}()
	errors := make([]int32, dwCount)
	for i := uint32(0); i < dwCount; i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*4))
		errors[i] = int32(errNo)
	}

	return errors, nil
}
