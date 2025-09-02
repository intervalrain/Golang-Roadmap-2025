package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// loggingMiddleware 是一個記錄請求日誌的中介軟體
func loggingMiddleware(next http.Handler) http.Handler {
	// http.HandlerFunc 是一個轉接器，讓普通函式可以作為 http.Handler 使用
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// 呼叫鏈中的下一個處理器 (可能是另一個中介軟體，或最終的處理函式)
		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}

// helloHandler 是我們最終的業務邏輯處理函式
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello with Middleware!")
}

func main() {
	// 建立一個最終的處理器
	finalHandler := http.HandlerFunc(helloHandler)

	// 使用 loggingMiddleware 包裹最終的處理器
	http.Handle("/hello", loggingMiddleware(finalHandler))

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
