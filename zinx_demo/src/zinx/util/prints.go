package util

import (
	"fmt"
	"time"
)

func PrintWaiting()  {
	fmt.Printf("Listening ")
	for true {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			fmt.Printf(".")
		}
		time.Sleep(time.Second)
		fmt.Printf("\b\b\b")
		fmt.Printf("   ")
		fmt.Printf("\b\b\b")
	}
}
