// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

package types

import "github.com/golang/protobuf/proto"
import "fmt"
import "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Log struct {
	// 日志级别，支持debug(dbug)/info/warn/error(eror)/crit
	Loglevel        string `protobuf:"bytes,1,opt,name=loglevel" json:"loglevel,omitempty"`
	LogConsoleLevel string `protobuf:"bytes,2,opt,name=logConsoleLevel" json:"logConsoleLevel,omitempty"`
	// 日志文件名，可带目录，所有生成的日志文件都放到此目录下
	LogFile string `protobuf:"bytes,3,opt,name=logFile" json:"logFile,omitempty"`
	// 单个日志文件的最大值（单位：兆）
	MaxFileSize uint32 `protobuf:"varint,4,opt,name=maxFileSize" json:"maxFileSize,omitempty"`
	// 最多保存的历史日志文件个数
	MaxBackups uint32 `protobuf:"varint,5,opt,name=maxBackups" json:"maxBackups,omitempty"`
	// 最多保存的历史日志消息（单位：天）
	MaxAge uint32 `protobuf:"varint,6,opt,name=maxAge" json:"maxAge,omitempty"`
	// 日志文件名是否使用本地事件（否则使用UTC时间）
	LocalTime bool `protobuf:"varint,7,opt,name=localTime" json:"localTime,omitempty"`
	// 历史日志文件是否压缩（压缩格式为gz）
	Compress bool `protobuf:"varint,8,opt,name=compress" json:"compress,omitempty"`
	// 是否打印调用源文件和行号
	CallerFile bool `protobuf:"varint,9,opt,name=callerFile" json:"callerFile,omitempty"`
	// 是否打印调用方法
	CallerFunction bool `protobuf:"varint,10,opt,name=callerFunction" json:"callerFunction,omitempty"`
}

func (m *Log) Reset()                    { *m = Log{} }
func (m *Log) String() string            { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()               {}
func (*Log) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *Log) GetLoglevel() string {
	if m != nil {
		return m.Loglevel
	}
	return ""
}

func (m *Log) GetLogConsoleLevel() string {
	if m != nil {
		return m.LogConsoleLevel
	}
	return ""
}

func (m *Log) GetLogFile() string {
	if m != nil {
		return m.LogFile
	}
	return ""
}

func (m *Log) GetMaxFileSize() uint32 {
	if m != nil {
		return m.MaxFileSize
	}
	return 0
}

func (m *Log) GetMaxBackups() uint32 {
	if m != nil {
		return m.MaxBackups
	}
	return 0
}

func (m *Log) GetMaxAge() uint32 {
	if m != nil {
		return m.MaxAge
	}
	return 0
}

func (m *Log) GetLocalTime() bool {
	if m != nil {
		return m.LocalTime
	}
	return false
}

func (m *Log) GetCompress() bool {
	if m != nil {
		return m.Compress
	}
	return false
}

func (m *Log) GetCallerFile() bool {
	if m != nil {
		return m.CallerFile
	}
	return false
}

func (m *Log) GetCallerFunction() bool {
	if m != nil {
		return m.CallerFunction
	}
	return false
}

func init() {
	proto.RegisterType((*Log)(nil), "types.Log")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 235 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcf, 0x4a, 0xc4, 0x30,
	0x10, 0x87, 0xe9, 0xae, 0xdb, 0x6d, 0xc7, 0x7f, 0x90, 0x83, 0x0c, 0x22, 0x52, 0x3c, 0x48, 0x4f,
	0x5e, 0x7c, 0x02, 0x15, 0x3c, 0xed, 0xa9, 0xfa, 0x02, 0x31, 0x8c, 0x21, 0x38, 0xe9, 0x94, 0xa6,
	0x2b, 0xd5, 0xb7, 0xf3, 0xcd, 0xa4, 0xc3, 0xba, 0x96, 0xbd, 0xf5, 0xfb, 0x7e, 0x1f, 0x34, 0x09,
	0x9c, 0x38, 0x69, 0xdf, 0x83, 0xbf, 0xeb, 0x7a, 0x19, 0xc4, 0xac, 0x86, 0xaf, 0x8e, 0xd2, 0xcd,
	0xcf, 0x02, 0x96, 0x1b, 0xf1, 0xe6, 0x12, 0x0a, 0x16, 0xcf, 0xf4, 0x49, 0x8c, 0x59, 0x95, 0xd5,
	0x65, 0xb3, 0x67, 0x53, 0xc3, 0x39, 0x8b, 0x7f, 0x92, 0x36, 0x09, 0xd3, 0x46, 0x93, 0x85, 0x26,
	0x87, 0xda, 0x20, 0xac, 0x59, 0xfc, 0x73, 0x60, 0xc2, 0xa5, 0x16, 0x7f, 0x68, 0x2a, 0x38, 0x8e,
	0x76, 0x9c, 0x3e, 0x5f, 0xc2, 0x37, 0xe1, 0x51, 0x95, 0xd5, 0xa7, 0xcd, 0x5c, 0x99, 0x6b, 0x80,
	0x68, 0xc7, 0x47, 0xeb, 0x3e, 0xb6, 0x5d, 0xc2, 0x95, 0x06, 0x33, 0x63, 0x2e, 0x20, 0x8f, 0x76,
	0x7c, 0xf0, 0x84, 0xb9, 0x6e, 0x3b, 0x32, 0x57, 0x50, 0xb2, 0x38, 0xcb, 0xaf, 0x21, 0x12, 0xae,
	0xab, 0xac, 0x2e, 0x9a, 0x7f, 0x31, 0xdd, 0xcb, 0x49, 0xec, 0x7a, 0x4a, 0x09, 0x0b, 0x1d, 0xf7,
	0x3c, 0xfd, 0xd1, 0x59, 0x66, 0xea, 0xf5, 0xc0, 0xa5, 0xae, 0x33, 0x63, 0x6e, 0xe1, 0x6c, 0x47,
	0xdb, 0xd6, 0x0d, 0x41, 0x5a, 0x04, 0x6d, 0x0e, 0xec, 0x5b, 0xae, 0x2f, 0x7a, 0xff, 0x1b, 0x00,
	0x00, 0xff, 0xff, 0x82, 0x02, 0x96, 0x7f, 0x61, 0x01, 0x00, 0x00,
}