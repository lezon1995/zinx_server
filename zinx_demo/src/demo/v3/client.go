package main

import (
	"fmt"
	"net"
	"time"
	."zinx/util"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	time.Sleep(time.Second)
	if err != nil {
		fmt.Println("client start error , exit")
		return
	}
	var i uint32 = 1
	for true {
		_, err := conn.Write([]byte(fmt.Sprintf("hello zinx v0.2 : %d", i)))
		if err != nil {
			fmt.Println("client write error ", err)
			return
		}
		buf := NewBuffer()
		count, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client read error ", err)
			return
		}
		fmt.Printf("server call back : %s count=%d \n", buf, count)
		time.Sleep(time.Second)
		i++
	}
}
