package opcda

import "fmt"

type OPCError struct {
	ErrorCode    int32
	ErrorMessage string
}

func (e *OPCError) Error() string {
	if e.ErrorMessage == "" {
		if msg, ok := opcErrors[e.ErrorCode]; ok {
			return fmt.Errorf("OPCError [0x%x]: %s", uint32(e.ErrorCode), msg).Error()
		}
		return fmt.Errorf("OPCError [0x%x]: %s", uint32(e.ErrorCode), "unknown error").Error()
	}
	return fmt.Errorf("OPCError [0x%x]: %s", uint32(e.ErrorCode), e.ErrorMessage).Error()
}

var opcErrors = map[int32]string{
	int32(OPCInvalidHandle):   "The value of the handle is invalid",
	int32(OPCBadType):         "The server cannot convert the data between the specified format/ requested data type and the canonical data type",
	int32(OPCPublic):          "The requested operation cannot be done on a public group",
	int32(OPCBadRights):       "The Items AccessRights do not allow the operation",
	int32(OPCUnknownItemID):   "The item ID is not defined in the server address space (on add or validate) or no longer exists in the server address space (for read or write). ",
	int32(OPCInvalidItemID):   "The item ID doesn't conform to the server's syntax",
	int32(OPCInvalidFilter):   "The filter string was not valid",
	int32(OPCUnknownPath):     "The item's access path is not known to the server",
	int32(OPCRange):           "The value was out of range",
	int32(OPCDuplicateName):   "Duplicate name not allowed",
	int32(OPCUnsupportedRate): "The server does not support the requested data rate but will use the closest available rate",
	int32(OPCClamp):           "A value passed to WRITE was accepted but the output was clamped",
	int32(OPCInuse):           "The operation cannot be performed because the object is bering referenced",
	int32(OPCInvalidConfig):   "The server's configuration file is an invalid format",
	int32(OPCNotFound):        "Requested Object was not found",
	int32(OPCInvalidPID):      "The passed property ID is not valid for the item",
}

var (
	OPCInvalidHandle   = uint32(0xC0040001)
	OPCBadType         = uint32(0xC0040004)
	OPCPublic          = uint32(0xC0040005)
	OPCBadRights       = uint32(0xC0040006)
	OPCUnknownItemID   = uint32(0xC0040007)
	OPCInvalidItemID   = uint32(0xC0040008)
	OPCInvalidFilter   = uint32(0xC0040009)
	OPCUnknownPath     = uint32(0xC004000A)
	OPCRange           = uint32(0xC004000B)
	OPCDuplicateName   = uint32(0xC004000C)
	OPCUnsupportedRate = uint32(0x0004000D)
	OPCClamp           = uint32(0x0004000E)
	OPCInuse           = uint32(0x0004000F)
	OPCInvalidConfig   = uint32(0xC0040010)
	OPCNotFound        = uint32(0xC0040011)
	OPCInvalidPID      = uint32(0xC0040203)
)
