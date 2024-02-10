package com

import (
	"fmt"
	"math"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func (v *VARIANT) Clear() error {
	//if v.IsArray() {
	//	safeArray := (*SafeArray)(unsafe.Pointer(uintptr(v.Val)))
	//	totalElements, _ := safeArray.TotalElements(0)
	//	vt, _ := safeArray.GetType()
	//	switch vt {
	//	case uint16(VT_BSTR):
	//
	//		var element *uint16
	//		for i := int32(0); i < totalElements; i++ {
	//			safeArrayGetElement(safeArray, i, unsafe.Pointer(&element))
	//			SysFreeString(element)
	//		}
	//	}
	//}
	return VariantClear(v)
}

func (v *VARIANT) IsArray() bool {
	return v.VT&VT_ARRAY == VT_ARRAY
}

func (v *VARIANT) Value() interface{} {
	if v.VT == VT_EMPTY || v.VT == VT_NULL {
		return nil
	}
	if v.IsArray() {
		safeArray := (*SafeArray)(unsafe.Pointer(uintptr(v.Val)))
		values, err := safeArray.ToValueArray()
		if err != nil {
			panic(err)
		}
		return values
	}
	switch v.VT {
	case VT_I1:
		return int8(v.Val)
	case VT_UI1:
		return uint8(v.Val)
	case VT_I2:
		return int16(v.Val)
	case VT_UI2:
		return uint16(v.Val)
	case VT_I4:
		return int32(v.Val)
	case VT_UI4:
		return uint32(v.Val)
	case VT_I8:
		return int64(v.Val)
	case VT_UI8:
		return uint64(v.Val)
	case VT_INT:
		return int(v.Val)
	case VT_UINT:
		return uint(v.Val)
	case VT_R4:
		return *(*float32)(unsafe.Pointer(&v.Val))
	case VT_R8:
		return *(*float64)(unsafe.Pointer(&v.Val))
	case VT_BSTR:
		return windows.UTF16PtrToString(*(**uint16)(unsafe.Pointer(&v.Val)))
	case VT_DATE:
		d := uint64(v.Val)
		date, err := GetVariantDate(d)
		if err != nil {
			panic(err)
		}
		return date
	//case VT_CY:
	//	return float64(v.Val) / 10000
	case VT_BOOL:
		return (v.Val & 0xffff) != 0
	}
	return nil
}

type VariantWrapper struct {
	Variant *VARIANT
	str     []*uint16
}

func NewVariant(val interface{}) (*VariantWrapper, error) {
	v := &VariantWrapper{Variant: &VARIANT{}}
	err := v.SetValue(val)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (vw *VariantWrapper) SetValue(val interface{}) error {
	v := vw.Variant
	switch val.(type) {
	case int8:
		v.VT = VT_I1
		v.Val = int64(val.(int8))
	case []int8:
		v.VT = VT_ARRAY | VT_I1
		values := val.([]int8)
		array, _ := safeArrayCreateVector(VT_I1, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case uint8:
		v.VT = VT_UI1
		v.Val = int64(val.(uint8))
	case []uint8:
		v.VT = VT_ARRAY | VT_UI1
		values := val.([]uint8)
		array, _ := safeArrayCreateVector(VT_UI1, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case int16:
		v.VT = VT_I2
		v.Val = int64(val.(int16))

	case []int16:
		v.VT = VT_ARRAY | VT_I2
		values := val.([]int16)
		array, _ := safeArrayCreateVector(VT_I2, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case uint16:
		v.VT = VT_UI2
		v.Val = int64(val.(uint16))
	case []uint16:
		v.VT = VT_ARRAY | VT_UI2
		values := val.([]uint16)
		array, _ := safeArrayCreateVector(VT_UI2, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case int32:
		v.VT = VT_I4
		v.Val = int64(val.(int32))
	case []int32:
		v.VT = VT_ARRAY | VT_I4
		values := val.([]int32)
		array, _ := safeArrayCreateVector(VT_I4, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case uint32:
		v.VT = VT_UI4
		v.Val = int64(val.(uint32))
	case []uint32:
		v.VT = VT_ARRAY | VT_UI4
		values := val.([]uint32)
		array, _ := safeArrayCreateVector(VT_UI4, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case int64:
		v.VT = VT_I8
		v.Val = int64(val.(int64))
	case []int64:
		v.VT = VT_ARRAY | VT_I8
		values := val.([]int64)
		array, _ := safeArrayCreateVector(VT_I8, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case uint64:
		v.VT = VT_UI8
		v.Val = int64(val.(uint64))
	case []uint64:
		v.VT = VT_ARRAY | VT_UI8
		values := val.([]uint64)
		array, _ := safeArrayCreateVector(VT_UI8, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case int:
		v.VT = VT_INT
		v.Val = int64(val.(int))

	case []int:
		v.VT = VT_ARRAY | VT_INT
		values := val.([]int)
		array, _ := safeArrayCreateVector(VT_INT, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case uint:
		v.VT = VT_UINT
		v.Val = int64(val.(uint))
	case []uint:
		v.VT = VT_ARRAY | VT_UINT
		values := val.([]uint)
		array, _ := safeArrayCreateVector(VT_UINT, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case float32:
		v.VT = VT_R4
		v.Val = int64(math.Float32bits(val.(float32)))
	case []float32:
		v.VT = VT_ARRAY | VT_R4
		values := val.([]float32)
		array, _ := safeArrayCreateVector(VT_R4, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case float64:
		v.VT = VT_R8
		v.Val = int64(math.Float64bits(val.(float64)))
	case []float64:
		v.VT = VT_ARRAY | VT_R8
		values := val.([]float64)
		array, _ := safeArrayCreateVector(VT_R8, 0, uint32(len(values)))
		for i, value := range values {
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&value)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))

	case string:
		v.VT = VT_BSTR
		ptr := SysAllocStringLen(val.(string))
		v.Val = int64(uintptr(unsafe.Pointer(ptr)))
	case []string:
		v.VT = VT_ARRAY | VT_BSTR
		values := val.([]string)
		array, _ := safeArrayCreateVector(VT_BSTR, 0, uint32(len(values)))
		for i, value := range values {
			ptr := SysAllocStringLen(value)
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(ptr)))
			vw.str = append(vw.str, ptr)
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case time.Time:
		v.VT = VT_DATE
		date, err := TimeToVariantDate(val.(time.Time))
		if err != nil {
			return err
		}
		v.Val = int64(date)
	case []time.Time:
		v.VT = VT_ARRAY | VT_DATE
		values := val.([]time.Time)
		array, _ := safeArrayCreateVector(VT_DATE, 0, uint32(len(values)))
		for i, value := range values {
			date, err := TimeToVariantDate(value)
			if err != nil {
				safeArrayDestroy(array)
				return err
			}
			safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&date)))
		}
		v.Val = int64(uintptr(unsafe.Pointer(array)))
	case bool:
		v.VT = VT_BOOL
		if val.(bool) {
			v.Val = -1
		} else {
			v.Val = 0
		}
	case []bool:
		v.VT = VT_ARRAY | VT_BOOL
		values := val.([]bool)
		array, _ := safeArrayCreateVector(VT_BOOL, 0, uint32(len(values)))
		t := int16(-1)
		f := int16(0)
		for i, value := range values {
			if value {
				safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&(t))))
			} else {
				safeArrayPutElement(array, int64(i), uintptr(unsafe.Pointer(&(f))))
			}
		}
	default:
		return fmt.Errorf("unsupported type: %T", val)
	}
	return nil
}

func (vw *VariantWrapper) Clear() error {
	if len(vw.str) > 0 {
		for _, ptr := range vw.str {
			SysFreeString(ptr)
		}
	}
	return vw.Variant.Clear()
}
