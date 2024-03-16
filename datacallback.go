package opcda

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/huskar-t/opcda/com"

	"golang.org/x/sys/windows"
)

var IID_IOPCDataCallback = windows.GUID{
	Data1: 0x39c13a70,
	Data2: 0x011e,
	Data3: 0x11d0,
	Data4: [8]byte{0x96, 0x75, 0x00, 0x20, 0xaf, 0xd8, 0xad, 0xb3},
}

type DataEventReceiver struct {
	lpVtbl                 *DataEventReceiverVtbl
	ref                    int32
	clsid                  *windows.GUID
	dataChangeReceiver     chan *CDataChangeCallBackData
	readCompleteReceiver   chan *CReadCompleteCallBackData
	writeCompleteReceiver  chan *CWriteCompleteCallBackData
	cancelCompleteReceiver chan *CCancelCompleteCallBackData
}

type DataEventReceiverVtbl struct {
	pQueryInterface   uintptr
	pAddRef           uintptr
	pRelease          uintptr
	pOnDataChange     uintptr
	pOnReadComplete   uintptr
	pOnWriteComplete  uintptr
	pOnCancelComplete uintptr
}

func NewDataEventReceiver(
	dataChangeReceiver chan *CDataChangeCallBackData,
	readCompleteReceiver chan *CReadCompleteCallBackData,
	writeCompleteReceiver chan *CWriteCompleteCallBackData,
	cancelCompleteReceiver chan *CCancelCompleteCallBackData,
) *DataEventReceiver {
	return &DataEventReceiver{
		lpVtbl: &DataEventReceiverVtbl{
			pQueryInterface:   syscall.NewCallback(DataQueryInterface),
			pAddRef:           syscall.NewCallback(DataAddRef),
			pRelease:          syscall.NewCallback(DataRelease),
			pOnDataChange:     syscall.NewCallback(DataOnDataChange),
			pOnReadComplete:   syscall.NewCallback(DataOnReadComplete),
			pOnWriteComplete:  syscall.NewCallback(DataOnWriteComplete),
			pOnCancelComplete: syscall.NewCallback(DataOnCancelComplete),
		},
		ref:                    0,
		clsid:                  &IID_IOPCDataCallback,
		dataChangeReceiver:     dataChangeReceiver,
		readCompleteReceiver:   readCompleteReceiver,
		writeCompleteReceiver:  writeCompleteReceiver,
		cancelCompleteReceiver: cancelCompleteReceiver,
	}
}

func DataQueryInterface(this unsafe.Pointer, iid *windows.GUID, punk *unsafe.Pointer) uintptr {
	er := (*DataEventReceiver)(this)
	*punk = nil
	if com.IsEqualGUID(iid, er.clsid) || com.IsEqualGUID(iid, com.IID_IUnknown) {
		DataAddRef(this)
		*punk = this
		return com.S_OK
	}
	return com.E_POINTER
}

func DataAddRef(this unsafe.Pointer) uintptr {
	er := (*DataEventReceiver)(this)
	er.ref++
	return uintptr(er.ref)
}

func DataRelease(this unsafe.Pointer) uintptr {
	er := (*DataEventReceiver)(this)
	er.ref--
	return uintptr(er.ref)
}

type CDataChangeCallBackData struct {
	TransID           uint32
	GroupHandle       uint32
	MasterQuality     int32
	MasterErr         int32
	ItemClientHandles []uint32
	Values            []interface{}
	Qualities         []uint16
	TimeStamps        []time.Time
	Errors            []int32
}

func DataOnDataChange(this unsafe.Pointer, dwTransid uint32, hGroup uint32, hrMasterquality int32, hrMastererror int32, dwCount uint32, phClientItems unsafe.Pointer, pvValues unsafe.Pointer, pwQualities unsafe.Pointer, pftTimeStamps unsafe.Pointer, pErrors unsafe.Pointer) uintptr {
	er := (*DataEventReceiver)(this)
	clientHandles := make([]uint32, dwCount)
	values := make([]interface{}, dwCount)
	qualities := make([]uint16, dwCount)
	timestamps := make([]time.Time, dwCount)
	errors := make([]int32, dwCount)
	for i := 0; i < int(dwCount); i++ {
		clientHandles[i] = *(*uint32)(unsafe.Pointer(uintptr(phClientItems) + uintptr(i)*unsafe.Sizeof(uint32(0))))
		variant := *(*com.VARIANT)(unsafe.Pointer(uintptr(pvValues) + uintptr(i)*unsafe.Sizeof(com.VARIANT{})))
		values[i] = variant.Value()
		qualities[i] = *(*uint16)(unsafe.Pointer(uintptr(pwQualities) + uintptr(i)*unsafe.Sizeof(uint16(0))))
		ft := *(*windows.Filetime)(unsafe.Pointer(uintptr(pftTimeStamps) + uintptr(i)*unsafe.Sizeof(windows.Filetime{})))
		timestamps[i] = time.Unix(0, ft.Nanoseconds())
		errors[i] = *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*unsafe.Sizeof(int32(0))))
	}
	cb := &CDataChangeCallBackData{
		TransID:           dwTransid,
		GroupHandle:       hGroup,
		MasterQuality:     hrMasterquality,
		MasterErr:         hrMastererror,
		ItemClientHandles: clientHandles,
		Values:            values,
		Qualities:         qualities,
		TimeStamps:        timestamps,
		Errors:            errors,
	}
	er.dataChangeReceiver <- cb
	return com.S_OK
}

