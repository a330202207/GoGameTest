package main

import (
	"fmt"
	"gFrame/protobufDemo/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	//定义一个Person结构对象

	person := &pb.Person{
		Name:   "Ned",
		Age:    28,
		Emails: []string{"ned@163.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "33333333",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "22222222",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "11111111",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	//将person对象序列化（就是将Protobuf的message进行序列化）,得到一个二进制的文件
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("Marshal Err:", err)
	}

	newData := &pb.Person{}

	err = proto.Unmarshal(data, newData)
	if err != nil {
		fmt.Println("Unmarshal Err:", err)
	}

	fmt.Println("源数据:", person)
	fmt.Println("解码之后的数据:", newData)
}
