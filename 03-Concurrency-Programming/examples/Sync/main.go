package main

import (
	"fmt"
	"sync"
)

// --- sync.Mutex Example ---

// SafeCounter 是一個線程安全的計數器
type SafeCounter struct {
	mu      sync.Mutex
	counter int
}

// Inc 方法安全地對計數器進行遞增
func (c *SafeCounter) Inc() {
	// 在修改計數器前鎖定，防止競爭條件
	c.mu.Lock()
	// 使用 defer 確保在函式退出時一定會解鎖
	defer c.mu.Unlock()
	c.counter++
}

// Value 方法安全地讀取計數器的值
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter
}

// --- sync.Once Example ---

var once sync.Once

func initialize() {
	fmt.Println("This will be printed only once.")
}

func main() {
	// --- Mutex Demo ---
	fmt.Println("---\ sync.Mutex Example ---")
	sc := SafeCounter{counter: 0}
	var wg sync.WaitGroup

	// 啟動 1000 個 goroutine 來並發地增加計數器
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sc.Inc()
		}()
	}

	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("Final counter value:", sc.Value()) // 如果沒有 Mutex 保護，結果將不確定

	// --- Once Demo ---
	fmt.Println("\n--- sync.Once Example ---")
	var onceWg sync.WaitGroup

	// 啟動 10 個 goroutine 來嘗試執行 initialize()
	for i := 0; i < 10; i++ {
		onceWg.Add(1)
		go func(id int) {
			defer onceWg.Done()
			fmt.Printf("Gorotuine %d trying to initialize...\n", id)
			// 儘管有 10 個 goroutine，initialize() 只會被執行一次
			once.Do(initialize)
		}(i)
	}

	onceWg.Wait()
	fmt.Println("Done.")
}
