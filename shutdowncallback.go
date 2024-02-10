package opcda

import (
	"syscall"
	"unsafe"

	"github.com/huskar-t/opcda/com"

	"golang.org/x/sys/windows"
)

type ShutdownEventReceiver struct {
	lpVtbl   *ShutdownEventReceiverVtbl
	ref      int32
	clsid    *windows.GUID
	receiver []chan string
}

type ShutdownEventReceiverVtbl struct {
	pQueryInterface  uintptr
	pAddRef          uintptr
	pRelease         uintptr
	pShutdownRequest uintptr
}

func NewShutdownEventReceiver() *ShutdownEventReceiver {
	return &ShutdownEventReceiver{
		lpVtbl: &ShutdownEventReceiverVtbl{
			pQueryInterface:  syscall.NewCallback(ShutdownQueryInterface),
			pAddRef:          syscall.NewCallback(ShutdownAddRef),
			pRelease:         syscall.NewCallback(ShutdownRelease),
			pShutdownRequest: syscall.NewCallback(ShutdownRequest),
		},
		ref:   0,
		clsid: &IID_IOPCShutdown,
	}
}

func (er *ShutdownEventReceiver) AddReceiver(ch chan string) {
	er.receiver = append(er.receiver, ch)
}

func ShutdownQueryInterface(this unsafe.Pointer, iid *windows.GUID, punk *unsafe.Pointer) uintptr {
	er := (*ShutdownEventReceiver)(this)
	*punk = nil
	if IsEqualGUID(iid, er.clsid) || IsEqualGUID(iid, com.IID_IUnknown) {
		ShutdownAddRef(this)
		*punk = this
		return com.S_OK
	}
	return com.E_POINTER
}

func ShutdownRequest(this *com.IUnknown, pReason *uint16) uintptr {
	er := (*ShutdownEventReceiver)(unsafe.Pointer(this))
	reason := windows.UTF16PtrToString(pReason)
	for _, ch := range er.receiver {
		select {
		case ch <- reason:
		default:
		}
	}
	return uintptr(com.S_OK)
}

func ShutdownAddRef(this unsafe.Pointer) uintptr {
	er := (*ShutdownEventReceiver)(this)
	er.ref++
	return uintptr(er.ref)
}

func ShutdownRelease(this unsafe.Pointer) uintptr {
	er := (*ShutdownEventReceiver)(this)
	er.ref--
	return uintptr(er.ref)
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
