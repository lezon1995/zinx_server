package zimpl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/util"
	"zinx/zinterface"
)

/**
针对IMessage进行TLV(type length value)格式的封装
解决封包/拆包/粘包问题
*/
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//dataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

/**
| dataLen | msgId | msg |
*/
func (d *DataPack) Pack(msg zinterface.IMessage) ([]byte, error) {
	//类似于netty中的byteBuffer 包装byte
	buffer := bytes.NewBuffer([]byte{})
	//低位写入
	binary.Write(buffer, binary.LittleEndian, msg.GetMsgLen())
	binary.Write(buffer, binary.LittleEndian, msg.GetMsgId())
	binary.Write(buffer, binary.LittleEndian, msg.GetData())
	return buffer.Bytes(), nil
}

func (d *DataPack) Unpack(rawBinaryData []byte) (zinterface.IMessage, error) {
	//创建一个从二进制数据读取的ioReader
	reader := bytes.NewReader(rawBinaryData)
	message := NewMessage()
	//读dataLen
	binary.Read(reader, binary.LittleEndian, &message.msgLen)
	//读msgId
	binary.Read(reader, binary.LittleEndian, &message.msgId)

	//判断dataLen是否超出配置的最大包长度
	if message.msgLen > util.GlobalConfig.BufferSize {
		return nil, errors.New("too large msg msg received")
	}
	return message, nil
}
