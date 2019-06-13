package main

import (
	"fmt"
	"net"
	"time"
)

//模拟客户端
func main() {

	fmt.Println("Client Start...")

	time.Sleep(1 * time.Second)

	//连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Client Start Error, Exit!")
		return
	}

	for {
		//调用写方法
		_, err := conn.Write([]byte("Hello World!V0.2"))
		if err != nil {
			fmt.Println("Write conn Error", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf Error", err)
			return
		}

		fmt.Printf("Server call back:%s, cnt = %d\n", buf, cnt)

		//CPU阻塞
		time.Sleep(1 * time.Second)

	}

}
