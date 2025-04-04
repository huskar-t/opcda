package com

type VT uint16

const (
	VT_EMPTY            VT = 0
	VT_NULL             VT = 1
	VT_I2               VT = 2
	VT_I4               VT = 3
	VT_R4               VT = 4
	VT_R8               VT = 5
	VT_CY               VT = 6
	VT_DATE             VT = 7
	VT_BSTR             VT = 8
	VT_DISPATCH         VT = 9
	VT_ERROR            VT = 10
	VT_BOOL             VT = 11
	VT_VARIANT          VT = 12
	VT_UNKNOWN          VT = 13
	VT_DECIMAL          VT = 14
	VT_I1               VT = 16
	VT_UI1              VT = 17
	VT_UI2              VT = 18
	VT_UI4              VT = 19
	VT_I8               VT = 20
	VT_UI8              VT = 21
	VT_INT              VT = 22
	VT_UINT             VT = 23
	VT_VOID             VT = 24
	VT_HRESULT          VT = 25
	VT_PTR              VT = 26
	VT_SAFEARRAY        VT = 27
	VT_CARRAY           VT = 28
	VT_USERDEFINED      VT = 29
	VT_LPSTR            VT = 30
	VT_LPWSTR           VT = 31
	VT_RECORD           VT = 36
	VT_INT_PTR          VT = 37
	VT_UINT_PTR         VT = 38
	VT_FILETIME         VT = 64
	VT_BLOB             VT = 65
	VT_STREAM           VT = 66
	VT_STORAGE          VT = 67
	VT_STREAMED_OBJECT  VT = 68
	VT_STORED_OBJECT    VT = 69
	VT_BLOB_OBJECT      VT = 70
	VT_CF               VT = 71
	VT_CLSID            VT = 72
	VT_VERSIONED_STREAM VT = 73
	VT_BSTR_BLOB        VT = 0xfff
	VT_VECTOR           VT = 0x1000
	VT_ARRAY            VT = 0x2000
	VT_BYREF            VT = 0x4000
	VT_RESERVED         VT = 0x8000
	VT_ILLEGAL          VT = 0xffff
	VT_ILLEGALMASKED    VT = 0xfff
	VT_TYPEMASK         VT = 0xfff
)

const (
	S_OK           = 0x00000000
	E_UNEXPECTED   = 0x8000FFFF
	E_NOTIMPL      = 0x80004001
	E_OUTOFMEMORY  = 0x8007000E
	E_INVALIDARG   = 0x80070057
	E_NOINTERFACE  = 0x80004002
	E_POINTER      = 0x80004003
	E_HANDLE       = 0x80070006
	E_ABORT        = 0x80004004
	E_FAIL         = 0x80004005
	E_ACCESSDENIED = 0x80070005
	E_PENDING      = 0x8000000A

	CO_E_CLASSSTRING = 0x800401F3
)

// authentication level constants
const (
	RPC_C_AUTHN_LEVEL_DEFAULT       uint32 = 0
	RPC_C_AUTHN_LEVEL_NONE          uint32 = 1
	RPC_C_AUTHN_LEVEL_CONNECT       uint32 = 2
	RPC_C_AUTHN_LEVEL_CALL          uint32 = 3
	RPC_C_AUTHN_LEVEL_PKT           uint32 = 4
	RPC_C_AUTHN_LEVEL_PKT_INTEGRITY uint32 = 5
	RPC_C_AUTHN_LEVEL_PKT_PRIVACY   uint32 = 6
)

// impersonation level constants
const (
	RPC_C_IMP_LEVEL_DEFAULT     uint32 = 0
	RPC_C_IMP_LEVEL_ANONYMOUS   uint32 = 1
	RPC_C_IMP_LEVEL_IDENTIFY    uint32 = 2
	RPC_C_IMP_LEVEL_IMPERSONATE uint32 = 3
	RPC_C_IMP_LEVEL_DELEGATE    uint32 = 4
)

// EOLE_AUTHENTICATION_CAPABILITIES
const (
	EOAC_NONE              uint32 = 0
	EOAC_MUTUAL_AUTH       uint32 = 0x1
	EOAC_STATIC_CLOAKING   uint32 = 0x20
	EOAC_DYNAMIC_CLOAKING  uint32 = 0x40
	EOAC_ANY_AUTHORITY     uint32 = 0x80
	EOAC_MAKE_FULLSIC      uint32 = 0x100
	EOAC_DEFAULT           uint32 = 0x800
	EOAC_SECURE_REFS       uint32 = 0x2
	EOAC_ACCESS_CONTROL    uint32 = 0x4
	EOAC_APPID             uint32 = 0x8
	EOAC_DYNAMIC           uint32 = 0x10
	EOAC_REQUIRE_FULLSIC   uint32 = 0x200
	EOAC_AUTO_IMPERSONATE  uint32 = 0x400
	EOAC_DISABLE_AAA       uint32 = 0x1000
	EOAC_NO_CUSTOM_MARSHAL uint32 = 0x2000
	EOAC_RESERVED1         uint32 = 0x4000
)
