package gnet

import "gFrame/giface"

//实现 Router 时，先嵌入这个BaseRouter基类，然后根据需要对这个基类方法进行重新
type BaseRouter struct{}

//处理 conn 业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request giface.IRequest) {}

//处理 conn 业务主方法Hook
func (br *BaseRouter) Handle(request giface.IRequest) {}

//处理 conn 业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request giface.IRequest) {}
