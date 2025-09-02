package main

import (
	"fmt"
	"log"
	"net/http"
)

// helloHandler 處理對 /hello 路徑的請求
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 檢查請求方法，只允許 GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Hello, World from net/http!")
}

func main() {
	// 將 helloHandler 函式註冊到 "/hello" 路徑
	http.HandleFunc("/hello", helloHandler)

	// 啟動伺服器，監聽 8080 埠
	fmt.Println("Server starting on http://localhost:8080")
	// ListenAndServe 會一直阻塞，直到發生錯誤
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
