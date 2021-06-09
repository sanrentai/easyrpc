package easyrpc

import (
	"time"

	"github.com/sanrentai/easyrpc/codec"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber    int           // 幻数，它可以用来标记文件或者协议的格式
	CodecType      codec.Type    // client may choose different Codec to encode body
	ConnectTimeout time.Duration // 0 means no limit
	HandleTimeout  time.Duration
}

var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: time.Second * 10,
}

// 一般来说，涉及协议协商的这部分信息，需要设计固定的字节来传输的。
// 但是为了实现上更简单，GeeRPC 客户端固定采用 JSON 编码 Option，
// 后续的 header 和 body 的编码方式由 Option 中的 CodeType 指定，
// 服务端首先使用 JSON 解码 Option，然后通过 Option 的 CodeType
// 解码剩余的内容。即报文将以这样的形式发送：
// | Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
// | <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|

// 在一次连接中，Option 固定在报文的最开始，Header 和 Body 可以有多个，即报文可能是这样的。
// | Option | Header1 | Body1 | Header2 | Body2 | ...
