package main

import (
	"fmt"
	"time"
)

func main() {
	// --- Basic Select ---
	fmt.Println("---"Basic Select Example---")
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// 等待 c1 和 c2 的消息，總共兩次
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("Received from c1:", msg1)
		case msg2 := <-c2:
			fmt.Println("Received from c2:", msg2)
		}
	}

	// --- Select with Timeout ---
	fmt.Println("\n---"Select with Timeout Example---")
	cr := make(chan string, 1)
	go func() {
		// 這個 Goroutine 需要 2 秒才能完成
		time.Sleep(2 * time.Second)
		cr <- "result"
	}()

	select {
	case res := <-cr:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		// 但我們只等待 1 秒
		fmt.Println("Timeout waiting for result")
	}

	// --- Select with Default (Non-blocking) ---
	fmt.Println("\n---"Select with Default (Non-blocking) Example---")
	messages := make(chan string)

	// 嘗試接收 messages，但沒有 Goroutine 在發送，所以會立即執行 default
	select {
	case msg := <-messages:
		fmt.Println("Received message:", msg)
	default:
		fmt.Println("No message received.")
	}

	// 嘗試發送，但沒有 Goroutine 在接收，也會立即執行 default
	select {
	case messages <- "hi":
		fmt.Println("Sent message")
	default:
		fmt.Println("No message sent.")
	}
}
