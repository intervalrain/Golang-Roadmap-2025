package main

import (
	"fmt"
	"time"
)

// --- Fan-Out, Fan-In Pattern ---

// worker 函式從 jobs channel 接收任務，並將結果發送到 results channel。
func worker(id int, jobs <-chan int, results chan<- string) {
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
		// 模擬耗時的計算
		time.Sleep(time.Second)
		resultStr := fmt.Sprintf("Worker %d finished job %d", id, j)
		results <- resultStr
	}
}

func fanOutFanIn() {
	fmt.Println("--- Fan-Out, Fan-In Example ---")
	numJobs := 10
	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	// Fan-Out: 啟動 3 個 worker goroutines 來並行處理任務
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// 將任務發送到 jobs channel
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Fan-In: 收集所有任務的結果
	for a := 1; a <= numJobs; a++ {
		fmt.Println("Result:", <-results)
	}
	fmt.Println("All jobs finished.")
}

// --- Rate Limiting Pattern ---

func rateLimiting() {
	fmt.Println("\n--- Rate Limiting Example ---")

	// 建立一個 Ticker，它會以固定的時間間隔向其 channel 發送事件
	// 這裡我們設定為每 500 毫秒一次
	ticker := time.NewTicker(500 * time.Millisecond)
	// 確保在函式結束時停止 ticker，以釋放資源
	defer ticker.Stop()

	// 模擬有 5 個請求需要處理
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	// 處理請求
	for req := range requests {
		// 等待 ticker 觸發。這會阻塞當前 goroutine，
		// 從而將請求處理的速率限制在每 500 毫秒一次。
		<-ticker.C
		fmt.Printf("Processing request %d at %v\n", req, time.Now().Format("15:04:05.000"))
	}
	fmt.Println("All requests processed.")
}

func main() {
	fanOutFanIn()
	rateLimiting()
}
