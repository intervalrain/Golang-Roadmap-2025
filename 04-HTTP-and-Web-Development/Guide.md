# Chapter 4.1: HTTP Basics

在深入 Go 的 Web 開發之前，理解 HTTP (Hypertext Transfer Protocol) 的基礎至關重要。HTTP 是一個客戶端-伺服器協議，它規範了網頁瀏覽器 (Client) 和網頁伺服器 (Server) 之間的溝通方式。

- **請求 (Request)**: 客戶端向伺服器發送一個請求。一個請求包含：
    - **方法 (Method)**: 定義了要執行的動作，例如 `GET` (獲取資源), `POST` (建立資源), `PUT` (更新資源), `DELETE` (刪除資源)。
    - **路徑 (Path)**: 指定了要操作的資源，例如 `/users/123`。
    - **標頭 (Headers)**: 包含額外的資訊，如 `Content-Type`, `Authorization`。
    - **主體 (Body)**: 包含要傳送的資料，例如 `POST` 或 `PUT` 請求中的 JSON 資料。

- **回應 (Response)**: 伺服器回傳一個回應給客戶端。一個回應包含：
    - **狀態碼 (Status Code)**: 一個三位數的數字，表示請求的結果。例如 `200 OK` (成功), `404 Not Found` (找不到資源), `500 Internal Server Error` (伺服器錯誤)。
    - **標頭 (Headers)**: 包含額外的資訊，如 `Content-Type`。
    - **主體 (Body)**: 包含請求的資源內容，例如 HTML 頁面或 JSON 資料。

---

# Chapter 4.2: net/http Package

Go 的標準庫 `net/http` 提供了建構 HTTP 伺服器和客戶端所需的一切。它的設計非常優雅且高效。

## 1. 建立一個簡單的伺服器

一個最基本的 Go Web 伺服器包含兩個步驟：

1.  **註冊處理函式 (Handler Function)**: 您需要編寫一個函式來處理來自特定 URL 路徑的請求。這個函式必須符合 `http.HandlerFunc` 型別，即 `func(http.ResponseWriter, *http.Request)`。
2.  **啟動伺服器**: 使用 `http.ListenAndServe` 來啟動伺服器並監聽指定的埠號。

```go
package main

import (
    "fmt"
    "net/http"
)

// helloHandler 處理對 /hello 路徑的請求
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    // 將 helloHandler 函式註冊到 "/hello" 路徑
    http.HandleFunc("/hello", helloHandler)

    // 啟動伺服器，監聽 8080 埠
    fmt.Println("Server starting on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
```

- `http.ResponseWriter`: 這是一個介面，伺服器用它來建構 HTTP 回應。我們使用 `fmt.Fprintf` 將 "Hello, World!" 寫入其中。
- `*http.Request`: 這是一個結構，包含了客戶端 HTTP 請求的所有資訊，如 URL、標頭和主體。

---

# Chapter 4.3: Routing

**Routing** (路由) 是指將不同的請求 URL 導向到不同處理函式的過程。`net/http` 套件提供了一個預設的路由器 `DefaultServeMux`，`http.HandleFunc` 就是在它上面註冊路徑。

```go
func main() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/users", usersHandler)
    http.ListenAndServe(":8080", nil)
}
```

標準庫的 `ServeMux` 功能比較基礎，它只支援靜態路徑（例如 `/users`）和以 `/` 結尾的子樹路徑（例如 `/static/`）。它不支援帶有變數的路徑（例如 `/users/:id`）。

對於更複雜的路由需求，Go 社群開發了許多優秀的第三方路由套件，例如 `gorilla/mux` 和 `chi`，它們也是許多 Web 框架的核心。

---

# Chapter 4.4: Middleware

**Middleware** (中介軟體) 是一個非常強大的概念。它是一個函式，包裹在您的主要處理函式之外，用於在處理請求之前或之後執行一些通用的操作，例如：

- 記錄請求日誌 (Logging)
- 身份驗證 (Authentication)
- 壓縮回應 (Compression)
- 新增安全標頭 (Security Headers)

Middleware 的本質是一個接收 `http.Handler` 並回傳一個新的 `http.Handler` 的函式。

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        fmt.Printf("Started %s %s\n", r.Method, r.URL.Path)

        // 呼叫下一個中介軟體或最終的處理函式
        next.ServeHTTP(w, r)

        fmt.Printf("Completed in %v\n", time.Since(start))
    })
}

func main() {
    finalHandler := http.HandlerFunc(helloHandler)
    http.Handle("/hello", loggingMiddleware(finalHandler))
    http.ListenAndServe(":8080", nil)
}
```

---

# Chapter 4.5: JSON Handling

在現代 Web 開發中，JSON (JavaScript Object Notation) 是 API 之間交換資料最常用的格式。Go 的 `encoding/json` 套件提供了完整的 JSON 支援。

- **編碼 (Encoding/Marshalling)**: 將 Go 的 `struct` 轉換為 JSON 字串。
- **解碼 (Decoding/Unmarshalling)**: 將 JSON 字串轉換為 Go 的 `struct`。

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// JSON 編碼範例
func getUser(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "Alice"}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// JSON 解碼範例
func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Printf("Created user: %+v\n", user)
    w.WriteHeader(http.StatusCreated)
}
```

`json:"..."` 這種標籤 (struct tag) 用於控制 `struct` 欄位在 JSON 中的表現形式。

---

# Chapter 4.6: Web Frameworks

雖然 Go 的標準庫非常強大，但在建構大型複雜應用時，從零開始可能會比較繁瑣。Web 框架提供了一套更高層次的抽象和工具，來簡化開發過程。

**為何使用框架？**

- 更強大和靈活的路由
- 更容易管理的中介軟體鏈
- 內建的資料驗證、渲染等功能
- 統一的專案結構

**流行的 Go Web 框架**

- **Gin**: 非常流行，以高效能和類似 Martini 的 API 而聞名。
- **Echo**: 高效能、可擴充、極簡的框架。
- **Fiber**: 受 Express.js 啟發，是 Go 領域中最快的框架之一。

### Gin 範例

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    r.Run() // 預設監聽在 :8080
}
```

---

## Conclusion

恭喜！您已經完成了第四章，掌握了在 Go 中進行 Web 開發的核心技能。從使用 `net/http` 建構基礎伺服器，到處理路由、中介軟體和 JSON，再到理解 Web 框架的價值，您現在已經具備了開發高效能 Web 服務和 API 的能力。

```