package smsg

import (
	"github.com/smallnest/rpcx/protocol"
	"github.com/wwengg/threego/core/utils"
)

type SerializeType byte

const (
	// SerializeNone uses raw []byte and don't serialize/deserialize
	SerializeNone SerializeType = iota
	// JSON for payload.
	JSON
	// ProtoBuffer for payload.
	ProtoBuffer
)

// CompressType defines decompression type.
type CompressType byte

const (
	// None does not compress.
	None CompressType = iota
	// Gzip uses gzip compression.
	Gzip
	Brotli
)

// MessageType is message type of requests and responses.
type MessageType byte

const (
	// Request is message type of request
	Request MessageType = iota
	// Response is message type of response
	Response
)

var Compressors = map[CompressType]protocol.Compressor{
	None:   &protocol.RawDataCompressor{},
	Gzip:   &protocol.GzipCompressor{},
	Brotli: &utils.BrotliCompressor{},
}
