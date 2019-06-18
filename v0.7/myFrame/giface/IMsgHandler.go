package giface

//消息管理抽象层
type IMsgHandler interface {
	//调度/执行对应的Router消息处理方法
	DoMsgHandler(reques IRequest)

	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
