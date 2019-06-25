package gnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//需要注释

//负责测试DataPack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	//模拟的服务

	//创建SocketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Server listen err:", err)
		return
	}

	//创建Go 承载 负责从客户端处理业务
	go func() {
		//从客户端读取数据，拆包处理

		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("Server Accept err:", err)
				return
			}

			go func(conn net.Conn) {
				//处理客户端的请求

				//拆包过程

				dp := NewDataPack()
				for {
					//第一次从Conn读，把包的Head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("Read Head err:", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("Server Unpack err:", err)
						return
					}

					//msg是有数据的，需要第二次读取
					if msgHead.GetMsgLen() > 0 {
						//第二次从Conn读，根据Head中DataLen再读取Data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据DataLen的长度再次从Io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("Server Unpack Data err:", err)
							return
						}

						fmt.Println("——>Recv MsgID", msg.Id, "DataLen=", msg.DataLen, "Data=", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Client Dial err:", err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'H', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("Client Pack msg1 err:", err)
		return
	}

	//封装第一个msg1包
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'W', 'o', 'r', 'l', 'd'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("Client Pack msg2 err:", err)
		return
	}

	//将两个包黏在一起
	sendData1 = append(sendData1, sendData2...)

	//一次性方式给服务端
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}
