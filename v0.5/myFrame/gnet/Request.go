package gnet

import "gFrame/giface"

type Request struct {
	//已经和客户端建立好的链接
	conn giface.IConnection

	//客户端请求的数据
	msg giface.IMessage
}

//得到当前链接
func (r *Request) GetConnection() giface.IConnection {
	return r.conn
}

//得到当前请求消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//得到请求的消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
