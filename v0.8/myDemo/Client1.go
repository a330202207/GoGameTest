package main

import (
	"fmt"
	"gFrame/gnet"
	"io"
	"net"
	"time"
)

//模拟客户端
func main() {

	fmt.Println("Client1 Start...")

	time.Sleep(1 * time.Second)

	//连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Client Start Error, Exit!")
		return
	}

	for {

		//发送封包的Message消息 MsgID:0
		dp := gnet.NewDataPack()
		binaryMsg, err := dp.Pack(gnet.NewMsgPackage(1, []byte("Client1 Test Message")))
		if err != nil {
			fmt.Println("Pack Error", err)
			return
		}

		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("Write Error", err)
			return
		}

		//服务器就应该回复一个Message数据，MsgID:1 pingpingping

		//先读取流中的Head部分得到ID 和DataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("Read Head Error", err)
			return
		}

		//将二进制的Head拆包到Msg 结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("Client Unpack MsgHead Error", err)
			return
		}

		if msgHead.GetMsgLen() > 0 {
			//再根据DataLen 进行二次读取，将Data读出来
			msg := msgHead.(*gnet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("Read Msg Data Error", err)
				return

			}

			fmt.Println("——》Recv Server Msg:ID=", msg.Id, ",Len=", msg.DataLen, ",Data=", string(msg.Data))
		}

		//CPU阻塞
		time.Sleep(1 * time.Second)

	}

}
