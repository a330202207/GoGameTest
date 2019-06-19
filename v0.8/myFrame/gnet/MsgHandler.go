package gnet

import (
	"fmt"
	"gFrame/giface"
	"gFrame/utils"
	"strconv"
)

//消息处理模块的实现
type MsgHandler struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]giface.IRouter

	//负责Worker取任务的消息队列
	TaskQueue []chan giface.IRequest

	//业务工作Worker池的工作数量
	WorkerPoolSize uint32
}

//初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]giface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan giface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

//启动一个Worker工作池(开启工作池的动作只能发生一次，框架只能有一个Worker工作池)
func (mh *MsgHandler) StartWorkerPool() {
	//根据WorkerPoolSize，分别开启Worker，每个Worker用一个Go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个Worker被启动
		//当前的Worker对应的Channel消息队列，开辟空间,第0个Worker，就用第0个Channel
		mh.TaskQueue[i] = make(chan giface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		//启动当前的Worker，阻塞等待消息从Channel传递进来

		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan giface.IRequest) {
	fmt.Println("WorkerID = ", workerID, " Is Started...")

	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给TaskQueue，由Worker处理
func (mh *MsgHandler) SendMsgToTaskQueue(request giface.IRequest) {
	//将消息平均分配给不通过的Worker
	//根据客户端简历的ConnID来进行分配
	//轮询法则
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" Request MsgID = ", request.GetMsgID(), " To WorkerID = ", workerID)

	//将消息发送给对应的Worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
