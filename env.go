package opcda

import (
	"golang.org/x/sys/windows"
)

// Initialize initialize COM with COINIT_MULTITHREADED
func Initialize() {
	windows.CoInitializeEx(0, windows.COINIT_MULTITHREADED)
}

// Uninitialize uninitialize COM
func Uninitialize() {
	windows.CoUninitialize()
}
