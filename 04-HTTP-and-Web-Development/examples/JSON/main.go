package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// User struct 用於定義我們的資料模型
type User struct {
	// struct tag 用於控制 JSON 的欄位名稱
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// usersHandler 處理對 /users 路徑的請求
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getUser 處理 GET 請求，回傳一個 User 的 JSON
func getUser(w http.ResponseWriter, r *http.Request) {
	user := User{ID: 1, Name: "Alice"}

	// 設定回應的 Content-Type 為 application/json
	w.Header().Set("Content-Type", "application/json")

	// 使用 json.NewEncoder 將 user struct 編碼為 JSON 並寫入 ResponseWriter
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// createUser 處理 POST 請求，從請求主體中解碼 User 的 JSON
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	// 使用 json.NewDecoder 從請求的 Body 中讀取並解碼 JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Created user: %+v\n", user)

	// 回傳 201 Created 狀態碼
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created: %s", user.Name)
}

func main() {
	http.HandleFunc("/users", usersHandler)

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println(`Try GET and POST on http://localhost:8080/users`)
	fmt.Println(`POST example: curl -X POST -d '{"name":"Bob"}' http://localhost:8080/users`)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
