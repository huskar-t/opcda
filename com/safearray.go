package com

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type SafeArray struct {
	Dimensions   uint16
	FeaturesFlag uint16
	ElementsSize uint32
	LocksAmount  uint32
	Data         uint32
	Bounds       [16]byte
}

type SafeArrayBound struct {
	Elements   uint32
	LowerBound int32
}

func (s *SafeArray) ToValueArray() (values []interface{}, err error) {
	totalElements, _ := s.TotalElements(0)
	values = make([]interface{}, totalElements)
	vt, _ := safeArrayGetVarType(s)

	for i := int32(0); i < totalElements; i++ {
		switch VT(vt) {
		case VT_BOOL:
			var v bool
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_I1:
			var v int8
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_I2:
			var v int16
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_I4:
			var v int32
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_I8:
			var v int64
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_UI1:
			var v uint8
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_UI2:
			var v uint16
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_UI4:
			var v uint32
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_UI8:
			var v uint64
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_R4:
			var v float32
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_R8:
			var v float64
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		case VT_BSTR:
			var element *uint16
			err = safeArrayGetElement(s, i, unsafe.Pointer(&element))
			if err != nil {
				return nil, err
			}
			values[i] = windows.UTF16PtrToString(element)
			SysFreeString(element)
		case VT_DATE:
			var v uint64
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			date, err := GetVariantDate(v)
			if err != nil {
				return nil, err
			}
			values[i] = date
		//case VT_CY:
		//	var v int64
		//	err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
		//	if err != nil {
		//		return nil, err
		//	}
		//	values[i] = float64(v) / 10000
		default:
			return nil, fmt.Errorf("unknown value type %x", VT(vt))
		}
	}

	return
}

func (s *SafeArray) TotalElements(index uint32) (totalElements int32, err error) {
	if index < 1 {
		index = 1
	}

	// Get array bounds
	var LowerBounds int32
	var UpperBounds int32

	LowerBounds, err = safeArrayGetLBound(s, index)
	if err != nil {
		return
	}

	UpperBounds, err = safeArrayGetUBound(s, index)
	if err != nil {
		return
	}

	totalElements = UpperBounds - LowerBounds + 1
	return
}

func (s *SafeArray) GetType() (varType uint16, err error) {
	return safeArrayGetVarType(s)
}

func (s *SafeArray) Release() error {
	return safeArrayDestroy(s)
}
