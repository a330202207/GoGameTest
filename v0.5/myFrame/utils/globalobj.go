package utils

import (
	"encoding/json"
	"gFrame/giface"
	"io/ioutil"
)

//存储一切有关框架的全局参数，供其他模块使用
//一些参数是可以通过json由用户进行配置

type GlobalObj struct {
	//Server
	TcpServer giface.IServer //当前全局Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	//框架
	Version        string //当前框架版本号
	MaxConn        int    //当前服务主机允许的最大连接数
	MaxPackageSize uint32 //当前框架数据包的最大值
}

//定义一个全局的对外Globalobject
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/gindex.json")
	if err != nil {
		panic(err)
	}
	//将json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//提供一个init方法，初始化当前的GlobalObject
func init() {

	//如果配置文件没有加载，默认值
	GlobalObject = &GlobalObj{
		Name:           "gInxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//加载用户自定义方法
	GlobalObject.Reload()
}
