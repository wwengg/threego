package sbus

import (
	"github.com/nsqio/go-nsq"
	"github.com/wwengg/threego/core/smsg"
)

type SMsg interface {
	GetMsgId() int32
	GetCmd() uint16 // Gets the ID of the message(获取消息ID)
	GetRet() uint16
	GetVersion() uint8
	GetSerializeType() smsg.SerializeType
	GetCompressType() smsg.CompressType
	GetMessageType() smsg.MessageType
	GetSeq() uint64
	GetMeta() map[string]string
	GetData() []byte // Gets the content of the message(获取消息内容)
	//
	SetNsqMessage(message *nsq.Message)
	GetNsqMessage() *nsq.Message
	//
	GetHasFrameDecoder() bool
	SetHasFrameDecoder(hasFrameDevoder bool)
}

type NSQMsg struct {
	//*Header
	//PkgLen        uint32  // nsq不存在粘包问题 不需要
	Cmd           uint16
	Ret           uint16
	Version       uint8
	SerializeType smsg.SerializeType
	CompressType  smsg.CompressType
	MessageType   smsg.MessageType
	Seq           uint64
	Metadata      map[string]string
	Data          []byte
	nsqMessage    *nsq.Message
}

func NewNSQMsg(Cmd uint16, ret uint16, sType smsg.SerializeType, md map[string]string, data []byte) *NSQMsg {
	return &NSQMsg{
		Cmd:           Cmd,
		Ret:           ret,
		Version:       1,
		SerializeType: sType,
		CompressType:  smsg.Gzip,
		MessageType:   smsg.Response,
		Seq:           123456789812322,
		Metadata:      md,
		Data:          data,
	}
}
func (m *NSQMsg) GetMsgId() int32 {
	return int32(m.Cmd)
}
func (m *NSQMsg) GetCmd() uint16 {
	return m.Cmd
}
func (m *NSQMsg) GetRet() uint16 {
	return m.Ret
}
func (m *NSQMsg) GetVersion() uint8 {
	return m.Version
}
func (m *NSQMsg) GetSerializeType() smsg.SerializeType {
	return m.SerializeType
}
func (m *NSQMsg) GetCompressType() smsg.CompressType {
	return m.CompressType
}
func (m *NSQMsg) GetMessageType() smsg.MessageType {
	return m.MessageType
}

func (m *NSQMsg) GetSeq() uint64 {
	return m.Seq
}

func (m *NSQMsg) GetMeta() map[string]string {
	return m.Metadata
}

func (m *NSQMsg) GetData() []byte {
	return m.Data
}

func (m *NSQMsg) SetNsqMessage(message *nsq.Message) {
	m.nsqMessage = message
}

func (m *NSQMsg) GetNsqMessage() *nsq.Message {
	return m.nsqMessage
}

func (m *NSQMsg) GetHasFrameDecoder() bool                { return false }
func (m *NSQMsg) SetHasFrameDecoder(hasFrameDevoder bool) {}
