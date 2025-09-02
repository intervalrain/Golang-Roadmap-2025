package main

import (
	"fmt"
	"sync"
	"time"
)

// worker 函式模擬一個需要一些時間來完成的工作。
// 它接收一個指向 sync.WaitGroup 的指標，以便在完成時通知主程式。
func worker(id int, wg *sync.WaitGroup) {
	// defer 確保在函式返回前，一定會呼叫 wg.Done()
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)

	// 模擬耗時的工作
	time.Sleep(time.Second)

	fmt.Printf("Worker %d done\n", id)
}

func main() {
	// WaitGroup 用於等待一組 Goroutine 完成。
	var wg sync.WaitGroup

	fmt.Println("--- WaitGroup Example ---")

	// 啟動 3 個 worker goroutines。
	for i := 1; i <= 3; i++ {
		// 每啟動一個 Goroutine，計數器就加 1。
		wg.Add(1)

		// 使用 go 關鍵字啟動一個新的 Goroutine。
		go worker(i, &wg)
	}

	// Wait() 會阻塞，直到 WaitGroup 的計數器變為 0。
	fmt.Println("Waiting for workers to finish...")
	wg.Wait()
	fmt.Println("All workers done.")

	fmt.Println("\n--- Anonymous Goroutine Example ---")
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("I am an anonymous goroutine!")
		time.Sleep(500 * time.Millisecond)
	}()

	fmt.Println("Waiting for the anonymous goroutine...")
	wg.Wait()
	fmt.Println("Anonymous goroutine finished.")
}
