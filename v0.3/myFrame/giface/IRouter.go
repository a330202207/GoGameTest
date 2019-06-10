package giface

/*
	定义路由的抽象层
	路由里的数据都是 IRequest
*/
type IRouter interface {
	//处理 conn 业务之前的钩子方法Hook
	PreHandle(request IRequest)

	//处理 conn 业务主方法Hook
	Handle(request IRequest)

	//处理 conn 业务之后的钩子方法Hook
	PostHandle(request IRequest)
}
