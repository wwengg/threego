package sbus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"runtime"

	"github.com/smallnest/rpcx/util"
	"github.com/wwengg/threego/core/slog"
	"github.com/wwengg/threego/core/smsg"
)

var nsqDataHeaderLen uint32 = 13

type NsqDataPack struct {
}

var NsqDataPackObj = new(NsqDataPack)

var ErrMetaKVMissing = errors.New("wrong metadata lines. some keys or values are missing")

//+------+-------+---------+---------------+--------------+-------------+-------+------------+--------------+-------------+------------+
//| CMD  |  Ret  | version | SerializeType | CompressType | messageType |  seq  |  meta len  |   meta data  |   data len  |    data    |
//| 2字节 |  2字节 |  1字节  |     4bit      |     2bit     |      2bit   | 8字节  |    4字节    |     n字节    |      4字节   |    n字节    |
//+------+-------+---------+---------------+--------------+-------------+-------+------------+--------------+-------------+------------+
//|                                 header                                      |
//|                                 13字节                                       |
//+-----------------------------------------------------------------------------+

func NewNsqDataPack() SDataPack { return &NsqDataPack{} }

func (dp *NsqDataPack) GetHeadLen() uint32 {
	//ID uint32(4 bytes) +  DataLen uint32(4 bytes)
	return nsqDataHeaderLen
}

// Pack packs the message (compresses the data)
// (封包方法,压缩数据)
func (dp *NsqDataPack) Pack(msg SMsg) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			slog.Ins().Errorf("panic in message decode: %v, stack: %s", err, errStack[:n])
		}
	}()
	// Create a buffer to store the bytes
	// (创建一个存放bytes字节的缓冲)
	dataBuff := bytes.NewBuffer([]byte{})

	// Write the cmd
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetCmd()); err != nil {
		return nil, err
	}

	// Write the ret
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetRet()); err != nil {
		return nil, err
	}

	// Write the version
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetVersion()); err != nil {
		return nil, err
	}
	var oneByte [1]byte
	// SerializeType
	oneByte[0] = (oneByte[0] &^ 0xF0) | (byte(msg.GetSerializeType()) << 4)
	// CompressType
	oneByte[0] = (oneByte[0] &^ 0x0C) | ((byte(msg.GetCompressType()) << 2) & 0x0C)
	// messageType
	oneByte[0] = (oneByte[0] &^ 0x03) | (byte(msg.GetMessageType()) & 0x03)
	// Write the oneByte
	if err := binary.Write(dataBuff, binary.BigEndian, oneByte); err != nil {
		return nil, err
	}
	// Write the seq
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetSeq()); err != nil {
		return nil, err
	}
	// Write the meta
	var bb = bytes.NewBuffer(make([]byte, 0, len(msg.GetMeta())*64))
	EncodeMetadata(msg.GetMeta(), bb)
	meta := bb.Bytes()
	// Write the meta len
	if err := binary.Write(dataBuff, binary.BigEndian, uint32(len(meta))); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, meta); err != nil {
		return nil, err
	}
	// Write the data
	if err := binary.Write(dataBuff, binary.BigEndian, uint32(len(msg.GetData()))); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack unpacks the message (decompresses the data)
// (拆包方法,解压数据)
func (dp *NsqDataPack) Unpack(binaryData []byte) (SMsg, error) {
	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			slog.Ins().Errorf("panic in message decode: %v, stack: %s", err, errStack[:n])

		}
	}()
	// Create an ioReader for the input binary data
	dataBuff := bytes.NewReader(binaryData)

	// Only unpack the header information to obtain the data length and message ID
	// (只解压head的信息，得到dataLen和msgID)
	msg := &NSQMsg{}

	// Read the Cmd
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Cmd); err != nil {
		return nil, err
	}

	// Read the Ret
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Ret); err != nil {
		return nil, err
	}
	// Read the Version
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Version); err != nil {
		return nil, err
	}
	onebyte := [1]byte{}
	if err := binary.Read(dataBuff, binary.BigEndian, &onebyte); err != nil {
		return nil, err
	}
	// Read the SerializeType
	msg.SerializeType = smsg.SerializeType((onebyte[0] & 0xF0) >> 4)
	// Read the CompressType
	msg.CompressType = smsg.CompressType((onebyte[0] & 0x0C) >> 2)
	// Read the MessageType
	msg.MessageType = smsg.MessageType(onebyte[0] & 0x03)
	// Read the seq
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Seq); err != nil {
		return nil, err
	}
	// Read the meta
	var metaLen uint32
	if err := binary.Read(dataBuff, binary.BigEndian, &metaLen); err != nil {
		return nil, err
	}
	if metaLen > 0 {
		metaData := make([]byte, metaLen)
		// Read the metaData
		if err := binary.Read(dataBuff, binary.BigEndian, &metaData); err != nil {
			return nil, err
		}
		if m, err := DecodeMetadata(metaLen, metaData); err != nil {
			return nil, err
		} else {
			msg.Metadata = m
		}
	}
	// Read the data
	var dataLen uint32
	if err := binary.Read(dataBuff, binary.BigEndian, &dataLen); err != nil {
		return nil, err
	}
	if dataLen > 0 {
		// Read the metaData
		// 包大小可能大于65535 所以全都返回吧
		msg.Data = binaryData[22+metaLen:]
	}

	// Only the header data needs to be unpacked, and then another data read is performed from the connection based on the header length
	// (这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据)
	return msg, nil
}

// len,string,len,string,......
func EncodeMetadata(m map[string]string, bb *bytes.Buffer) {
	if len(m) == 0 {
		return
	}
	d := make([]byte, 4)
	for k, v := range m {
		binary.BigEndian.PutUint32(d, uint32(len(k)))
		bb.Write(d)
		bb.Write(util.StringToSliceByte(k))
		binary.BigEndian.PutUint32(d, uint32(len(v)))
		bb.Write(d)
		bb.Write(util.StringToSliceByte(v))
	}
}

func DecodeMetadata(l uint32, data []byte) (map[string]string, error) {
	m := make(map[string]string, 10)
	n := uint32(0)
	for n < l {
		// parse one key and value
		// key
		sl := binary.BigEndian.Uint32(data[n : n+4])
		n = n + 4
		if n+sl > l-4 {
			return m, ErrMetaKVMissing
		}
		k := string(data[n : n+sl])
		n = n + sl

		// value
		sl = binary.BigEndian.Uint32(data[n : n+4])
		n = n + 4
		if n+sl > l {
			return m, ErrMetaKVMissing
		}
		v := string(data[n : n+sl])
		n = n + sl
		m[k] = v
	}

	return m, nil
}
