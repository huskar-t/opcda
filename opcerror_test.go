package opcda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOPCError_Error(t *testing.T) {
	type fields struct {
		ErrorCode    int32
		ErrorMessage string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "TestOPCError_Error",
			fields: fields{
				ErrorCode:    int32(16),
				ErrorMessage: "Unspecified error",
			},
			want: "OPCError [0x10]: Unspecified error",
		},
		{
			name: "TestOPCError_Error",
			fields: fields{
				ErrorCode:    int32(-1073479679),
				ErrorMessage: "",
			},
			want: "OPCError [0xc0040001]: The value of the handle is invalid",
		},
		{
			name: "TestOPCError_Error",
			fields: fields{
				ErrorCode:    int32(-1),
				ErrorMessage: "",
			},
			want: "OPCError [0xffffffff]: unknown error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &OPCError{
				ErrorCode:    tt.fields.ErrorCode,
				ErrorMessage: tt.fields.ErrorMessage,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}
