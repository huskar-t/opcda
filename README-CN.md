# Go OPC DA 客户端

[English](README.md) | 简体中文

这是一个用 Go 语言编写的 OPC DA 客户端，它允许你与 OPC DA 服务器进行通信并获取数据。OPC DA
是一种常用的工业自动化通信协议，它允许设备和控制系统之间的数据交换。

## 功能

- 获取全部 OPC DA 服务器
- 连接到 OPC DA 服务器
- 浏览 OPC DA 服务器上的标签
- 同步读取标签的值
- 异步读取标签的值
- 同步写入标签的值
- 异步写入标签的值
- 订阅标签的实时数据变化

## 先决条件

在开始使用本客户端之前，确保满足以下先决条件：

- amd64/i386 架构的 Windows 操作系统
- Go 版本 1.20 或更高版本

**Go 1.20 是支持 Microsoft Windows 7 / 8 / Server 2008 / Server 2012 的最后一个版本，为了保证兼容性，会一直保持对 Go 1.20
的支持。**

**测试使用 Go1.20 和 Go 最新版本并进行 32 位和 64 位测试。**

## 安装

使用以下命令安装本客户端：

```shell
go get github.com/huskar-t/opcda
```

## 类型

本客户端提供了以下类型支持：

| OPC 类型            | GO 类型       | 说明          |
|-------------------|-------------|-------------|
| VT_BOOL           | bool        | 布尔值         |
| VT_I1             | int8        | 8 位有符号整数    |
| VT_I2             | int16       | 16 位有符号整数   |
| VT_I4             | int32       | 32 位有符号整数   |
| VT_I8             | int64       | 64 位有符号整数   |
| VT_UI1            | uint8       | 8 位无符号整数    |
| VT_UI2            | uint16      | 16 位无符号整数   |
| VT_UI4            | uint32      | 32 位无符号整数   |
| VT_UI8            | uint64      | 64 位无符号整数   |
| VT_R4             | float32     | 32 位浮点数     |
| VT_R8             | float64     | 64 位浮点数     |
| VT_BSTR           | string      | 字符串         |
| VT_DATE           | time.Time   | 日期时间        |
| VT_ARRAY\|VT_BOOL | []bool      | 布尔值数组       | 
| VT_ARRAY\|VT_I1   | []int8      | 8 位有符号整数数组  |
| VT_ARRAY\|VT_I2   | []int16     | 16 位有符号整数数组 |
| VT_ARRAY\|VT_I4   | []int32     | 32 位有符号整数数组 |
| VT_ARRAY\|VT_I8   | []int64     | 64 位有符号整数数组 |
| VT_ARRAY\|VT_UI1  | []uint8     | 8 位无符号整数数组  |
| VT_ARRAY\|VT_UI2  | []uint16    | 16 位无符号整数数组 |
| VT_ARRAY\|VT_UI4  | []uint32    | 32 位无符号整数数组 |
| VT_ARRAY\|VT_UI8  | []uint64    | 64 位无符号整数数组 |
| VT_ARRAY\|VT_R4   | []float32   | 32 位浮点数数组   |
| VT_ARRAY\|VT_R8   | []float64   | 64 位浮点数数组   |
| VT_ARRAY\|VT_BSTR | []string    | 字符串数组       |
| VT_ARRAY\|VT_DATE | []time.Time | 日期时间数组      |

其他类型暂未支持。

## 使用示例

- [获取全部 OPC DA 服务器](./example/serverlist)
- [浏览标签](./example/browse)
- [订阅标签](./example/subscribe)
- [同步读取标签](./example/read)
- [异步读取标签](./example/asyncread)
- [同步写入标签](./example/write)
- [异步写入标签](./example/asyncwrite)

## API 文档

全部 API 见 [API 文档](https://pkg.go.dev/github.com/huskar-t/opcda)

## 为什么选择本客户端

1. 目前 Go 语言生态中并没有成熟的 OPC DA 客户端，本客户端是为了填补这一空白而开发的。
2. 本客户端使用纯 Go 语言编写并且支持 32 位与 64 位。
3. 本客户端还提供了相对完整的功能，包括读取、写入和订阅标签的实时数据变化等。

### 为什么不是 [konimarti/opc](https://github.com/konimarti/opc)

1. [konimarti/opc](https://github.com/konimarti/opc) 使用 OPC DA Automation Wrapper 接口，无法对 Wrapper 进行修改无法保证
   bug 的修复。例如使用 Graybox DA Automation Wrapper 无法编译 64 位运行（HIGHENTROPYVA 默认开启，Wrapper 内部内存地址转换越界）
2. 读字符串类型存在内存泄漏问题。
3. 功能不够完善，如无法订阅标签的实时数据变化。
4. 未升级 [go-ole](https://github.com/go-ole/go-ole) [#252](https://github.com/go-ole/go-ole/pull/252)，无法获得毫秒时间戳

## 常见问题

1. 跨平台支持

   本客户端目前仅支持 Windows 操作系统。

2. 多平台编译

   本客户端无法在非 Windows 平台编译，因为它依赖于 Windows 平台的 COM 接口。如果程序需要支持多个平台，可以封装接口，并将非
   Windows 平台的接口设置为空实现。

3. 内存泄漏

   本客户端使用了 COM 接口，内部已经处理了内存释放，并进行了所有支持类型的疲劳测试。但是在使用过程中，如果发现内存泄漏问题，可以提交
   issue 并提供复现步骤。
