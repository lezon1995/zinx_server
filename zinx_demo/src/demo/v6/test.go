package main

import "time"

func main() {
	//intn := rand.Intn(1)
	//println(intn)
	for true {
		println(time.Now().Second())
		time.Sleep(time.Second)
	}
}
