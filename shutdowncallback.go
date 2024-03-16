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
	if com.IsEqualGUID(iid, er.clsid) || com.IsEqualGUID(iid, com.IID_IUnknown) {
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
