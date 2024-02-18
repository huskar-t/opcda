package com

import (
	"fmt"
	"time"
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

func (s *SafeArray) ToValueArray() (interface{}, error) {
	var err error
	totalElements, _ := s.TotalElements(0)
	vt, _ := safeArrayGetVarType(s)

	switch VT(vt) {
	case VT_BOOL:
		values := make([]bool, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v int16
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = (v & 0xff) != 0
		}
		return values, nil
	case VT_I1:
		values := make([]int8, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v int8
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_I2:
		values := make([]int16, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v int16
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_I4:
		values := make([]int32, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v int32
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_I8:
		values := make([]int64, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v int64
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_UI1:
		values := make([]uint8, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v uint8
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_UI2:
		values := make([]uint16, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v uint16
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_UI4:
		values := make([]uint32, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v uint32
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_UI8:
		values := make([]uint64, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v uint64
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_INT:
		values := make([]int, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v int
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_UINT:
		values := make([]uint, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v uint
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_R4:
		values := make([]float32, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v float32
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_R8:
		values := make([]float64, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var v float64
			err = safeArrayGetElement(s, i, unsafe.Pointer(&v))
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil
	case VT_BSTR:
		values := make([]string, totalElements)
		for i := int32(0); i < totalElements; i++ {
			var element *uint16
			err = safeArrayGetElement(s, i, unsafe.Pointer(&element))
			if err != nil {
				return nil, err
			}
			values[i] = windows.UTF16PtrToString(element)
			SysFreeString(element)
		}
		return values, nil
	case VT_DATE:
		values := make([]time.Time, totalElements)
		for i := int32(0); i < totalElements; i++ {
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
		}
		return values, nil
	default:
		return nil, fmt.Errorf("unknown value type %x", VT(vt))
	}
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
