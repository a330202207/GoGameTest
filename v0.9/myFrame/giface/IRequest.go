package giface

//将客户端请求的链接信息 和 请求的数据 包装到一个Request中
type IRequest interface {
	//得到当前链接
	GetConnection() IConnection

	//得到当前请求消息数据
	GetData() []byte

	//得到请求的消息ID
	GetMsgID() uint32
}
