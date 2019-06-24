package gnet

import (
	"bytes"
	"encoding/binary"
	"gFrame/giface"
	"gFrame/utils"
	"github.com/pkg/errors"
)

//封包、拆包具体操作
//|DateLen|MsgId|Data|
type DataPack struct{}

//拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

//封包方法
func (dp *DataPack) Pack(msg giface.IMessage) ([]byte, error) {
	//创建一个存放byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将DataLen写入DataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将MsgId 写入DataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将Data数据写入DataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法(讲包的Head信息读出来)之后再根据Head信息里面的Data长度，再进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (giface.IMessage, error) {
	//创建一个从输入二进制数据的IoReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压Head信息，得到DataLen喝MsgID
	msg := &Message{}

	//读DataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断DataLen是否超过允许最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("Too Large Msg Data Recv!")
	}

	return msg, nil
}
