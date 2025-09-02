package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		// time.Sleep 讓 Goroutine 暫停一下，以便觀察交錯執行的效果
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println("---\nGoroutines Intro---")

	// 使用 `go` 關鍵字啟動一個新的 Goroutine
	// `say("World")` 將會和 `main` 函式並發執行
	go say("World")

	// 在主 Goroutine 中執行 say("Hello")
	// 主 Goroutine 的執行會給 `say("World")` Goroutine 一些執行的時間
	say("Hello")

	fmt.Println("\n---\nGoroutine Exit Demo---")
	// 在這個例子中，`main` 函式可能在 Goroutine 開始執行前就退出了
	// 因此您可能看不到 "I am a goroutine" 的輸出
	go func() {
		fmt.Println("I am a goroutine")
	}()

	// 我們在這裡短暫睡眠，只是為了演示目的，以增加看到上面 Goroutine 輸出的機會。
	// 這不是一個可靠的同步方法！
	time.Sleep(50 * time.Millisecond)

	fmt.Println("Main function finished.")
}