type CReadCompleteCallBackData struct {
	TransID           uint32
	GroupHandle       uint32
	MasterQuality     int32
	MasterErr         int32
	ItemClientHandles []uint32
	Values            []interface{}
	Qualities         []uint16
	TimeStamps        []time.Time
	Errors            []int32
}

func DataOnReadComplete(this unsafe.Pointer, dwTransid uint32, hGroup uint32, hrMasterquality int32, hrMastererror int32, dwCount uint32, phClientItems unsafe.Pointer, pvValues unsafe.Pointer, pwQualities unsafe.Pointer, pftTimeStamps unsafe.Pointer, pErrors unsafe.Pointer) uintptr {
	er := (*DataEventReceiver)(this)
	clientHandles := make([]uint32, dwCount)
	values := make([]interface{}, dwCount)
	qualities := make([]uint16, dwCount)
	timestamps := make([]time.Time, dwCount)
	errors := make([]int32, dwCount)
	for i := 0; i < int(dwCount); i++ {
		clientHandles[i] = *(*uint32)(unsafe.Pointer(uintptr(phClientItems) + uintptr(i)*unsafe.Sizeof(uint32(0))))
		variant := *(*com.VARIANT)(unsafe.Pointer(uintptr(pvValues) + uintptr(i)*unsafe.Sizeof(com.VARIANT{})))
		values[i] = variant.Value()
		qualities[i] = *(*uint16)(unsafe.Pointer(uintptr(pwQualities) + uintptr(i)*unsafe.Sizeof(uint16(0))))
		ft := *(*windows.Filetime)(unsafe.Pointer(uintptr(pftTimeStamps) + uintptr(i)*unsafe.Sizeof(windows.Filetime{})))
		timestamps[i] = time.Unix(0, ft.Nanoseconds())
		errors[i] = *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*unsafe.Sizeof(int32(0))))
	}
	cb := &CReadCompleteCallBackData{
		TransID:           dwTransid,
		GroupHandle:       hGroup,
		MasterQuality:     hrMasterquality,
		MasterErr:         hrMastererror,
		ItemClientHandles: clientHandles,
		Values:            values,
		Qualities:         qualities,
		TimeStamps:        timestamps,
		Errors:            errors,
	}
	er.readCompleteReceiver <- cb
	return com.S_OK
}

type CWriteCompleteCallBackData struct {
	TransID           uint32
	GroupHandle       uint32
	MasterErr         int32
	ItemClientHandles []uint32
	Errors            []int32
}

func DataOnWriteComplete(this unsafe.Pointer, dwTransid uint32, hGroup uint32, hrMastererr int32, dwCount uint32, pClienthandles unsafe.Pointer, pErrors unsafe.Pointer) uintptr {
	er := (*DataEventReceiver)(this)
	clientHandles := make([]uint32, dwCount)
	errors := make([]int32, dwCount)
	for i := 0; i < int(dwCount); i++ {
		clientHandles[i] = *(*uint32)(unsafe.Pointer(uintptr(pClienthandles) + uintptr(i)*unsafe.Sizeof(uint32(0))))
		errors[i] = *(*int32)(unsafe.Pointer(uintptr(pErrors) + uintptr(i)*unsafe.Sizeof(int32(0))))
	}
	cb := &CWriteCompleteCallBackData{
		TransID:           dwTransid,
		GroupHandle:       hGroup,
		MasterErr:         hrMastererr,
		ItemClientHandles: clientHandles,
		Errors:            errors,
	}
	er.writeCompleteReceiver <- cb
	return com.S_OK
}

type CCancelCompleteCallBackData struct {
	TransID     uint32
	GroupHandle uint32
}

func DataOnCancelComplete(this unsafe.Pointer, dwTransid uint32, hGroup uint32) uintptr {
	er := (*DataEventReceiver)(this)
	cb := &CCancelCompleteCallBackData{
		TransID:     dwTransid,
		GroupHandle: hGroup,
	}
	er.cancelCompleteReceiver <- cb
	return com.S_OK
}
