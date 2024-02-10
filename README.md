# Go OPC DA Client

English | [简体中文](README-CN.md)

This is an OPC DA client written in Go language, allowing you to communicate with OPC DA servers and retrieve data. OPC DA is a commonly used industrial automation communication protocol that enables data exchange between devices and control systems.

## Features

- Get all OPC DA servers
- Connect to OPC DA server
- Browse tags on OPC DA server
- Synchronously read tag values
- Asynchronously read tag values
- Synchronously write tag values
- Asynchronously write tag values
- Subscribe to real-time data changes of tags

## Prerequisites

Before using this client, make sure you meet the following prerequisites:

- Windows operating system with amd64/i386 architecture
- Go version 1.20 or higher

**Go 1.20 is the last version that supports Microsoft Windows 7 / 8 / Server 2008 / Server 2012. To ensure compatibility, this client will continue to support Go 1.20.**

**Testing is done with both Go 1.20 and the latest version of Go, covering both 32-bit and 64-bit testing.**

## Installation

Use the following command to install this client:

```shell
go get github.com/huskar-t/opcda
```

## Types

This client provides support for the following types:

| OPC Type          | Go Type     | Description                        |
|-------------------|-------------|------------------------------------|
| VT_BOOL           | bool        | Boolean value                      |
| VT_I1             | int8        | 8-bit signed integer               |
| VT_I2             | int16       | 16-bit signed integer              |
| VT_I4             | int32       | 32-bit signed integer              |
| VT_I8             | int64       | 64-bit signed integer              |
| VT_UI1            | uint8       | 8-bit unsigned integer             |
| VT_UI2            | uint16      | 16-bit unsigned integer            |
| VT_UI4            | uint32      | 32-bit unsigned integer            |
| VT_UI8            | uint64      | 64-bit unsigned integer            |
| VT_R4             | float32     | 32-bit floating point number       |
| VT_R8             | float64     | 64-bit floating point number       |
| VT_BSTR           | string      | String                             |
| VT_DATE           | time.Time   | Date time                          |
| VT_ARRAY\|VT_BOOL | []bool      | Boolean array                      |
| VT_ARRAY\|VT_I1   | []int8      | 8-bit signed integer array         |
| VT_ARRAY\|VT_I2   | []int16     | 16-bit signed integer array        |
| VT_ARRAY\|VT_I4   | []int32     | 32-bit signed integer array        |
| VT_ARRAY\|VT_I8   | []int64     | 64-bit signed integer array        |
| VT_ARRAY\|VT_UI1  | []uint8     | 8-bit unsigned integer array       |
| VT_ARRAY\|VT_UI2  | []uint16    | 16-bit unsigned integer array      |
| VT_ARRAY\|VT_UI4  | []uint32    | 32-bit unsigned integer array      |
| VT_ARRAY\|VT_UI8  | []uint64    | 64-bit unsigned integer array      |
| VT_ARRAY\|VT_R4   | []float32   | 32-bit floating point number array |
| VT_ARRAY\|VT_R8   | []float64   | 64-bit floating point number array |
| VT_ARRAY\|VT_BSTR | []string    | String array                       |
| VT_ARRAY\|VT_DATE | []time.Time | Date time array                    |

Other types are not currently supported.

## Usage Examples

- [Get all OPC DA servers](./example/serverlist)
- [Browse tags](./example/browse)
- [Subscribe to tags](./example/subscribe)
- [Synchronously read tags](./example/read)
- [Asynchronously read tags](./example/asyncread)
- [Synchronously write tags](./example/write)
- [Asynchronously write tags](./example/asyncwrite)

## API Documentation

All APIs can be found in the [API documentation](https://pkg.go.dev/github.com/huskar-t/opcda).

## Why Choose This Client

1. There is currently no mature OPC DA client in the Go language ecosystem, and this client was developed to fill this gap.
2. This client is written purely in Go language and supports both 32-bit and 64-bit architectures.
3. This client also provides relatively complete functionality, including reading, writing, and subscribing to real-time data changes of tags.

### Why Not [konimarti/opc](https://github.com/konimarti/opc)

1. [konimarti/opc](https://github.com/konimarti/opc) uses the OPC DA Automation Wrapper interface, which cannot be modified for bug fixes. For example, using the Graybox DA Automation Wrapper cannot compile for 64-bit operation (HIGHENTROPYVA is enabled by default, leading to out-of-bounds memory address conversion within the Wrapper).
2. There are memory leaks when reading string types.
3. Insufficient functionality, such as inability to subscribe to real-time data changes of tags.
4. [go-ole](https://github.com/go-ole/go-ole) has not been updated [#252](https://github.com/go-ole/go-ole/pull/252), making it unable to obtain millisecond timestamps.

## Frequently Questions

1. Cross-platform support

   This client only supports the Windows operating system.

2. Multi-platform compilation

   This client cannot be compiled on non-Windows platforms because it relies on Windows platform COM interfaces. If the program needs to support multiple platforms, interfaces can be encapsulated, and non-Windows platform interfaces can be set to empty implementations.

3. Memory leaks

   This client uses COM interfaces, and memory release has been handled internally, with fatigue testing done for all supported types. However, if memory leak issues are found during use, you can submit an issue and provide reproduction steps.
