package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
	"zinx/util"
	"zinx/zimpl"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	time.Sleep(time.Second)
	if err != nil {
		fmt.Println("client start error , exit")
		return
	}
	pack := zimpl.NewDataPack()
	go func() {
		for true {
			intn := rand.Intn(time.Now().Second())
			//模拟粘包过程 封装两个IMessage一起发送

			//封装第一个IMessage
			bytes1 := util.NewLenBuffer(uint32(intn))
			for i := 0; i < intn; i++ {
				r1 := rand.Intn(time.Now().Second() + 10)
				i2 := r1 % 26
				bytes1[i] = byte(97 + i2)
			}
			message1 := zimpl.NewMessage1(
				1,
				bytes1,
			)
			data1, _ := pack.Pack(message1)

			//封装第二个IMessage
			intn = rand.Intn(100)
			bytes2 := util.NewLenBuffer(uint32(intn))
			for i := 0; i < intn; i++ {
				r1 := rand.Intn(time.Now().Second() + 10)
				i2 := r1 % 26
				bytes2[i] = byte(97 + i2)
			}
			message2 := zimpl.NewMessage1(
				2,
				bytes2,
			)
			data2, _ := pack.Pack(message2)

			//将两个包粘在一起
			data := append(data1, data2...)

			_, _ = conn.Write(data)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		dataPack := zimpl.NewDataPack()
		for true {
			headData := util.NewLenBuffer(dataPack.GetHeadLen())
			io.ReadFull(conn, headData)

			message, _ := dataPack.Unpack(headData)
			var data []byte
			if message.GetMsgLen() > 0 {
				data = util.NewLenBuffer(message.GetMsgLen())
				io.ReadFull(conn, data)
			}
			message.SetData(data)
			fmt.Println("--->client receive MsgId=", message.GetMsgId(), "msg=", string(message.GetData()))
		}
	}()

	select {}
}
