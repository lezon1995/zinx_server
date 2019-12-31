package zimpl

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"testing"
	"time"
	"zinx/util"
)

func TestDataPack_GetHeadLen(t *testing.T) {

}

func TestDataPack(t *testing.T) {

	//启动server
	listen, _ := net.Listen("tcp", "127.0.0.1:7777")

	//循环监听
	go func() {
		for true {
			conn, _ := listen.Accept()

			//开协程模拟server handler
			go func(conn net.Conn) {
				dataPack := NewDataPack()

				for true {
					//初始化包头
					headData := util.NewLenBuffer(dataPack.GetHeadLen())

					//等待读取 读满headData
					_, _ = io.ReadFull(conn, headData)

					//拆包头 获取IMessage对象
					message, _ := dataPack.Unpack(headData)

					//判断dataLen是否大于0
					if message.GetMsgLen() > 0 {
						//代表有数据需要读取 需要二次读取
						buffer := util.NewLenBuffer(message.GetMsgLen())

						message.SetData(buffer)

						//从conn中二次读取MsgLen长度的数据
						_, _ = io.ReadFull(conn, message.GetData())

						fmt.Println("---> receive MsgId=", message.GetMsgId(), "dataLen=", message.GetMsgLen(), "msg=", string(message.GetData()))
					}
				}
			}(conn)

		}
	}()

	//模拟客户端
	conn1, _ := net.Dial("tcp", "127.0.0.1:7777")

	pack := NewDataPack()

	//一次性发送给服务端
	for true {
		intn := rand.Intn(time.Now().Second())
		//模拟粘包过程 封装两个IMessage一起发送

		//封装第一个IMessage
		bytes1 := make([]byte, intn)
		for i := 0; i < intn; i++ {
			r1 := rand.Intn(time.Now().Second() + 1)
			i2 := r1 % 26
			bytes1[i] = byte(97 + i2)
		}
		message1 := NewMessage1(
			1,
			bytes1,
		)
		data1, _ := pack.Pack(message1)

		//封装第二个IMessage
		intn = rand.Intn(100)
		bytes2 := make([]byte, intn)
		for i := 0; i < intn; i++ {
			r1 := rand.Intn(time.Now().Second() + 1)
			i2 := r1 % 26
			bytes2[i] = byte(97 + i2)
		}
		message2 := NewMessage1(
			2,
			bytes2,
		)
		data2, _ := pack.Pack(message2)

		//将两个包粘在一起
		data := append(data1, data2...)

		_, _ = conn1.Write(data)
		time.Sleep(time.Second)
	}

	//阻塞
	select {}
}
