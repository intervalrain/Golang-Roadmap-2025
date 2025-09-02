package main

import (
	"fmt"
	"time"
)

// producer 函式會向 channel 中發送 5 個整數，然後關閉 channel
func producer(ch chan int) {
	fmt.Println("Producer: starting")
	for i := 0; i < 5; i++ {
		fmt.Printf("Producer: sending %d\n", i)
		ch <- i
		time.Sleep(100 * time.Millisecond)
	}
	// 當所有值都發送完畢後，關閉 channel
	// 這會通知接收方不會再有新的值傳入
	close(ch)
	fmt.Println("Producer: channel closed")
}

func main() {
	// --- Unbuffered Channel ---
	fmt.Println("--- Unbuffered Channel Example ---")
	// 建立一個無緩衝的 channel
	messages := make(chan string)

	go func() {
		// 這個傳送操作會阻塞，直到 main goroutine 準備好接收
		messages <- "ping"
	}()

	// 從 channel 接收一個值，這個操作會阻塞，直到有 goroutine 傳入值
	msg := <-messages
	fmt.Println("Received message:", msg)

	// --- Buffered Channel ---
	fmt.Println("\n--- Buffered Channel Example ---")
	// 建立一個緩衝大小為 2 的 channel
	bufferedChan := make(chan string, 2)

	// 因為有緩衝，這兩次傳送不會阻塞
	bufferedChan <- "buffered"
	bufferedChan <- "channel"
	fmt.Println("Sent 2 values to buffered channel without blocking")

	// 取出值
	fmt.Println("Received:", <-bufferedChan)
	fmt.Println("Received:", <-bufferedChan)

	// --- Range and Close ---
	fmt.Println("\n--- Range and Close Example ---")
	ch := make(chan int, 5)
	go producer(ch)

	// for-range 會持續從 channel 中接收值，直到 channel 被關閉
	fmt.Println("Consumer: waiting for values")
	for val := range ch {
		fmt.Printf("Consumer: received %d\n", val)
	}
	fmt.Println("Consumer: finished receiving")
}
