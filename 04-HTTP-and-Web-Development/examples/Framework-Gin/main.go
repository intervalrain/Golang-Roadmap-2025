package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// User struct 用於定義我們的資料模型
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// 初始化 Gin 引擎，使用預設的中介軟體 (Logger 和 Recovery)
	r := gin.Default()

	// 定義一個 GET 路由
	r.GET("/ping", func(c *gin.Context) {
		// c.JSON 是一個方便的函式，可以將 struct 或 map 序列化為 JSON 並回傳
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 帶有路徑參數的路由
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "Getting user with ID: " + id,
		})
	})

	// 定義一個 POST 路由
	r.POST("/users", func(c *gin.Context) {
		var user User

		// c.ShouldBindJSON 會將請求的 JSON 主體綁定到 user struct 上
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "User created: " + user.Name,
		})
	})

	// 啟動伺服器，預設監聽在 :8080
	// 在執行前，請確保您已經使用 `go get github.com/gin-gonic/gin` 安裝了 Gin
	r.Run()
}
