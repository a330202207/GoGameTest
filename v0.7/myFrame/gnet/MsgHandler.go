package gnet

import (
	"fmt"
	"gFrame/giface"
	"strconv"
)

//消息处理模块的实现
type MsgHandler struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]giface.IRouter
}

//初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]giface.IRouter),
	}
}

//调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request giface.IRequest) {
	//从Request中找到MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("API MsgID=", request.GetMsgID(), "Is Not Fount! Need Register")
	}

	//根据MsgID调用对应Router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router giface.IRouter) {
	//判断 当前Msg绑定的API的处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("Repeat API, MsgID=" + strconv.Itoa(int(msgID)))
	}

	//添加Msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add API MsgID=", msgID, " Success")
}
