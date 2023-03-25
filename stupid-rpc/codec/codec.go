package codec

import "io"

type Header struct {
	Seqnum        uint64
	ServiceMethod string
	Error         string
}

/*
header包含 方法名，调用序列号，和错误，都由基本类型组成
*/

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(header *Header, data interface{}) error
}

type Type string

const (
	GobType  = "application/gob"
	JsonType = "application/json"
)

type NewCodecFunc func(io.ReadWriteCloser) Codec

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
