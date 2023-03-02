package main

import (
	"fmt"
	"github.com/lemoba/mmo-game/proto/pb"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 定义一个Person结构体
	person := &pb.Person{
		Name:   "ranen",
		Age:    27,
		Emails: []string{"lemoba@qq.com", "ranen1024@gmail.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "15549959926",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "15549959927",
				Type:   pb.PhoneType_WORK,
			},
			&pb.PhoneNumber{
				Number: "15549959928",
				Type:   pb.PhoneType_HOME,
			},
		},
	}

	// 编码
	// 将对象进行序列化，得到一个二进制文件，进行传输
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("Marshal error: ", err)
	}

	// 解码
	newData := &pb.Person{}
	if err := proto.Unmarshal(data, newData); err != nil {
		fmt.Println("unmarshal error: ", err)
	}

	fmt.Println("源数据: ", person)
	fmt.Println("newData: ", newData)
}
