package opcda

import (
	"fmt"
	"reflect"
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

func TestOPCWrapperError_Error(t *testing.T) {
	type fields struct {
		Err  error
		Info string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test with error and info",
			fields: fields{
				Err:  fmt.Errorf("test error"),
				Info: "test info",
			},
			want: "test info: test error",
		},
		{
			name: "Test with error and no info",
			fields: fields{
				Err:  fmt.Errorf("test error"),
				Info: "",
			},
			want: ": test error",
		},
		{
			name: "Test with no error and info",
			fields: fields{
				Err:  nil,
				Info: "test info",
			},
			want: "test info: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &OPCWrapperError{
				Err:  tt.fields.Err,
				Info: tt.fields.Info,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("OPCWrapperError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOPCWrapperError(t *testing.T) {
	type args struct {
		info string
		err  error
	}
	tests := []struct {
		name string
		args args
		want *OPCWrapperError
	}{
		{
			name: "Test with error and info",
			args: args{
				err:  fmt.Errorf("test error"),
				info: "test info",
			},
			want: &OPCWrapperError{
				Err:  fmt.Errorf("test error"),
				Info: "test info",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOPCWrapperError(tt.args.info, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOPCWrapperError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOPCWrapperError_Unwrap(t *testing.T) {
	type fields struct {
		Err  error
		Info string
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "Test with error",
			fields: fields{
				Err:  fmt.Errorf("test error"),
				Info: "test info",
			},
			want: fmt.Errorf("test error"),
		},
		{
			name: "Test with no error",
			fields: fields{
				Err:  nil,
				Info: "test info",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &OPCWrapperError{
				Err:  tt.fields.Err,
				Info: tt.fields.Info,
			}
			if got := e.Unwrap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OPCWrapperError.Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}
